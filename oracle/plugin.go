package main

import (
	"context"
	"os"
	//"encoding/json"
	"fmt"
	//"math"
	"math/big"
	"sort"
	//"strings"
	"time"

	//"github.com/ethereum/go-ethereum/accounts/abi"
	//"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/common/hexutil"
	//"github.com/ethereum/go-ethereum/rpc" 
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/quorumhelper"
)

// =============================================================================
// REPORTING PLUGIN (OCR3 Application Logic)
// =============================================================================
type attributionPluginFactory struct {
	//oracleID 	int
	//numOracles 	int
	//rpcUrl		string
	//contractAddr common.Address 
}

// NewReportingPlugin is invoked by libocr during node bootstrap.
// It initializes the plugin with the protocol limits and injects the 
// custom dependencies required for the attribution process.
func (f attributionPluginFactory) NewReportingPlugin(_ context.Context, cfg ocr3types.ReportingPluginConfig) (ocr3types.ReportingPlugin[struct{}], ocr3types.ReportingPluginInfo, error) {
	return &attributionPlugin{
		cfg:      cfg,
		//oracleIndex: f.oracleID,
		//numOracles: f.numOracles,
		//rpcUrl:		f.rpcUrl,
		//contractAddr: f.contractAddr,
	}, ocr3types.ReportingPluginInfo{
		Name: "attribution-matrix-median",
		Limits: ocr3types.ReportingPluginLimits{
			MaxQueryLength:       10 * 1024,
			MaxObservationLength: 256 * 1024, // 256KB limit
			MaxOutcomeLength:     256 * 1024,
			MaxReportLength:      256 * 1024,
			MaxReportCount:       1,
		},
	}, nil
}

type attributionPlugin struct {
	cfg      ocr3types.ReportingPluginConfig
	//oracleIndex int
	//numOracles int
	//rpcUrl 		string
	//contractAddr common.Address
}


// =============================================================================
// PHASE 1: QUERY (Leader proposes work)
// =============================================================================

// Query is executed exclusively by the elected round leader.
// It inspects the local thread-safe cache to identify pending or completed jobs.
// The FIFO queue ensures fair scheduling and prevents head-of-line blocking.
func (p *attributionPlugin) Query(ctx context.Context, _ ocr3types.OutcomeContext) (ocrtypes.Query, error) {
    start := time.Now()
    defer func() {
        fmt.Printf("[METRIC-OCR] Phase: QUERY | Leader Node | Time: %v\n", time.Since(start))
    }()

	// Enforce mutual exclusion while managing the FIFO queue state
	JobCache.Lock()
    defer JobCache.Unlock()

	// Iteratively process the queue to flush out stale or failed jobs
	for len(JobCache.queue) > 0 {
        headID := JobCache.queue[0]
        job, exists := JobCache.jobs[headID]

        if !exists || job.Processed {
            // Job already finalized or missing. Dequeue and proceed
            JobCache.queue = JobCache.queue[1:]
            continue
        }

		// Fault Tolerance: If the asynchronous worker encounters a fatal error, 
        // the job is dequeued to prevent consensus stalling.
        if job.State == StateFailed {
            fmt.Printf("[WARN] Job #%s failed permanently. Removing from queue.\n", job.JobID)
            job.Processed = true // Lo marchiamo processato per non vederlo più
            JobCache.queue = JobCache.queue[1:]
            continue
        }

		// Active job found. Pack the RequestID and IPFS CID into an ABI-encoded query.
        query, err := QueryArgs.Pack(job.JobID, job.CID)
		if err != nil {
        return nil, fmt.Errorf("Query Pack failed: %w", err)
    	}

        fmt.Printf("QUERY (Leader): Proposing Job #%s (CID: %s)\n", job.JobID, job.CID)
        return query, err
    }

    return nil, nil  // Empty queue
}

// =============================================================================
// PHASE 2: OBSERVATION (Nodes fetch data)
// =============================================================================

