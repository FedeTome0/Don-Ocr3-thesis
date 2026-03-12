const hre = require("hardhat");

/**
 * customerRequest.js
 * Simulates a customer placing a request by uploading data to IPFS
 * and submitting a payable transaction to the OracleQueue contract.
 */
async function main() {
    console.log("[CUSTOMER] Starting request and payment workflow...");

    // --- PHASE 1: IPFS UPLOAD ---
    const { create } = await import('kubo-rpc-client');
    const ipfs = create({ url: process.env.IPFS_API_URL || 'http://127.0.0.1:5001' });
    
    const payload = "Sample AI Input Data: " + Date.now();
    let cid;

    try {
        const result = await ipfs.add(payload);
        cid = result.cid.toString();
        console.log(`[IPFS] Upload Successful. CID: ${cid}`);
    } catch (error) {
        console.error("[IPFS] Upload failed:", error.message);
        return;
    }

    // --- PHASE 2: BLOCKCHAIN SUBMISSION ---
    const queueAddress = process.env.QUEUE_ADDRESS;
    const signers = await hre.ethers.getSigners();
    const customerWallet = signers[10]; // Standard test customer index
    
    const queueContract = (await hre.ethers.getContractAt("OracleQueue", queueAddress)).connect(customerWallet);

    try {
        const payment = hre.ethers.parseEther("0.02");
        console.log(`[CHAIN] Sending request with ${hre.ethers.formatEther(payment)} ETH payment...`);
        
        const tx = await queueContract.requestAttribution(cid, { value: payment });
        console.log(`[CHAIN] Transaction broadcasted. Hash: ${tx.hash}`);

        await tx.wait();
        console.log("[SUCCESS] Request accepted by the Smart Contract. Job state: PENDING.");
    } catch (error) {
        console.error("[CHAIN] Contract call failed:", error.message);
    }
}

main().catch(console.error);