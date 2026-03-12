package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/networking"
	"github.com/smartcontractkit/libocr/offchainreporting2plus"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"
)

// =============================================================================
// RUNNERS (Node Bootstrapping & P2P Networking)
// =============================================================================

// runBootstrap initializes the Bootstrap Node.
// This node does not participate in the OCR3 consensus (it doesn't propose or sign data).
// It handles the Oracles discovery of peer by communicating with them and giving the other peer informations
func runBootstrap(ctx context.Context, n, f int, seed int64, listen, announce string) error {
	cc, _, err := buildContractConfig(n, f, seed)
	if err != nil { return err }
	
	_, oracles, err := deriveAllNodes(n, seed)
	if err != nil {	return err }
	
	// Collect the Peer IDs of all active oracles to authorize them on the network
	peerIDs := make([]string, 0, n)
	for _, o := range oracles { 
		peerIDs = append(peerIDs, o.PeerID) 
	}
	
	// DNS Resolution: crucial for Docker container networking mapping
	announceIP, err := resolveHostnamePortToIP(announce)
	if err != nil {
		return fmt.Errorf("Resolve bootstrap_announce %q: %w", announce, err)
	}

	// Deterministically derive the Bootstrapper's cryptographic identity
	priv := deriveBootstrapPriv(seed)
	pid, err := ragetypes.PeerIDFromPrivateKey(priv)
	if err != nil {
		return err
	}

	// Initialize the networking layer (libp2p)
	peer, err := networking.NewPeer(networking.PeerConfig{
		PrivKey: 				priv, 
		Logger: 				quietLogger{log.New(os.Stdout, "", 0)}, 
		V2ListenAddresses: 		[]string{listen}, 
		V2AnnounceAddresses: 	[]string{announceIP},
		V2DeltaReconcile:		10 * time.Second,
		V2EndpointConfig: networking.EndpointConfigV2{
			IncomingMessageBufferSize: 100,
			OutgoingMessageBufferSize: 50,
		},
		MetricsRegisterer:		prometheus.NewRegistry(),
	})
	if err != nil { return err }

	bootstrapper, err := peer.OCR2BootstrapperFactory().NewBootstrapper(cc.ConfigDigest, peerIDs, nil, f)
	if err != nil {
		return err
	}
	if err := bootstrapper.Start(); err != nil {
		return err
	}
	// Ensure graceful resource teardown upon context cancellation
	defer func() { _ = bootstrapper.Close() }()

	// Block the main thread until the application context is canceled (e.g., SIGTERM)
	fmt.Printf("Bootstrap peerId=%s listen=%s announce=%s configDigest=%s\n", pid.String(), listen, announceIP, cc.ConfigDigest.Hex())
	<-ctx.Done()
	return nil
}