// Observation is executed by all network nodes upon receiving the Leader's Query.
// It retrieves the AI computation results from the local cache. If the asynchronous
// worker is still processing, it gracefully yields an empty observation.
func (p *attributionPlugin) Observation(ctx context.Context, _ ocr3types.OutcomeContext, query ocrtypes.Query) (ocrtypes.Observation, error) {
	start := time.Now()
	defer func() {
        fmt.Printf("[METRIC-OCR] Phase: OBSERVATION | Node: %d | Time: %v\n", p.cfg.OracleID, time.Since(start))
    }()

	// If the leader proposed an empty query, return empty observation
	if len(query) == 0 { return ocrtypes.Observation("{}"), nil }
	
	qVals, err := QueryArgs.Unpack(query)
	if err != nil { return ocrtypes.Observation("{}"), nil }
	
	reqID := qVals[0].(*big.Int)
	jobId64 := reqID.Uint64()

	// Thread-safe read from the shared cache
	JobCache.RLock()
	job, exists := JobCache.jobs[jobId64]
	if !exists {
		JobCache.RUnlock()
		return ocrtypes.Observation("{}"), nil
	}
	// Copy the variables fields while you still have the lock
	state := job.State
	result := job.Result
	JobCache.RUnlock()

	// Case 1: Asynchronous AI computation is in progress.
    // Returning an empty array prevents triggering libocr's MaxDurationObservation timeouts.
	if state == StatePending {
		return ocrtypes.Observation("{}"), nil
	}

	// CASE 2: The background worker encountered an error
	if state == StateFailed {
		fmt.Printf("OBSERVATION Oracle=%d: Job %d Failed: %v\n", p.cfg.OracleID, jobId64, job.Err)
		return ocrtypes.Observation("{}"), nil
	}

	// Case 3: Byzantine Fault Tolerance (BFT) Attack Simulation.
    // Alters the generated attribution vector to test the resilience of the Outcome phase.
	maliciousMode := os.Getenv("MALICIOUS_MODE")

	if maliciousMode == "alter" {
		fmt.Printf("[ALERT] Oracle=%d: MALICIOUS MODE ACTIVE. Altering vector by +5%%.\n", p.cfg.OracleID)
        
        // Create a fake vector with the same length of the original
        fakeVector := make([]*big.Int, len(result))
        
        for i, originalVal := range result {
            fakeVal := new(big.Int).Set(originalVal) 
            
			// Introduce a +5% deviation: (val * 105) / 100
            fakeVal.Mul(fakeVal, big.NewInt(105))
            fakeVal.Div(fakeVal, big.NewInt(100))
            
            fakeVector[i] = fakeVal
        }
		fmt.Printf("OBSERVATION Malicious Oracle=%d: Generated Vector len=%d. First: %s\n", p.cfg.OracleID, len(fakeVector), fakeVector[0].String())    
        
        packedObs, err := ObservationArgs.Pack(fakeVector)
        if err != nil { return nil, err }
        return ocrtypes.Observation(packedObs), nil
	}

	// Malicious node (timeout) or timeout/break of node
	if maliciousMode == "timeout" {
		fmt.Printf("[ALERT] Oracle=%d: MALICIOUS MODE (TIMEOUT) ACTIVE. Simulating crash/offline status.\n", p.cfg.OracleID)
		// simulates a node that takes infinite time to respond
		time.Sleep(1*time.Hour)
	}

	// Case 4: Honest Execution. 
    // Pack the ABI-encoded true attribution vector into the observation.
	fmt.Printf("OBSERVATION Oracle=%d: Generated Vector len=%d. First: %s\n", p.cfg.OracleID, len(result), result[0].String())    
	packedObs, err := ObservationArgs.Pack(result)
	if err != nil { return nil, err }

	return ocrtypes.Observation(packedObs), nil
}

// ValidateObservation is called on every received observation to reject malformed
// data early.
// Returning nil here means "accept everything".
func (p *attributionPlugin) ValidateObservation(context.Context, ocr3types.OutcomeContext, ocrtypes.Query, ocrtypes.AttributedObservation) error { return nil }


// =============================================================================
// PHASE 3: CONSENSUS & OUTCOME (Leader aggregates data)
// =============================================================================

// ObservationQuorum decides when the leader has enough valid observations to
// propose an outcome. The default OCR2-style rule is "2f+1 observations".
//
// With n=4, f=1 this means 3 valid observations are required to proceed.
func (p *attributionPlugin) ObservationQuorum(_ context.Context, _ ocr3types.OutcomeContext, _ ocrtypes.Query, aos []ocrtypes.AttributedObservation) (bool, error) {
	return quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumTwoFPlusOne, p.cfg.N, p.cfg.F, aos), nil
}

