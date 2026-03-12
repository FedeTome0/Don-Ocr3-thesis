const hre = require("hardhat");
const fs = require("fs");

/**
 * benchmark.js
 * Automated latency evaluation script for the Decentralized Oracle Network (DON).
 * Measures performance across four distinct phases: IPFS storage, on-chain 
 * request, creator approval, and OCR3 BFT consensus/fulfillment.
 */
async function main() {
    console.log("[BENCHMARK] Starting automated DON latency evaluation...\n");

    // Dynamic import for 'kubo-rpc-client' (ESM module compatibility)
    const { create } = await import('kubo-rpc-client');
    const ipfsUrl = process.env.IPFS_API_URL || 'http://127.0.0.1:5001';
    const ipfs = create({ url: ipfsUrl });

    const queueAddress = process.env.QUEUE_ADDRESS;
    const verifierAddress = process.env.VERIFIER_ADDRESS;

    // Use Signer #10 as the 'Customer' to prevent nonce collisions with Oracles or Creator
    const signers = await hre.ethers.getSigners();
    const customerWallet = signers[10];

    const queueContract = await hre.ethers.getContractAt("OracleQueue", queueAddress, customerWallet);
    const verifierContract = await hre.ethers.getContractAt("OracleVerifier", verifierAddress);

    // Initialize CSV Telemetry Data
    const csvFile = "benchmark_results.csv";
    const csvHeader = "Iteration,Timestamp,OffChain_Storage(ms),OnChain_RequestTx(ms),OnChain_ApprovalTx(ms),OCR_Consensus_And_FulfillmentTx(ms),Total_Latency(ms)\n";
    fs.writeFileSync(csvFile, csvHeader);

    for (let i = 1; i <= 20; i++) {
        console.log(`========================================`);
        console.log(` ITERATION ${i} / 20`);
        console.log(`========================================`);

        // IMPORTANT: Removed '#' character to prevent AI Tokenizer KeyError
        const simplePayload = `Automated Benchmark Prompt Payload number ${i}...`;
        const t0 = performance.now();

        // ---------------------------------------------------------------------
        // PHASE 1: Off-Chain Storage (IPFS Upload)
        // ---------------------------------------------------------------------
        process.stdout.write('[1/4] Off-Chain Storage (IPFS Upload)... ');
        const { cid } = await ipfs.add(simplePayload);
        const cidString = cid.toString();
        const tIpfs = performance.now();
        console.log(`[+] ${(tIpfs - t0).toFixed(2)} ms | CID: ${cidString}`);

        // ---------------------------------------------------------------------
        // PHASE 2: Customer Request (On-Chain Transaction)
        // ---------------------------------------------------------------------
        // Pre-register listener for the Approval event to avoid race conditions
        const approvalPromise = new Promise((resolve, reject) => {
            const timeout = setTimeout(() => reject(new Error("Timeout: LogNewJobForOracles event missed")), 45000);
            
            queueContract.once("LogNewJobForOracles", (jobId) => {
                clearTimeout(timeout);
                resolve(jobId);
            });
        });

        process.stdout.write(`[2/4] On-Chain Customer Request (Tx)... `);
        const paymentAmount = hre.ethers.parseEther("0.02");
        
        const tx = await queueContract.requestAttribution(cidString, { value: paymentAmount });
        const receipt = await tx.wait();
        const tPhase1 = performance.now();
        console.log(`[+] ${(tPhase1 - tIpfs).toFixed(0)} ms`);

        // Extract RequestID from the transaction logs
        let currentJobId = null;
        for (const log of receipt.logs) {
            try {
                const parsed = queueContract.interface.parseLog(log);
                if (parsed.name === "LogNewCustomerRequest") {
                    currentJobId = parsed.args[0];
                    break;
                }
            } catch (e) { /* Skip unrelated logs */ }
        }

        // ---------------------------------------------------------------------
        // PHASE 3: Validation & Approval (Model Creator Tx)
        // ---------------------------------------------------------------------
        process.stdout.write(`[3/4] On-Chain Approval (Model Creator Tx)... `);
        await approvalPromise; // Resolves when the separate Creator script approves the job
        const tPhase2 = performance.now();
        console.log(`[+] ${(tPhase2 - tPhase1).toFixed(2)} ms`);

        // ---------------------------------------------------------------------
        // PHASE 4: OCR Network Execution (AI Inference + P2P Consensus)
        // ---------------------------------------------------------------------
        console.log(`[4/4] Off-Chain Reporting (AI + BFT Consensus)... Waiting`);
        let winnerAddress = "";

        await new Promise((resolve, reject) => {
            const timeout = setTimeout(() => {
                reject(new Error(`Timeout: JobCompleted event missed for Job #${currentJobId}`));
            }, 300000); // 5-minute threshold for AI computation and P2P rounds

            const fulfillmentListener = (jobId, submitter) => {
                if (currentJobId !== null && jobId.toString() === currentJobId.toString()) {
                    clearTimeout(timeout);
                    winnerAddress = submitter;
                    verifierContract.off("JobCompleted", fulfillmentListener);
                    resolve();
                }
            };

            verifierContract.on("JobCompleted", fulfillmentListener);
        });

        const tPhase3 = performance.now();
        console.log(`      OCR Fulfillment On-Chain Tx detected in ${(tPhase3 - tPhase2).toFixed(2)} ms`);
        console.log(`      Transmitter Node Identity: ${winnerAddress}\n`);

        // ---------------------------------------------------------------------
        // TELEMETRY EXPORT
        // ---------------------------------------------------------------------
        const timeIpfs     = (tIpfs   - t0).toFixed(2);
        const timePhase1   = (tPhase1 - tIpfs).toFixed(2);
        const timePhase2   = (tPhase2 - tPhase1).toFixed(2);
        const timePhase3   = (tPhase3 - tPhase2).toFixed(2);
        const totalTime    = (tPhase3 - t0).toFixed(2);

        const csvRow = `${i},${Date.now()},${timeIpfs},${timePhase1},${timePhase2},${timePhase3},${totalTime}\n`;
        fs.appendFileSync(csvFile, csvRow);
        
        console.log(`[OK] Iteration ${i} saved. Total Latency: ${totalTime} ms\n`);

        // RPC Cooldown period
        await new Promise(resolve => setTimeout(resolve, 3000));
    }
    
    console.log(`\n Benchmark completed successfully! Data exported to: ${csvFile}`);
}

main().catch((error) => {
    console.error("Benchmark encountered a fatal error:", error.message);
    process.exit(1);
});