// runOracle handles the initialization and lifecycle of an active OCR3 Oracle node.
// It sets up the persistent network connections, injects dependencies into the 
// libocr core engine, and manages the execution of background asynchronous workers.
func runOracle(ctx context.Context, id, n, f int, seed int64, bootAddr, listen, announce string) error {
	// ==========================================
	// 1. Environment Configuration & SECRETS
	// ==========================================	
	privKeyHex := os.Getenv("PRIVATE_KEY")
	if privKeyHex == "" {
		log.Fatal("Critical Error: PRIVATE KEY variable not found. Check the docker-compose.yml")
	}

	// Strip the "0x" prefix if present to standardize the hex format
    privKeyHex = strings.TrimPrefix(privKeyHex, "0x")

	// Retrieval of the smart contract addresses from .env
	verifierAddressHex := os.Getenv("VERIFIER_ADDRESS")
	queueAddressHex := os.Getenv("QUEUE_ADDRESS") 
	ipfsUrl := getEnvironment("IPFS_API_URL", "http://ipfs:5001")

	if verifierAddressHex == "" || queueAddressHex == "" {
		log.Fatal("Critical: Contract addresses missing in .env")
	}

	// Change into Address type
	verifierAddress := common.HexToAddress(verifierAddressHex)
	rpc := getEnvironment("CHAIN_RPC", "http://localhost:8545")

	// ==========================================
	// 2. Off-Chain Infrastructure Boot
	// ==========================================	
	// Instantiate the IPFS Shell client for decentralized storage interaction
	ipfsShell := shell.NewShell(ipfsUrl)
	
	// Spawn the asynchronous Event Listener in a dedicated goroutine.
    // This decoupled architecture ensures that the blockchain subscription stream
    // does not block the OCR consensus operations.
	go startChainListener(ctx, rpc, verifierAddressHex, queueAddressHex, ipfsShell)

	
	// ==========================================
	// 3. P2P Networking Setup
	// ==========================================	
	announceIP, err := resolveHostnamePortToIP(announce)
	if err != nil { return err }
	bootIP, err := resolveHostnamePortToIP(bootAddr)
	if err != nil { return err }

	// Generate deterministic contract configuration and digest for the consensus
	cc, digester, err := buildContractConfig(n, f, seed)
	if err != nil { return err }

	nodes, _, err := deriveAllNodes(n, seed)
	if err != nil {
		return err
    }
	me := nodes[id]

	// Instantiate the P2P libp2p node for this oracle
	peer, err := networking.NewPeer(networking.PeerConfig{
		PrivKey: 				me.offKR.offPriv, 
		Logger: 				quietLogger{log.New(os.Stdout, "", 0)},
		V2ListenAddresses: 		[]string{listen}, 
		V2AnnounceAddresses: 	[]string{announceIP}, 
		V2EndpointConfig: networking.EndpointConfigV2{
			IncomingMessageBufferSize: 100,
			OutgoingMessageBufferSize: 50,
		},                      
        V2DeltaReconcile:            10 * time.Second,           
        MetricsRegisterer:           prometheus.NewRegistry(),   
	})

	if err != nil { return err }
	
	bootPriv := deriveBootstrapPriv(seed)
	bootPID, _ := ragetypes.PeerIDFromPrivateKey(bootPriv)

	// ==========================================
    // 4. EthTransmitter Initialization 
    // ==========================================
	// Pre-allocate RPC clients and cryptographic keys to guarantee low-latency 
    // on-chain transaction broadcasting during the transmission phase.
	transmitter, err := NewEthTransmitter(id, n, rpc, privKeyHex, verifierAddress, JobCache)
	if err != nil {
		return fmt.Errorf("Failed to initialize ethTransmitter: %w", err)
	}

	// ==========================================
	// 5. OCR3 Node Assembly (Dependency Injection)
	// ==========================================
	// Build the configuration arguments required by the libocr core state machine.
	args := offchainreporting2plus.OCR3OracleArgs[struct{}]{
		BinaryNetworkEndpointFactory: peer.OCR2BinaryNetworkEndpointFactory(),
		// Connect to the Bootstrap node for peer discovery		
		V2Bootstrappers: []commontypes.BootstrapperLocator{
			{PeerID: bootPID.String(), Addrs: []string{bootIP}},
		},
		ContractConfigTracker: staticTracker{cc},
		// Injecting our custom on-chain interaction logic
		ContractTransmitter: transmitter,
		Database: &memDB3{},
		// Tuning the protocol limits to handle large array submissions
		LocalConfig: ocrtypes.LocalConfig{
			DevelopmentMode: 					ocrtypes.EnableDangerousDevelopmentMode,
			BlockchainTimeout: 					30*time.Second, // Increased from 15 to 30
			ContractConfigConfirmations:		1, 
			ContractConfigTrackerPollInterval: 	5*time.Second, 
			ContractConfigLoadTimeout: 			10*time.Second,
			ContractTransmitterTransmitTimeout: 60*time.Second,  // Increased from 30 to 60
			DatabaseTimeout: 					10*time.Second,
		},
		Logger: 				quietLogger{log.New(os.Stdout, "", 0)},
		MetricsRegisterer: 		prometheus.NewRegistry(),
		MonitoringEndpoint: 	noopMonitoring{},
		OffchainConfigDigester: digester,
		OffchainKeyring: 		me.offKR, 
		OnchainKeyring: 		me.onKR,
		// Inject the custom Reporting Plugin containing our specific business logic (Attribution)
		ReportingPluginFactory: attributionPluginFactory{
			//oracleID:	id, 
			//numOracles: n,
			//rpcUrl: 	rpc,
			//contractAddr: verifierAddress,
		},
	}

	oracle, err := offchainreporting2plus.NewOracle(args)
	if err != nil { panic(err) }

	// Starts the internal state machine and consensus loops
	oracle.Start()
	fmt.Printf("ORACLE %d STARTED\n", id)

	// Block the main thread until the application context signals a shutdown
	<-ctx.Done()

	// Graceful teardown of network connections
	//if err := transmitter.Close(); err != nil {
	//	log.Printf("Error closing transmitter: %v", err)
	//}
	return nil
}


func main() {

	// CLI flags parsing for node routing and network configuration
	mode := flag.String("mode", "oracle", "bootstrap|oracle")
	n := flag.Int("n", 4, "")
	//f := flag.Int("f", 2, "")
	f := flag.Int("f", 1, "")
	id := flag.Int("oracle_id", 0, "")
	seed := flag.Int64("seed", 1, "")
	bootL := flag.String("bootstrap_listen", "", "")
	bootA := flag.String("bootstrap_announce", "", "")
	bootAddr := flag.String("bootstrap_addr", "", "")
	p2pL := flag.String("p2p_listen", "", "")
	p2pA := flag.String("p2p_announce", "", "")
	flag.Parse()
	
	// OS Signal trapping for graceful shutdown (SIGTERM, SIGINT)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Node Role routing
	if *mode == "bootstrap" {
		runBootstrap(ctx, *n, *f, *seed, *bootL, *bootA)
	} else {
		runOracle(ctx, *id, *n, *f, *seed, *bootAddr, *p2pL, *p2pA)
	}
}