// Outcome acts as the final consensus gatekeeper.
// It deterministically aggregates observations using a median calculation 
// to nullify the impact of Byzantine (malicious) vectors.
func (p *attributionPlugin) Outcome(ctx context.Context, _ ocr3types.OutcomeContext, query ocrtypes.Query, attrObservation []ocrtypes.AttributedObservation) (ocr3types.Outcome, error) {
	start := time.Now()
	defer func() {
        fmt.Printf("[METRIC-OCR] Phase: OUTCOME (Median Calc) | Leader Node | Time: %v\n", time.Since(start))
    }()

	if len(query) == 0 { 
		return ocr3types.Outcome([]byte{}), nil 
	}

	// Important if we have two identical vectors (very unprobable)
	sort.Slice(attrObservation, func(i, j int) bool {
		return attrObservation[i].Observer < attrObservation[j].Observer
	})
	
	// Unpack Query to get the ID
	qVals, err := QueryArgs.Unpack(query)
	if err != nil { return ocr3types.Outcome([]byte{}), nil }
	reqID := qVals[0].(*big.Int)

	// Decode and store all the received vectors
	// by casting as ObservationArgs the elements of aos (array of observations)
	var candidates [][]*big.Int

	for _, ao := range attrObservation {
        // Unpack Observation ABI
		vals, err := ObservationArgs.Unpack(ao.Observation)
		if err != nil { continue }
        
        // Extract the vector (index 0)
		vector := vals[0].([]*big.Int)
		if len(vector) > 0 {
			candidates = append(candidates, vector)
		}
	}

	// Enforce strict BFT quorum bounds before computing the median
	if len(candidates) < (2*p.cfg.F + 1) {
		fmt.Println("WARN: Not enough observations")
		return ocr3types.Outcome([]byte{}), nil
	}

	// ==================================================================
	// DETERMINISTIC MEDIAN COMPUTATION
	// ==================================================================

	// 1. Take the length of the first vector as reference
	vectorLen := len(candidates[0])
	winner := make([]*big.Int, vectorLen)

	// 2. For every position in the vector
	for i := 0; i < vectorLen; i++ {
		var valuesAtPosition []*big.Int

		// Take all the values for that position
		for _, vec := range candidates {
			// Security check: if an oracle send a vector less long, ignore it for this index
			if len(vec) > i {
				valuesAtPosition = append(valuesAtPosition, vec[i])
			}
		}
		// 3. Order the values from the smallest to the biggest (useful for the median)
		sort.Slice(valuesAtPosition, func(a, b int) bool {
			return valuesAtPosition[a].Cmp(valuesAtPosition[b]) < 0
		})

		// 4. Take the central value
		medianIndex := len(valuesAtPosition)/2

		// Assign this value to the winner
		if len(valuesAtPosition) > 0 {
			winner[i] = valuesAtPosition[medianIndex]
		} else {
			// Fallback if the list is empty (must not happen)
			winner[i] = big.NewInt(0)
		}
	}

	fmt.Printf("OUTCOME (Leader): Selected Vector len %d\n", len(winner))

    // PACK ABI OUTCOME
	packedOutcome, err := OutcomeArgs.Pack(reqID, winner)
	if err != nil { return nil, err }
	return ocr3types.Outcome(packedOutcome), nil
}

// =============================================================================
// PHASE 4: REPORT GENERATION
// =============================================================================

// Reports converts the agreed Outcome into a report to be signed by off-chain keys.
func (p *attributionPlugin) Reports(_ context.Context, _ uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportPlus[struct{}], error) {
	return []ocr3types.ReportPlus[struct{}]{{ReportWithInfo: ocr3types.ReportWithInfo[struct{}]{Report: ocrtypes.Report(outcome), Info: struct{}{}}}}, nil
}

// ShouldAcceptAttestedReport is called after a report has collected enough
// onchain signatures to be considered "attested" (typically f+1 signers).
//
// Returning true means "this report is acceptable for potential transmission".
func (p *attributionPlugin) ShouldAcceptAttestedReport(context.Context, uint64, ocr3types.ReportWithInfo[struct{}]) (bool, error) { return true, nil }

// ShouldTransmitAcceptedReport is called right before transmission.
// Returning true means "actually transmit now (subject to the schedule)".
func (p *attributionPlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, r ocr3types.ReportWithInfo[struct{}]) (bool, error) {
	// Start stopwatch 
	// One of the last phases before the transmission to the blockchain, useful to measure it
	start := time.Now()
    defer func() {
        fmt.Printf("[METRIC-OCR] Phase: TRANSMIT_CHECK | Node: %d | Time: %v\n", p.cfg.OracleID, time.Since(start))
    }()
	
	// Baseline size verification for the ABI payload
	if len(r.Report) < 32 { return false, nil }

    // Unpack report
	values, err := OutcomeArgs.Unpack(r.Report)
	if err != nil { return false, nil }
	
	// Sanity check on the associated Request ID
    requestID := values[0].(*big.Int)
	if requestID == nil { return false, nil }

	return true, nil
}
func (p *attributionPlugin) Close() error { return nil }
