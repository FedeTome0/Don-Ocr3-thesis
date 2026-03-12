FROM node:18-alpine

WORKDIR /app

# Init Json and install hardhat and other libraries
RUN npm init -y
RUN npm install --save-dev hardhat@2.19.1 @nomicfoundation/hardhat-ethers@3.0.5 ethers@6.11.1 dotenv kubo-rpc-client

# 4. Copy the files
COPY hardhat.config.js ./
COPY contracts ./contracts
COPY scripts ./scripts

EXPOSE 8545

RUN apk add --no-cache curl

RUN echo '#!/bin/sh' > /start.sh && \
    echo 'echo "--- STARTING HARDHAT NODE ---"' >> /start.sh && \
    echo 'npx hardhat node --hostname 0.0.0.0 | grep -v "eth_" &' >> /start.sh && \
    echo 'echo "--- WAITING FOR CHAIN ---"' >> /start.sh && \
    echo 'until curl -sf -X POST http://localhost:8545 -H "Content-Type: application/json" -d '"'"'{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'"'"' > /dev/null; do' >> /start.sh && \
    echo '  echo "Chain not ready, retrying..."' >> /start.sh && \
    echo '  sleep 2' >> /start.sh && \
    echo 'done' >> /start.sh && \
    echo 'echo "--- DEPLOYING ALL CONTRACTS ---"' >> /start.sh && \
    echo 'npx hardhat run scripts/deployContracts.js --network localhost || echo "DEPLOYMENT FAILED!"' >> /start.sh && \
    echo 'echo "--- CHAIN READY ---"' >> /start.sh && \
    echo 'wait' >> /start.sh && \
    chmod +x /start.sh

    # 5. Script for the start of hardhat and the deployment of the two smart contract
    #RUN echo '#!/bin/sh' > /start.sh && \
    #echo 'echo "--- STARTING HARDHAT NODE ---"' >> /start.sh && \
    #echo 'npx hardhat node --hostname 0.0.0.0 &' >> /start.sh && \
    #echo 'echo "--- WAITING 5s FOR CHAIN TO BOOT ---"' >> /start.sh && \
    #echo 'sleep 5' >> /start.sh && \
    #echo 'echo "--- DEPLOYING VERIFIER ---"' >> /start.sh && \
    #echo 'npx hardhat run scripts/deploy.js --network localhost || echo "VERIFIER DEPLOY FAILED!"' >> /start.sh && \
    #echo 'echo "--- DEPLOYING ORACLE QUEUE ---"' >> /start.sh && \
    #echo 'npx hardhat run scripts/deploy_queue.js --network localhost || echo "QUEUE DEPLOY FAILED!"' >> /start.sh && \
    #echo 'echo "--- CHAIN READY ---"' >> /start.sh && \
    #echo 'wait' >> /start.sh && \
    #chmod +x /start.sh

CMD ["/start.sh"]