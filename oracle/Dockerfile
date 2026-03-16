FROM golang:1.24-alpine AS builder

WORKDIR /app

# 1. Install git to sync all the libraries 
RUN apk add --no-cache git

# 2. Copy all
COPY . .

# Then in the oracle directory
WORKDIR /app
RUN go mod tidy

# 3. Compile everything in the package main 
RUN go build -ldflags="-w -s" -o /out/oracle .
RUN go build -ldflags="-w -s" -o /out/checker ./check/status/main.go
RUN go build -ldflags="-w -s" -o /out/manual_checker ./check/manual/main.go

# --- STAGE 2: Runtime ---
FROM alpine:latest

WORKDIR /root/

# Install iproute2 to use tc command
RUN apk add --no-cache iproute2

COPY --from=builder /out/oracle /usr/local/bin/oracle
COPY --from=builder /out/checker /usr/local/bin/checker

# Copy the script in the container and make it executable
COPY wait-for-deploy.sh /usr/local/bin/wait-for-deploy.sh
RUN chmod +x /usr/local/bin/wait-for-deploy.sh

# Copy the network script in the container
COPY setup_network.sh /usr/local/bin/setup_network.sh
RUN chmod +x /usr/local/bin/setup_network.sh

# This script is going to be executed every time the container starts
# It will be this script to call wait-for-deploy
# Latency mode
ENTRYPOINT ["/usr/local/bin/setup_network.sh"]
# 0 Latency mode
#ENTRYPOINT ["/usr/local/bin/wait-for-deploy.sh"]

CMD ["oracle"]