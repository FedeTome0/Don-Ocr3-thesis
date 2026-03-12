package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/libocr/commontypes"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"


)

// =============================================================================
// JSON PAYKIAD STRUCTURES 
// =============================================================================

// AttributeRequest represents the JSON payload sent via POST to the Python AI service.
type AttributeRequest struct {
    JobId  string `json:"job_id"`  // Unique identifier for the attribution job
    Text string `json:"text"` // Content downloaded from IPFS
}

// AttributeResponse represents the expected JSON response from the Python AI polling endpoint.
type AttributeResponse struct {
    Status string   `json:"status"` // Job lifecycle state (e.g., "processing", "completed")
    Result []string `json:"result"` // String representations of the calculated BigInt vector
    Error  string   `json:"error,omitempty"` // Populated only if the AI computation fails
}

// =============================================================================
// ON-CHAIN ABI SCHEMA DEFINITIONS
// =============================================================================

// The following structures define the encoding rules for OCR3 consensus payloads.
// Strict ABI encoding ensures native data compatibility with the EVM Smart Contract.
var (
	// QueryArgs: Data proposed by the Round Leader to initialize consensus.
	// Query: uint256 (JobId), string (IPFS CID)
	QueryArgs = abi.Arguments{
		{Type: MustParseType("uint256")},
		{Type: MustParseType("string")},
	}

	// ObservationArgs: The calculated data returned by individual Oracle nodes.
    // Encodes: int128[] (Attribution Vector)
    // Note: Optimized from int256[] to int128[] to leverage Solidity storage packing,
    // significantly reducing gas costs during the final state transition.
	ObservationArgs = abi.Arguments{
		{Type: MustParseType("int128[]")},
	}

	// OutcomeArgs: The final, aggregated BFT consensus result.
    // Encodes: uint256 (JobID), int128[] (Median Attribution Vector)
	OutcomeArgs = abi.Arguments{
		{Type: MustParseType("uint256")},
		{Type: MustParseType("int128[]")},
	}
)

// MustParseType dynamically parses a Solidity type string into a valid Go-Ethereum ABI type.
// It panics upon initialization if the type is invalid, enforcing strict type safety.
func MustParseType(t string) abi.Type {
	ty, err := abi.NewType(t, "", nil)
	if err != nil {
		panic(fmt.Sprintf("Errore tipo ABI %s: %v", t, err))
	}
	return ty
}


// =============================================================================
// INFRASTRUCTURE STUBS
// =============================================================================

// quietLogger implements libocr's logger interface but intentionally drops most
// log levels to keep the console output focused on the experiment's fmt.Printf
// lines (QUERY/OBS/TRANSMIT).

type quietLogger struct{ *log.Logger }
func (l quietLogger) Trace(string, commontypes.LogFields)          {}
func (l quietLogger) Debug(string, commontypes.LogFields)          {}
func (l quietLogger) Info(string, commontypes.LogFields)           {}
func (l quietLogger) Warn(string, commontypes.LogFields)           {}
func (l quietLogger) Error(msg string, f commontypes.LogFields)    { l.Println("ERROR", msg, f) }
func (l quietLogger) Critical(msg string, f commontypes.LogFields) { l.Println("CRIT", msg, f) }

type noopMonitoring struct{}
// noopMonitoring is a stub MonitoringEndpoint. In production you would export
// libocr telemetry/logs to your monitoring system.
func (noopMonitoring) SendLog([]byte) {}

type memDB3 struct {
	mu    sync.Mutex
	cfg   *ocrtypes.ContractConfig
	state map[string]map[string][]byte
}

// ReadConfig returns the latest contract config the protocol should run with.
// In a real deployment this would come from chain (via ContractConfigTracker)
// and be persisted in a durable DB.
func (m *memDB3) ReadConfig(context.Context) (*ocrtypes.ContractConfig, error) { 
	return m.cfg, nil 
}

