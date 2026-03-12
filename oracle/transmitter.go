package main

import (
	"context"
	"crypto/ecdsa"
	//"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc" 
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"OCR3-thesis/contracts"
)
/*
// ethTransmitter handles the on-chain submission of the OCR3 consensus outcome.
// It maintains persistent RPC connections and cryptographic materials to ensure
// low-latency and gas-efficient transaction broadcasting.
type ethTransmitter struct {
	oracleID		int
	nOracles 		int
	contractAddr	common.Address
	jobCache		*jobCache // pointer to the shared thread-safe cache
	// Pre-initialized resources to avoid overhead during the transmission phase	
	rpcClient  *rpc.Client
	ethClient  *ethclient.Client
	verifier   *contracts.OracleVerifier
	parsedABI  abi.ABI
	privateKey *ecdsa.PrivateKey
	address    common.Address
}

// NewEthTransmitter initializes the Ethereum transmitter instance.
// It derives the public identity, establishes persistent RPC connections,
// and prepares the ABI structures required for contract interactions.
func NewEthTransmitter(id, n int, rpcUrl, privKeyHex string, addr common.Address, cache *jobCache) (*ethTransmitter, error) {
	// 1. Cryptographic Identity Setup
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}
	pubKey := privKey.Public().(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*pubKey)

	// 2. Persistent RPC connetions, initialized here
	rpcClient, err := rpc.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial raw RPC: %w", err)
	}

	ethClient, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial eth client: %w", err)
	}

	// 3. Smart Contract Binding & ABI Parsing
	verifier, err := contracts.NewOracleVerifier(addr, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to bind contract: %w", err)
	}

	// Pre-parsing the ABI for the view function to save CPU cycles during execution
	const isCompletedABI = `[{"inputs":[{"internalType":"uint256","name":"_requestId","type":"uint256"}],"name":"isCompleted","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]`
	parsedABI, _ := abi.JSON(strings.NewReader(isCompletedABI))

	return &ethTransmitter{
		oracleID:     id,
		nOracles:     n,
		contractAddr: addr,
		jobCache:     cache,
		rpcClient:    rpcClient,
		ethClient:    ethClient,
		verifier:     verifier,
		parsedABI:    parsedABI,
		privateKey:   privKey,
		address:      address,
	}, nil
}

// FromAccount retrieves the Ethereum identity of the oracle node.
// LibOCR uses this method to verify which on-chain account corresponds
// to this transmitter during the signature verification phase.
func (t *ethTransmitter) FromAccount(context.Context) (ocrtypes.Account, error) {
    // Return the pre-computed address instead of calculating it on-the-fly
    return ocrtypes.Account(t.address.Hex()), nil
}


// Transmit broadcasts the consensus outcome to the Ethereum blockchain.
// It includes multi-layered duplicate prevention (Cache + On-chain View) 
// to ensure gas efficiency.
func (t *ethTransmitter) Transmit(ctx context.Context, cd ocrtypes.ConfigDigest, seq uint64, report ocr3types.ReportWithInfo[struct{}], sigs []ocrtypes.AttributedOnchainSignature) error {

	// Start the stopwatch of the total transmission
	// In this section are used more clocks: one for all the transmission and few for granular analysis
	startTotal := time.Now()
    defer func() {
        fmt.Printf("[METRIC-OCR] Phase: TRANSMIT_TOTAL | Oracle: %d | Time: %v\n", t.oracleID, time.Since(startTotal))
    }()

	// =========================================================================
    // 1. REPORT DECODING
    // Unpack the ABI-encoded report to extract the RequestID and Result Vector
    // =========================================================================
	if len(report.Report) < 32 {
		return nil
	}

	values, err := OutcomeArgs.Unpack(report.Report)
	if err != nil {
		return fmt.Errorf("failed to unpack report ABI: %v", err)
	}

	requestID := values[0].(*big.Int)
	vector := values[1].([]*big.Int)
	jobId64 := requestID.Uint64()

	// =========================================================================
    // 2. LOCAL CACHE CHECK
    // Check the thread-safe memory cache before executing any network calls.
    // If the asynchronous listener already marked this job as processed, skip.
    // =========================================================================

	t.jobCache.RLock()
	job, exists := t.jobCache.jobs[jobId64]
	isProcessedLocal := exists && job.Processed
	t.jobCache.RUnlock()

	if isProcessedLocal {
		return nil
	}

	// =========================================================================
    // 3. ON-CHAIN VIEW CHECK
    // Synchronous double-check directly against the blockchain state.
    // This acts as a fallback to prevent race conditions if the WebSocket
    // listener experiences latency.
    // Note: We use raw JSON-RPC (eth_call) to bypass known Hardhat view revert bugs.
    // =========================================================================
	inputData, err := t.parsedABI.Pack("isCompleted", requestID)
	if err == nil {
		callMsg := map[string]string{
			"to":   t.contractAddr.Hex(),
			"data": hexutil.Encode(inputData),
		}

		var resultHex string

		startRPC := time.Now()
		err = t.rpcClient.CallContext(ctx, &resultHex, "eth_call", callMsg, "latest")
		rpcDuration := time.Since(startRPC)

		// Print the evaluation only if the call doesn't return an error
		if err == nil {
			fmt.Printf("[METRIC-OCR] Phase: RPC_VIEW_CALL | Oracle: %d | Time: %v\n", t.oracleID, rpcDuration)

			// Decode the raw hex response back into a Go boolean
			resultBytes, decodeErr := hexutil.Decode(resultHex)
			if decodeErr == nil {
				var isDone bool
				t.parsedABI.UnpackIntoInterface(&isDone, "isCompleted", resultBytes)

				if isDone {
					// Job already finalized on-chain, no need to submit again.
					fmt.Printf("[Oracle %d] Job #%s already completed on-chain. Skipping transmission.\n", t.oracleID, requestID)
					MarkJobAsProcessed(jobId64)
					return nil
				}
			}
		} else {
			fmt.Printf("[Oracle %d] View check failed: %v.\n", t.oracleID, err)
		}
	}

	// =========================================================================
    // 4. TRANSACTION BUILDING & BROADCASTING
    // =========================================================================
	// Add of transmit fail logic to show the transmission schedule working
	if os.Getenv("MALICIOUS_MODE") == "transmit_fail" {
		fmt.Printf("[ALERT] Oracle=%d: MALICIOUS MODE (TRANSMIT_FAIL) ACTIVE. Dropping on-chain transaction.\n", t.oracleID)
		// Return an error to say to the ocr protcol that the transmission failed
		return fmt.Errorf("Simulated RPC network failure during Transmit")
	}
	// ------------- End transmit fail logic ----------------------
	
	fmt.Printf("[Oracle %d] Transmitting Job #%s... \n", t.oracleID, requestID)
	
	chainID, _ := t.ethClient.NetworkID(ctx)
	auth, err := bind.NewKeyedTransactorWithChainID(t.privateKey, chainID)
	if err != nil { return err }

	// Dynamic gas estimation
	auth.GasPrice, _ = t.ethClient.SuggestGasPrice(ctx)
	auth.GasFeeCap = nil
	auth.GasTipCap = nil
	auth.GasLimit = 30000000

	startTx := time.Now()
	tx, err := t.verifier.SaveOutcome(auth, requestID, vector)
	
	if err != nil {
		// Contract-level duplicate detection (Revert on "already fulfilled")
		if strings.Contains(err.Error(), "Request already fulfilled") || strings.Contains(err.Error(), "revert") {
			fmt.Printf("[Oracle %d] Race condition intercettata dal contratto (Revert).\n", t.oracleID)
			MarkJobAsProcessed(jobId64)
			return nil
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	fmt.Printf("[METRIC-OCR] Phase: TX_SUBMIT | Oracle: %d | Time: %v | Hash: %s\n", t.oracleID, time.Since(startTx), tx.Hash().Hex())
	// Optimistically update the local cache
	MarkJobAsProcessed(jobId64)

	return nil
}

// Close gracefully shuts down the persistent network connections.
func (t *ethTransmitter) Close() error {
    t.rpcClient.Close()
    t.ethClient.Close()
    return nil
}
*/

