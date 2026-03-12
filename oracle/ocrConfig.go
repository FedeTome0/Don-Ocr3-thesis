package main

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// =============================================================================
// OCR SETUP CONFIGURATION
// =============================================================================

// buildContractConfig generates a deterministic OCR3 contract config (including
// onchain/offchain config blobs) for this experiment.
//
// In real OCR3, this config would be set on-chain via `setConfig(...)` and the
// oracles would learn it through ContractConfigTracker + their DB.
func buildContractConfig(n, f int, seed int64) (ocrtypes.ContractConfig, evmutil.EVMOffchainConfigDigester, error) {
	_, oracles, err := deriveAllNodes(n, seed)
	if err != nil {
		return ocrtypes.ContractConfig{}, evmutil.EVMOffchainConfigDigester{}, err
	}

	// IMPORTANT: Scheduling transmission
	//transmissionSchedule := make([]int, n)
	//for i := 0; i < n; i++ {
	//	transmissionSchedule[i] = 1
	//}

	// With 7 nodes 
	transmissionSchedule := []int{1,2,4}

	// Deterministic configuration of OCR3 standards
	// These values control the speed of the consensus rounds.
	signers, transmitters, fOut, onchain, ver, offchain, err := ocr3confighelper.ContractSetConfigArgsDeterministic(
		sha256.Sum256([]byte(fmt.Sprintf("sk|%d", seed))),               //ephemeralSk
		[16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}, //sharedSecret
		40*time.Second, // deltaProgress: Max time for a round leader to drive progress
		20*time.Second, // deltaResend: Time before resending messages
		20*time.Second, // deltaInitial: Initial round timeout
		10*time.Second, // deltaRound: Timeout for a single round
		2*time.Second,  // deltaGrace: Grace period for slow nodes
		2*time.Second,  // deltaCertifiedCommitRequest
		10*time.Second, // deltaStage: Time per stage (Prepare/Commit/Accept)
		5,              // rMax: Max rounds
		transmissionSchedule,   // Schedule (S) FIX
		oracles,        // List of oracle identities
		[]byte("{}"),   // Reporting Plugin Config (Empty for now)
		nil,
		20*time.Second, // MaxDurQuery
		30*time.Second, // MaxDurObs (Ideally kept short, hence our async architecture)
		20*time.Second, // MaxDurSat
		90*time.Second, // MaxDurST
		f,              // Fault tolerance threshold (f)
		nil,
	)
	if err != nil {
		return ocrtypes.ContractConfig{}, evmutil.EVMOffchainConfigDigester{}, err
	}

	// The digester determines how configDigest is computed. We use an EVM digester
	// with a dummy (but fixed) chainID + contract address so that all nodes agree.
	addr := common.HexToAddress("0x" + strings.Repeat("11", 20))
	digester := evmutil.EVMOffchainConfigDigester{ChainID: 1337, ContractAddress: addr}
	cc := ocrtypes.ContractConfig{
		ConfigCount:           1,
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     fOut,
		OnchainConfig:         onchain,
		OffchainConfigVersion: ver,
		OffchainConfig:        offchain,
	}
	cc.ConfigDigest, err = digester.ConfigDigest(context.Background(), cc)
	if err != nil {
		return ocrtypes.ContractConfig{}, evmutil.EVMOffchainConfigDigester{}, err
	}
	return cc, digester, nil
}

// deriveBootstrapPriv deterministically derives the bootstrapper's Ed25519 key.
// This makes the bootstrap peerID stable across runs (handy for docker-compose).
func deriveBootstrapPriv(seed int64) ed25519.PrivateKey {
	h := sha256.Sum256([]byte(fmt.Sprintf("bootstrap|%d", seed)))
	return ed25519.NewKeyFromSeed(h[:])
}