// WriteConfig stores the current contract config.
func (m *memDB3) WriteConfig(_ context.Context, c ocrtypes.ContractConfig) error { 
	m.cfg = &c; return nil 
}

// ReadProtocolState returns protocol-internal persisted state (per configDigest).
// libocr uses this to survive restarts and avoid safety bugs.
func (m *memDB3) ReadProtocolState(_ context.Context, d ocrtypes.ConfigDigest, key string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.state == nil {
		return nil, nil
	}
	kv := m.state[d.Hex()]
	if kv == nil {
		return nil, nil
	}
	v := kv[key]
	if v == nil {
		return nil, nil
	}
	out := make([]byte, len(v))
	copy(out, v)
	return out, nil
}

// WriteProtocolState stores protocol-internal persisted state (per configDigest).
func (m *memDB3) WriteProtocolState(_ context.Context, d ocrtypes.ConfigDigest, key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.state == nil {
		m.state = map[string]map[string][]byte{}
	}
	if m.state[d.Hex()] == nil {
		m.state[d.Hex()] = map[string][]byte{}
	}
	if value == nil {
		delete(m.state[d.Hex()], key)
		return nil
	}
	v := make([]byte, len(value))
	copy(v, value)
	m.state[d.Hex()][key] = v
	return nil
}

type staticTracker struct{ 
	cfg ocrtypes.ContractConfig 
}

// staticTracker is a tiny ContractConfigTracker implementation that always
// returns the single config we computed in buildContractConfig().
//
// Real OCR runs with a ContractConfigTracker that watches chain logs / RPC for
// config changes.
func (t staticTracker) Notify() <-chan struct{} { return nil }
func (t staticTracker) LatestConfigDetails(context.Context) (uint64, ocrtypes.ConfigDigest, error) { return 1, t.cfg.ConfigDigest, nil }
func (t staticTracker) LatestConfig(context.Context, uint64) (ocrtypes.ContractConfig, error) { return t.cfg, nil }
func (t staticTracker) LatestBlockHeight(context.Context) (uint64, error) { return 1, nil }

// =============================================================================
// ASYNCHRONOUS STATE MANAGEMENT & LOCAL CACHE
// =============================================================================

// JobState tracks the lifecycle of an off-chain computation.
type JobState int
const (
	StatePending JobState = iota	// Background AI computation is in progress
	StateCompleted					// Result vector is ready in local RAM
	StateFailed						// Async task encountered a fatal error
)

// JobData encapsulates all context required to process a specific request
type JobData struct {
	JobID     *big.Int
	CID       string
	State     JobState
	Result    []*big.Int
	Err       error
	Processed bool // If true, it has been sent on chain
}

// Reduce the requestId to uint64 because it is only an
// incremental integer and is not going to be so large
// and it is useful for the map structure (it can't handle uint128 o uint256)
type jobCache struct {
	sync.RWMutex
	LatestJobID *big.Int
	jobs        map[uint64]*JobData
	queue 		[]uint64 // FIFO array to maintain determinstic processing order for the Leader
}

// Global instance initialized at startup
var JobCache = &jobCache{
	LatestJobID: big.NewInt(-1),
	jobs:        make(map[uint64]*JobData),
	queue:		 make([]uint64, 0),
}

// enqueue safely registers a newly detected job, appending it to the processing queue.
func (c *jobCache) enqueue(id uint64, job *JobData) {
	if _, exists := c.jobs[id]; !exists {
		c.jobs[id] = job
		c.queue = append(c.queue, id)
	}
}

// MarkJobAsProcessed updates the job status atomically
// Used by the Transmitter to signal that a transaction has been successfully sent
func MarkJobAsProcessed(jobIdUint uint64) {
	JobCache.Lock()
	defer JobCache.Unlock()
	if job, exists := JobCache.jobs[jobIdUint]; exists {
		job.Processed = true
	}
}