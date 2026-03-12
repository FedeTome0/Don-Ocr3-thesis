const hre = require("hardhat");
const { ethers } = hre;

/**
 * ModelCreatorApprove.js
 * Acts as the 'Model Creator' authority. It monitors the OracleQueue for new
 * customer requests and formally approves them to trigger the DON consensus.
 */
async function main() {
  console.log("[MODEL CREATOR] Initializing validation and approval service...");

  const queueAddress = process.env.QUEUE_ADDRESS;
  const verifierAddress = process.env.VERIFIER_ADDRESS;

  const [creatorWallet] = await hre.ethers.getSigners();

  // Attach to contracts using the Creator's identity
  const queueContract = (await hre.ethers.getContractAt("OracleQueue", queueAddress)).connect(creatorWallet);
  const verifierContract = await hre.ethers.getContractAt("OracleVerifier", verifierAddress);

  console.log(`[MODEL CREATOR] Monitoring 'LogNewCustomerRequest' events...`);

  const MAX_RETRIES = 3;
  
  // Sequential Job Queue: ensures jobs are approved and finalized one by one
  // to prevent nonce collisions and maintain deterministic benchmark results.
  let jobProcessingPipeline = Promise.resolve();

  queueContract.on("LogNewCustomerRequest", async (requestId, ipfsCid, customer, payment) => {
    console.log(`\n[EVENT] New Job Detected: #${requestId}`);
    console.log(`       CID:      ${ipfsCid}`);
    console.log(`       Value:    ${ethers.formatEther(payment)} ETH`);

    // Add job to the sequential pipeline
    jobProcessingPipeline = jobProcessingPipeline.then(async () => {
      
      // --- STEP 1: ON-CHAIN APPROVAL ---
      for (let attempt = 1; attempt <= MAX_RETRIES; attempt++) {
        try {
          console.log(`[PROCESS] Approving job #${requestId} (Attempt ${attempt}/${MAX_RETRIES})...`);
          const tx = await queueContract.approveJob(requestId);
          await tx.wait();
          console.log(`[SUCCESS] Job #${requestId} approved. Oracles notified via LogNewJobForOracles.`);
          break;
        } catch (error) {
          if (attempt === MAX_RETRIES) {
            console.error(`[ERROR] Job #${requestId} approval failed permanently: ${error.message}`);
            return;
          }
          await new Promise(res => setTimeout(res, 2000));
        }
      }

      // --- STEP 2: WAIT FOR DON FULFILLMENT ---
      // We use a Promise with .once() to wait for the consensus result to land on-chain.
      console.log(`[WAIT] Awaiting OCR consensus for job #${requestId}...`);
      await new Promise((resolve, reject) => {
        const timeout = setTimeout(() => reject(new Error("OCR fulfillment timeout (10m)")), 600000);

        verifierContract.once("JobCompleted", (completedId, submitter) => {
          if (completedId.toString() === requestId.toString()) {
            clearTimeout(timeout);
            console.log(`[DONE] Job #${requestId} finalized by Oracle: ${submitter}.`);
            resolve();
          }
        });
      });
    });
  });

  // Keep the process alive
  await new Promise(() => {});
}

main().catch((error) => {
  console.error("Fatal service error:", error);
  process.exit(1);
});