// ethTransmitter handles the on-chain submission of the OCR3 consensus outcome.
// It maintains persistent RPC connections and cryptographic materials to ensure
// low-latency and gas-efficient transaction broadcasting.
type ethTransmitter struct {
    oracleID        int
    nOracles        int
    contractAddr    common.Address
    jobCache        *jobCache // pointer to the shared thread-safe cache
    // Pre-initialized resources to avoid overhead during the transmission phase    
    rpcClient  *rpc.Client
    ethClient  *ethclient.Client
    verifier   *contracts.OracleVerifier
    parsedABI  abi.ABI
    privateKey *ecdsa.PrivateKey
    address    common.Address
}

// NewEthTransmitter initializes the Ethereum transmitter instance.
// It derives the public identity, establishes persistent RPC connections,
// and prepares the ABI structures required for contract interactions.
func NewEthTransmitter(id, n int, rpcUrl, privKeyHex string, addr common.Address, cache *jobCache) (*ethTransmitter, error) {
    // 1. Cryptographic Identity Setup
    privKey, err := crypto.HexToECDSA(privKeyHex)
    if err != nil {
        return nil, fmt.Errorf("invalid private key: %w", err)
    }
    pubKey := privKey.Public().(*ecdsa.PublicKey)
    address := crypto.PubkeyToAddress(*pubKey)

    // 2. Persistent RPC connetions, initialized here
    rpcClient, err := rpc.Dial(rpcUrl)
    if err != nil {
        return nil, fmt.Errorf("failed to dial raw RPC: %w", err)
    }

    ethClient, err := ethclient.Dial(rpcUrl)
    if err != nil {
        return nil, fmt.Errorf("failed to dial eth client: %w", err)
    }

    // 3. Smart Contract Binding & ABI Parsing
    verifier, err := contracts.NewOracleVerifier(addr, ethClient)
    if err != nil {
        return nil, fmt.Errorf("failed to bind contract: %w", err)
    }

    // Pre-parsing the ABI for the view function to save CPU cycles during execution
    const isCompletedABI = `[{"inputs":[{"internalType":"uint256","name":"_requestId","type":"uint256"}],"name":"isCompleted","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]`
    parsedABI, _ := abi.JSON(strings.NewReader(isCompletedABI))

    return &ethTransmitter{
        oracleID:     id,
        nOracles:     n,
        contractAddr: addr,
        jobCache:     cache,
        rpcClient:    rpcClient,
        ethClient:    ethClient,
        verifier:     verifier,
        parsedABI:    parsedABI,
        privateKey:   privKey,
        address:      address,
    }, nil
}

