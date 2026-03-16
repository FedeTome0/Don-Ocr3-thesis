package main

import (
	"OCR3-thesis/contracts"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	shell "github.com/ipfs/go-ipfs-api"
)


// startChainListener initializes the WebSocket connection to the blockchain and
// listens for OracleQueue events. It implements a resilient reconnection loop
// to ensure the oracle never misses an event even if the RPC node temporarily drops.
func startChainListener(ctx context.Context, rpcUrl, verifierContractAddress, queueContractAddr string, ipfs *shell.Shell) {
	// Convert standard HTTP RPC URL to WebSocket for real-time event streaming
	wsUrl := strings.Replace(rpcUrl, "http://", "ws://", 1)

	fmt.Printf("Listener: Connecting to %s\n    -> Watching Queue: %s\n", wsUrl, queueContractAddr)

	// Resilient reconnection loop: if runListener returns an error, wait and retry
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Listener: Context cancelled, shutting down")
			return
		default:
			// Blocking call that listens for events. Returns only on error.
			if err := runListener(ctx, wsUrl, verifierContractAddress, queueContractAddr, ipfs); err != nil {
				log.Printf("Listener crashed: %v. Reconnecting in 0.5s...", err)
				// Changed from 5 sec to 500ms because of time evaluation
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

// runListener handles ABI binding, log subscription, and event parsing.
// It returns an error on any failure, allowing startChainListener to reconnect.
func runListener(ctx context.Context, wsUrl, verifierContractAddress, queueContractAddr string, ipfs *shell.Shell) error {
	client, err := ethclient.Dial(wsUrl)
	if err != nil {
		return fmt.Errorf("Failed RPC connection: %w", err)
	}
	defer client.Close()

	qAddr := common.HexToAddress(queueContractAddr)
	vAddr := common.HexToAddress(verifierContractAddress)

	// Bind the OracleQueue contract ABI
	queueContract, err := contracts.NewOracleQueue(qAddr, client)
	if err != nil {
		return fmt.Errorf("Failed to bind OracleQueue: %w", err)
	}

	// Bind the OracleVerifier contract ABI
	verifierContract, err := contracts.NewOracleVerifier(vAddr, client)
	if err != nil {
		return fmt.Errorf("Failed to bind OracleVerifier: %w", err)
	}

	logs := make(chan gethtypes.Log, 100)

	// Subscribe to logs from both the Queue and Verifier contracts
	query := ethereum.FilterQuery{
		Addresses: []common.Address{qAddr, vAddr},
	}

	sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		return fmt.Errorf("Failed to subscribe: %w", err)
	}
	defer sub.Unsubscribe()

	fmt.Println("Listener: Listening for LogNewJobForOracles events...")

	for {
		// Inactivity timer: force reconnect if no events are received within the window.
		// This guards against silent WebSocket drops that don't emit a proper close frame.
		timeout := time.NewTimer(20 * time.Second)

		select {
		case err := <-sub.Err():
			timeout.Stop()
			return fmt.Errorf("Subscription error: %w", err)

		case vLog := <-logs:
			timeout.Stop()

			// Attempt to parse as LogNewJobForOracles
			newJob, err := queueContract.ParseLogNewJobForOracles(vLog)
			if err == nil {
				// Handle the job asynchronously to avoid blocking the listener loop
				if err := handleNewJob(ctx, newJob, ipfs); err != nil {
					log.Printf("Error handling event: %v", err)
				}
				continue
			}

			// Attempt to parse as JobCompleted
			completed, err := verifierContract.ParseJobCompleted(vLog)
			if err == nil {
				jobId64 := completed.JobId.Uint64()
				MarkJobAsProcessed(jobId64)
				fmt.Printf("[Listener] Job #%d completed on-chain by %s. Cache updated.\n",
					jobId64, completed.Submitter.Hex())
				continue
			}

		case <-timeout.C:
			// No events received within the timeout window.
			// The WebSocket connection may have silently dropped — force a reconnect.
			fmt.Println("[Listener] No events received in 60s. Forcing WebSocket reconnection...")
			return fmt.Errorf("websocket silent drop timeout")

		case <-ctx.Done():
			timeout.Stop()
			fmt.Println("Listener: Shutdown signal received")
			return nil
		}
	}
}

// =============================================================================
// JOB HANDLING & CACHE MANAGEMENT
// =============================================================================

// handleNewJob receives a parsed LogNewJobForOracles event, registers the job
// in the shared thread-safe cache as PENDING, and delegates heavy processing
// to a background goroutine so the listener can return immediately to the
// WebSocket stream.
func handleNewJob(ctx context.Context, event *contracts.OracleQueueLogNewJobForOracles, ipfs *shell.Shell) error {
	jobId := event.JobId
	ipfsCid := event.IpfsCid
	jobId64 := jobId.Uint64()

	fmt.Printf("\nListener: Detected Job #%s | CID: %s\n", jobId.String(), ipfsCid)

	// -------------------------------------------------------------------------
	// 1. THREAD-SAFE CACHE UPDATE
	// Lock the RWMutex before writing to prevent race conditions with the OCR3
	// plugin, which may be reading the cache concurrently.
	// -------------------------------------------------------------------------
	JobCache.Lock()

	// Skip if this node has already registered this job
	if _, exists := JobCache.jobs[jobId64]; exists {
		JobCache.Unlock()
		return nil
	}

	// Track the highest job ID seen so far
	if JobCache.LatestJobID.Cmp(big.NewInt(-1)) == 0 || jobId.Cmp(JobCache.LatestJobID) > 0 {
		JobCache.LatestJobID = jobId
	}

	// Enqueue the job as PENDING. OCR3 Observation() will return empty bytes
	// until the state transitions to Completed.
	JobCache.enqueue(jobId64, &JobData{
		JobID:     jobId,
		CID:       ipfsCid,
		State:     StatePending,
		Processed: false,
	})
	JobCache.Unlock()

	fmt.Printf("[Listener] Job #%d saved as PENDING. Starting async IPFS & AI task...\n", jobId64)

	// -------------------------------------------------------------------------
	// 2. ASYNC DELEGATION
	// Spawn a goroutine for all I/O-bound and CPU-bound work (IPFS download,
	// AI inference polling). The listener loop is not blocked.
	// -------------------------------------------------------------------------
	go processJobAsync(ctx, jobId64, ipfsCid, ipfs)

	return nil
}

// =============================================================================
// ASYNCHRONOUS WORKER (IPFS & AI Computation)
// =============================================================================

// MODEL_SERVICE_URL is the base URL of the Python inference service.
const MODEL_SERVICE_URL = "http://host.docker.internal:9090"

// processJobAsync is the background worker responsible for:
//  1. Downloading the job payload from IPFS
//  2. Submitting it to the Python AI service and polling for the result
//  3. Committing the final result vector to the shared cache
func processJobAsync(ctx context.Context, jobIdUint uint64, ipfsCid string, ipfs *shell.Shell) {
	// -------------------------------------------------------------------------
	// PHASE 1: IPFS DOWNLOAD
	// Fetch the raw input data associated with this job from IPFS.
	// -------------------------------------------------------------------------
	ipfsReader, err := ipfs.Cat(ipfsCid)
	if err != nil {
		setJobError(jobIdUint, fmt.Errorf("IPFS download error: %w", err))
		return
	}
	defer ipfsReader.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(ipfsReader); err != nil {
		setJobError(jobIdUint, fmt.Errorf("IPFS read error: %w", err))
		return
	}
	inputText := buf.String()
	fmt.Printf("[Async Task] Job #%d - IPFS data downloaded: %d bytes\n", jobIdUint, len(inputText))

	// -------------------------------------------------------------------------
	// PHASE 2: AI ATTRIBUTION (Python service)
	// Send the input text to the external attribution model and poll for the
	// result. The computation can take several minutes, so running it here
	// (outside the OCR3 round) prevents libocr timeout loops.
	// -------------------------------------------------------------------------
	// Convert the numerical ID to a string format expected by the Python backend
	jobIdStr := fmt.Sprintf("%d", jobIdUint)

	fmt.Printf("[Async Task] Job #%d - Submitting to AI service...\n", jobIdUint)

	// Prepare the payload struct with the job details
	reqData := AttributeRequest{
		JobId: jobIdStr,
		Text:  inputText,
	}
	// Serialize the request payload into JSON
	jsonData, err := json.Marshal(reqData)
	if err != nil {
		setJobError(jobIdUint, fmt.Errorf("json marshal error: %w", err))
		return
	}

	// =========================================================================
    // 1. ASYNCHRONOUS JOB SUBMISSION (POST)
    // =========================================================================
	// Use NewRequestWithContext to bind the HTTP request to the goroutine's lifecycle.
    // This prevents goroutine leaks if the parent context is canceled during the request.
	req, err := http.NewRequestWithContext(ctx, "POST", MODEL_SERVICE_URL+"/attribute", bytes.NewBuffer(jsonData))
	if err != nil {
		setJobError(jobIdUint, fmt.Errorf("failed to build request: %w", err))
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the POST request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		setJobError(jobIdUint, fmt.Errorf("failed to call Python server: %w", err))
		return
	}
	// Defer body closure to prevent TCP connection and memory leaks
	defer res.Body.Close()

	// The Python server returns 202 for new jobs and 200 for deduplicated/existing jobs
	if res.StatusCode != 200 && res.StatusCode != 202 {
		setJobError(jobIdUint, fmt.Errorf("Python server returned status: %d", res.StatusCode))
		return
	}

	// =========================================================================
    // 2. NON-BLOCKING POLLING LOOP (GET)
    // =========================================================================
	fmt.Printf("[Async Task] Job #%d - Waiting for AI result...\n", jobIdUint)

	// Initialize a ticker to poll the endpoint every 5 seconds.
    // Deferring Stop() ensures the timer is released from memory when the function exits.
	ticker := time.NewTicker(4 * time.Second) // Timer changed from 15 to 5, is based on the overall time of the att. method
	defer ticker.Stop()

	var finalResultStrings []string
	statusUrl := fmt.Sprintf("%s/result/%s", MODEL_SERVICE_URL, jobIdStr)

	// Infinite loop using a select statement for safe multiplexing of channels
	for {
		select {
		// CASE A: The parent context signals a shutdown or timeout.
        // We catch this to gracefully abort the polling and exit the goroutine.
		case <-ctx.Done():
			fmt.Printf("[Async Task] Job #%d - Cancelled (shutdown).\n", jobIdUint)
			return

		// CASE B: The ticker asks every 15s. Time to check the job status.
		case <-ticker.C:
			req, err := http.NewRequestWithContext(ctx, "GET", statusUrl, nil)
			if err != nil {
				continue
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// If the network temporarily drops, we log the warning but keep the loop alive
				fmt.Printf("[Async Task] Job #%d - Status check failed (%v). Retrying...\n", jobIdUint, err)
				continue
			}

			var pyResp AttributeResponse
			if err := json.NewDecoder(resp.Body).Decode(&pyResp); err != nil {
				resp.Body.Close()
				fmt.Printf("[Async Task] Job #%d - JSON decode error. Retrying...\n", jobIdUint)
				continue
			}
			resp.Body.Close()

			// Cases evaluation based on the Python server's response
			switch pyResp.Status {
			case "completed":
				// Target state reached: extract results and break out of the infinite loop
				finalResultStrings = pyResp.Result
				goto donePolling
			case "error":
				// Fatal state reached: log the AI error and terminate the goroutine
				setJobError(jobIdUint, fmt.Errorf("Python AI error: %s", pyResp.Error))
				return
			default:
				// Intermediate states ("queued", "processing"): log and wait for the next tick
				fmt.Printf("[Async Task] Job #%d - Status: %s. Waiting...\n", jobIdUint, pyResp.Status)
			}
		}
	}
donePolling:

	// -------------------------------------------------------------------------
	// PHASE 3: CACHE COMMIT
	// Parse the result strings returned by Python into []*big.Int (Solidity
	// compatible) and expose the final vector to the OCR3 consensus round.
	// -------------------------------------------------------------------------
	var finalVector []*big.Int
	for _, strVal := range finalResultStrings {
		val := new(big.Int)
		val, success := val.SetString(strVal, 10)
		if !success {
			setJobError(jobIdUint, fmt.Errorf("failed to parse BigInt from string: %s", strVal))
			return
		}
		finalVector = append(finalVector, val)
	}

	JobCache.Lock()
	if job, exists := JobCache.jobs[jobIdUint]; exists {
		job.State = StateCompleted // Mark as ready for consensus
		job.Result = finalVector
	}
	JobCache.Unlock()

	fmt.Printf("[Async Task] Job #%d - COMPLETED. Result vector: %d elements.\n", jobIdUint, len(finalVector))
}