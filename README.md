# Decentralized Oracle Network (DON) for LLM Royalties

This repository contains my Bachelor's thesis project, focused on the development of a Decentralized Oracle Network (DON) utilizing Chainlink's architecture and libraries.

The project proposes a novel decentralized application designed to address fair compensation in the AI industry. Specifically, it attributes and distributes royalties to the original owners of the data used to train Large Language Models (LLMs). By leveraging the off-chain computation capabilities of the DON, custom oracles securely execute an attribution algorithm to calculate the correct royalty distribution, ensuring a transparent, trustless, and decentralized reward system.

## Project Structure

* **chain**: Contains the Solidity smart contracts (`OracleQueue.sol`, `OracleVerifier.sol`), Hardhat configuration, and deployment scripts.
* **oracle**: Contains the off-chain Oracle node backend written in Go, including the custom OCR3 (Off-Chain Reporting) plugin and listener.
* **IpfsAgent**: Scripts and configurations for handling data storage and retrieval via IPFS.
* **docker-compose.yml**: Orchestrates the Docker containers for the oracle network and testing environment.

## Tech Stack
* **Blockchain/Smart Contracts:** Solidity, Hardhat, Ethers.js
* **Oracle Infrastructure:** Chainlink DON, Go (Golang), Docker
* **Decentralized Storage:** IPFS

---

## Prerequisites

To run this project locally, ensure you have the following installed on your machine:
* **Node.js** (v24.13.0 or higher)
* **Go** (v1.23.4 linux/amd64 or higher)
* **Docker** and **Docker Compose** (v29.2.0 and Docker Compose version v5.0.2)
* **Git**

---

## Installation & Setup

1. **Clone the repository**
    git clone https://github.com/FedeTome0/Don-Ocr3-thesis.git
    cd Don-Ocr3-thesis
2. **Install Smart Contract dependencies (Hardhat)**
    cd chain
    npm install
    cd ..

---

## How to Run the Project locally 
Thanks to the Docker orchestration, booting the entire ecosystem is highly automated. 
The docker-compose setup will:

- Start the Hardhat local node
- Wait for it to be ready
- Deploy the smart contracts 
- Boot the off-chain oracle nodes.

To test the full system, follow this chronological sequence.

---
## Step 1: Start the AI Backend and SSH Tunnel

The oracles communicate with the Python AI server hosted on the remote machine (`satoshi`).

### 1. Access the remote server and start the Python model service (All in satoshi)

```bash
python model_service_v2.py
```

### 2. Open a new local terminal and create an SSH tunnel

```bash
ssh -L 50100:127.0.0.1:50100 tomelliniT@131.114.50.205
```

Keep this terminal open to maintain the connection.

---

## Step 2: Boot the Core Infrastructure

Open a new terminal in the root directory of the project and run:

```bash
docker compose up --build
```

### What happens under the hood

- A local Hardhat blockchain is initialized
- The `OracleVerifier` and `OracleQueue` smart contracts are compiled and deployed automatically
- IPFS nodes are started for decentralized storage
- The Go-based OCR3 oracle nodes are built and begin listening to the blockchain

Wait until you see the message:

```
--- CHAIN READY ---
```

in the Docker logs before proceeding.

---

## Step 3: Simulate the Workflow

Now you simulate the two actors of the system:

- Model Creator
- Customer

---

### 1. Start the Model Creator Listener

Open a new terminal:

```bash
cd chain
npx hardhat run scripts/modelCreatorApprove.js --network localhost
```

This script simulates the model creator listening for OracleQueue request events.

---

### 2. Trigger the Customer Request

Open another terminal:

```bash
cd chain
npx hardhat run scripts/customerRequest.js --network localhost
```

At this point:

- The computation starts
- The oracle consensus begins
- Royalties are calculated

You can monitor everything in the main Docker terminal.

---

## How to Stop the Environment

To stop execution:

Press `CTRL+C` in the Docker terminal.

To completely reset everything:

```bash
docker compose down -v
```