// FromAccount retrieves the Ethereum identity of the oracle node.
// LibOCR uses this method to verify which on-chain account corresponds
// to this transmitter during the signature verification phase.
func (t *ethTransmitter) FromAccount(context.Context) (ocrtypes.Account, error) {
    // Return the pre-computed address instead of calculating it on-the-fly
    return ocrtypes.Account(t.address.Hex()), nil
}

//function transmit(bytes32 configDigest, uint64 seqNr, bytes report, bytes attestation)

// Transmit broadcasts the consensus outcome to the Ethereum blockchain.
// La firma rispetta al 100% l'interfaccia ocr3types.ContractTransmitter.
func (t *ethTransmitter) Transmit(
    ctx context.Context, 
	digest ocrtypes.ConfigDigest,
	seqNr uint64,
    report ocr3types.ReportWithInfo[struct{}], 
    sigs []ocrtypes.AttributedOnchainSignature,
) error {
	
    // Start the stopwatch of the total transmission
    startTotal := time.Now()
    defer func() {
        fmt.Printf("[METRIC-OCR] Phase: TRANSMIT_TOTAL | Oracle: %d | Time: %v\n", t.oracleID, time.Since(startTotal))
    }()

    // =========================================================================
    // 1. REPORT DECODING
    // Unpack the ABI-encoded report to extract the RequestID for caching
    // =========================================================================
    if len(report.Report) < 32 {
        return nil
    }

    values, err := OutcomeArgs.Unpack(report.Report)
    if err != nil {
        return fmt.Errorf("failed to unpack report ABI: %v", err)
    }

    requestID := values[0].(*big.Int)
    jobId64 := requestID.Uint64()

    // =========================================================================
    // 2. LOCAL CACHE CHECK
    // =========================================================================
    t.jobCache.RLock()
    job, exists := t.jobCache.jobs[jobId64]
    isProcessedLocal := exists && job.Processed
    t.jobCache.RUnlock()

    if isProcessedLocal {
        return nil
    }

    // =========================================================================
    // 3. ON-CHAIN VIEW CHECK
    // Synchronous double-check directly against the blockchain state.
    // =========================================================================
    inputData, err := t.parsedABI.Pack("isCompleted", requestID)
    if err == nil {
        callMsg := map[string]string{
            "to":   t.contractAddr.Hex(),
            "data": hexutil.Encode(inputData),
        }

        var resultHex string
        startRPC := time.Now()
        err = t.rpcClient.CallContext(ctx, &resultHex, "eth_call", callMsg, "latest")
        rpcDuration := time.Since(startRPC)

        if err == nil {
            fmt.Printf("[METRIC-OCR] Phase: RPC_VIEW_CALL | Oracle: %d | Time: %v\n", t.oracleID, rpcDuration)

            resultBytes, decodeErr := hexutil.Decode(resultHex)
            if decodeErr == nil {
                var isDone bool
                t.parsedABI.UnpackIntoInterface(&isDone, "isCompleted", resultBytes)

                if isDone {
                    fmt.Printf("[Oracle %d] Job #%s already completed on-chain. Skipping transmission.\n", t.oracleID, requestID)
                    MarkJobAsProcessed(jobId64)
                    return nil
                }
            }
        } else {
            fmt.Printf("[Oracle %d] View check failed: %v.\n", t.oracleID, err)
        }
    }

    // =========================================================================
    // 4. TRANSACTION BUILDING & BROADCASTING
    // =========================================================================
    if os.Getenv("MALICIOUS_MODE") == "transmit_fail" {
        fmt.Printf("[ALERT] Oracle=%d: MALICIOUS MODE (TRANSMIT_FAIL) ACTIVE. Dropping on-chain transaction.\n", t.oracleID)
        return fmt.Errorf("Simulated RPC network failure during Transmit")
    }
    
    fmt.Printf("[Oracle %d] Transmitting Job #%s with ECDSA Signatures... \n", t.oracleID, requestID)
    
    chainID, _ := t.ethClient.NetworkID(ctx)
    auth, err := bind.NewKeyedTransactorWithChainID(t.privateKey, chainID)
    if err != nil { return err }

    auth.GasPrice, _ = t.ethClient.SuggestGasPrice(ctx)
    auth.GasFeeCap = nil
    auth.GasTipCap = nil
    auth.GasLimit = 30000000 // High gas limit to account for ecrecover loop
    
    // --- CRYTOGRAPHIC SIGNATURE UNPACKING FOR EVM ---
    var rs [][32]byte
    var ss [][32]byte
    var rawVs [32]byte // Fixed 32-byte array to compactly store up to 31 'v' values

    // Inizializza un contatore fuori dal ciclo
    validSigCount := 0

    for _, sig := range sigs {
        if len(sig.Signature) != 65 { 
            return fmt.Errorf("invalid signature length") 
        }
        
        var r, s [32]byte
        copy(r[:], sig.Signature[0:32])
        copy(s[:], sig.Signature[32:64])
        
        v := sig.Signature[64]
        if v < 27 { 
            v += 27 // EVM normalization (0/1 -> 27/28)
        }
        
        rs = append(rs, r)
        ss = append(ss, s)
        
        // FIX: Riempiamo rawVs in ordine sequenziale, 
        // esattamente come rs e ss.
        if validSigCount < len(rawVs) {
            rawVs[validSigCount] = v 
            validSigCount++
        }
    }

    startTx := time.Now()
    
    // Execute the contract method using the generated abigen bindings
    tx, err := t.verifier.Transmit(auth, digest, seqNr, report.Report, rs, ss, rawVs)
    
    if err != nil {
        if strings.Contains(err.Error(), "already fulfilled") || strings.Contains(err.Error(), "revert") {
            fmt.Printf("[Oracle %d] Race condition intercettata dal contratto (Revert).\n", t.oracleID)
            MarkJobAsProcessed(jobId64)
            return nil
        }
        return fmt.Errorf("transaction failed: %w", err)
    }

    fmt.Printf("[METRIC-OCR] Phase: TX_SUBMIT | Oracle: %d | Time: %v | Hash: %s\n", t.oracleID, time.Since(startTx), tx.Hash().Hex())
    MarkJobAsProcessed(jobId64)
    return nil
}