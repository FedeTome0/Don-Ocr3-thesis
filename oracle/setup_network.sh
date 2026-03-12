#!/bin/sh

echo "[NETWORK] Traffic Control (tc) initialization for ORACLE_ID=${ORACLE_ID}..."

# Define the static IP of the subnet
IP_CHAIN="10.5.0.10"
IP_IPFS="10.5.0.11"
IP_OR0="10.5.0.20"
IP_OR1="10.5.0.21"
IP_OR2="10.5.0.22"
IP_OR3="10.5.0.23"
# More Oracles
IP_OR4="10.5.0.24"
IP_OR5="10.5.0.25"
IP_OR6="10.5.0.26"

# Clean previous rules and create the root
tc qdisc del dev eth0 root 2>/dev/null || true
#tc qdisc add dev eth0 root handle 1: prio bands 6
# For more oracles
tc qdisc add dev eth0 root handle 1: prio bands 9


# Apply the logic based on the current node
if [ "$ORACLE_ID" = "0" ]; then
    echo "[NETWORK] Node location: Milan"
    
    # Creation of the routes with specific delays
    #tc qdisc add dev eth0 parent 1:1 handle 10: netem delay 15ms  # To Chain/IPFS
    tc qdisc add dev eth0 parent 1:2 handle 20: netem delay 90ms  # To NY
    tc qdisc add dev eth0 parent 1:3 handle 30: netem delay 120ms # To Mumbai
    tc qdisc add dev eth0 parent 1:4 handle 40: netem delay 175ms # To Johannesburg
    tc qdisc add dev eth0 parent 1:5 handle 50: netem delay 50ms # To Lisbon
    tc qdisc add dev eth0 parent 1:6 handle 60: netem delay 45ms # To Moscow
    tc qdisc add dev eth0 parent 1:7 handle 70: netem delay 110ms # To Toronto


    # Filters the traffic based on the IP
    #tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst $IP_CHAIN flowid 1:1
    #tc filter add dev eth0 protocol ip parent 1:0 prio 2 u32 match ip dst $IP_IPFS flowid 1:1
    tc filter add dev eth0 protocol ip parent 1:0 prio 3 u32 match ip dst $IP_OR1 flowid 1:2
    tc filter add dev eth0 protocol ip parent 1:0 prio 4 u32 match ip dst $IP_OR2 flowid 1:3
    tc filter add dev eth0 protocol ip parent 1:0 prio 5 u32 match ip dst $IP_OR3 flowid 1:4
    tc filter add dev eth0 protocol ip parent 1:0 prio 6 u32 match ip dst $IP_OR4 flowid 1:5
    tc filter add dev eth0 protocol ip parent 1:0 prio 7 u32 match ip dst $IP_OR5 flowid 1:6
    tc filter add dev eth0 protocol ip parent 1:0 prio 8 u32 match ip dst $IP_OR6 flowid 1:7


elif [ "$ORACLE_ID" = "1" ]; then
    echo "[NETWORK] Node location: New York"
    #tc qdisc add dev eth0 parent 1:1 handle 10: netem delay 80ms  # To Chain/IPFS
    tc qdisc add dev eth0 parent 1:2 handle 20: netem delay 90ms  # To Milano
    tc qdisc add dev eth0 parent 1:3 handle 30: netem delay 190ms # To Mumbai
    tc qdisc add dev eth0 parent 1:4 handle 40: netem delay 230ms # To Johannesburg
    tc qdisc add dev eth0 parent 1:5 handle 50: netem delay 115ms # To Lisbon
    tc qdisc add dev eth0 parent 1:6 handle 60: netem delay 120ms # To Moscow
    tc qdisc add dev eth0 parent 1:7 handle 70: netem delay 17ms # To Toronto


    #tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst $IP_CHAIN flowid 1:1
    #tc filter add dev eth0 protocol ip parent 1:0 prio 2 u32 match ip dst $IP_IPFS flowid 1:1
    tc filter add dev eth0 protocol ip parent 1:0 prio 3 u32 match ip dst $IP_OR0 flowid 1:2
    tc filter add dev eth0 protocol ip parent 1:0 prio 4 u32 match ip dst $IP_OR2 flowid 1:3
    tc filter add dev eth0 protocol ip parent 1:0 prio 5 u32 match ip dst $IP_OR3 flowid 1:4
    tc filter add dev eth0 protocol ip parent 1:0 prio 6 u32 match ip dst $IP_OR4 flowid 1:5
    tc filter add dev eth0 protocol ip parent 1:0 prio 7 u32 match ip dst $IP_OR5 flowid 1:6
    tc filter add dev eth0 protocol ip parent 1:0 prio 8 u32 match ip dst $IP_OR6 flowid 1:7

elif [ "$ORACLE_ID" = "2" ]; then
    echo "[NETWORK] Node location: Mumbai"
    #tc qdisc add dev eth0 parent 1:1 handle 10: netem delay 80ms  # To Chain/IPFS
    tc qdisc add dev eth0 parent 1:2 handle 20: netem delay 120ms  # To Milano
    tc qdisc add dev eth0 parent 1:3 handle 30: netem delay 190ms # To NewYork
    tc qdisc add dev eth0 parent 1:4 handle 40: netem delay 290ms # To Johannesburg
    tc qdisc add dev eth0 parent 1:5 handle 50: netem delay 150ms # To Lisbon
    tc qdisc add dev eth0 parent 1:6 handle 60: netem delay 175ms # To Moscow
    tc qdisc add dev eth0 parent 1:7 handle 70: netem delay 240ms # To Toronto


    #tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst $IP_CHAIN flowid 1:1
    #tc filter add dev eth0 protocol ip parent 1:0 prio 2 u32 match ip dst $IP_IPFS flowid 1:1
    tc filter add dev eth0 protocol ip parent 1:0 prio 3 u32 match ip dst $IP_OR0 flowid 1:2
    tc filter add dev eth0 protocol ip parent 1:0 prio 4 u32 match ip dst $IP_OR1 flowid 1:3
    tc filter add dev eth0 protocol ip parent 1:0 prio 5 u32 match ip dst $IP_OR3 flowid 1:4
    tc filter add dev eth0 protocol ip parent 1:0 prio 6 u32 match ip dst $IP_OR4 flowid 1:5
    tc filter add dev eth0 protocol ip parent 1:0 prio 7 u32 match ip dst $IP_OR5 flowid 1:6
    tc filter add dev eth0 protocol ip parent 1:0 prio 8 u32 match ip dst $IP_OR6 flowid 1:7

elif [ "$ORACLE_ID" = "3" ]; then
    echo "[NETWORK] Node location: Johannesburg"
    #tc qdisc add dev eth0 parent 1:1 handle 10: netem delay 80ms  # To Chain/IPFS
    tc qdisc add dev eth0 parent 1:2 handle 20: netem delay 175ms  # To Milano
    tc qdisc add dev eth0 parent 1:3 handle 30: netem delay 230ms # To NewYork
    tc qdisc add dev eth0 parent 1:4 handle 40: netem delay 290ms # To Mumbai
    tc qdisc add dev eth0 parent 1:5 handle 50: netem delay 210ms # To Lisbon
    tc qdisc add dev eth0 parent 1:6 handle 60: netem delay 200ms # To Moscow
    tc qdisc add dev eth0 parent 1:7 handle 70: netem delay 225ms # To Toronto

    #tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst $IP_CHAIN flowid 1:1
    #tc filter add dev eth0 protocol ip parent 1:0 prio 2 u32 match ip dst $IP_IPFS flowid 1:1
    tc filter add dev eth0 protocol ip parent 1:0 prio 3 u32 match ip dst $IP_OR0 flowid 1:2
    tc filter add dev eth0 protocol ip parent 1:0 prio 4 u32 match ip dst $IP_OR1 flowid 1:3
    tc filter add dev eth0 protocol ip parent 1:0 prio 5 u32 match ip dst $IP_OR2 flowid 1:4
    tc filter add dev eth0 protocol ip parent 1:0 prio 6 u32 match ip dst $IP_OR4 flowid 1:5
    tc filter add dev eth0 protocol ip parent 1:0 prio 7 u32 match ip dst $IP_OR5 flowid 1:6
    tc filter add dev eth0 protocol ip parent 1:0 prio 8 u32 match ip dst $IP_OR6 flowid 1:7

elif [ "$ORACLE_ID" = "4" ]; then
    echo "[NETWORK] Node location: Lisbon"
    #tc qdisc add dev eth0 parent 1:1 handle 10: netem delay 80ms  # To Chain/IPFS
    tc qdisc add dev eth0 parent 1:2 handle 20: netem delay 50ms  # To Milano
    tc qdisc add dev eth0 parent 1:3 handle 30: netem delay 110ms # To NewYork
    tc qdisc add dev eth0 parent 1:4 handle 40: netem delay 150ms # To Mumbai
    tc qdisc add dev eth0 parent 1:5 handle 50: netem delay 210ms # To Johannesburg
    tc qdisc add dev eth0 parent 1:6 handle 60: netem delay 80ms # To Moscow
    tc qdisc add dev eth0 parent 1:7 handle 70: netem delay 120ms # To Toronto


    #tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst $IP_CHAIN flowid 1:1
    #tc filter add dev eth0 protocol ip parent 1:0 prio 2 u32 match ip dst $IP_IPFS flowid 1:1
    tc filter add dev eth0 protocol ip parent 1:0 prio 3 u32 match ip dst $IP_OR0 flowid 1:2
    tc filter add dev eth0 protocol ip parent 1:0 prio 4 u32 match ip dst $IP_OR1 flowid 1:3
    tc filter add dev eth0 protocol ip parent 1:0 prio 5 u32 match ip dst $IP_OR2 flowid 1:4
    tc filter add dev eth0 protocol ip parent 1:0 prio 6 u32 match ip dst $IP_OR3 flowid 1:5
    tc filter add dev eth0 protocol ip parent 1:0 prio 7 u32 match ip dst $IP_OR5 flowid 1:6
    tc filter add dev eth0 protocol ip parent 1:0 prio 8 u32 match ip dst $IP_OR6 flowid 1:7

elif [ "$ORACLE_ID" = "5" ]; then
    echo "[NETWORK] Node location: Moscow"
    #tc qdisc add dev eth0 parent 1:1 handle 10: netem delay 80ms  # To Chain/IPFS
    tc qdisc add dev eth0 parent 1:2 handle 20: netem delay 45ms  # To Milano
    tc qdisc add dev eth0 parent 1:3 handle 30: netem delay 120ms # To NewYork
    tc qdisc add dev eth0 parent 1:4 handle 40: netem delay 175ms # To Mumbai
    tc qdisc add dev eth0 parent 1:5 handle 50: netem delay 200ms # To Johannesburg
    tc qdisc add dev eth0 parent 1:6 handle 60: netem delay 80ms # To Lisbon
    tc qdisc add dev eth0 parent 1:7 handle 70: netem delay 140ms # To Toronto

    #tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst $IP_CHAIN flowid 1:1
    #tc filter add dev eth0 protocol ip parent 1:0 prio 2 u32 match ip dst $IP_IPFS flowid 1:1
    tc filter add dev eth0 protocol ip parent 1:0 prio 3 u32 match ip dst $IP_OR0 flowid 1:2
    tc filter add dev eth0 protocol ip parent 1:0 prio 4 u32 match ip dst $IP_OR1 flowid 1:3
    tc filter add dev eth0 protocol ip parent 1:0 prio 5 u32 match ip dst $IP_OR2 flowid 1:4
    tc filter add dev eth0 protocol ip parent 1:0 prio 6 u32 match ip dst $IP_OR3 flowid 1:5
    tc filter add dev eth0 protocol ip parent 1:0 prio 7 u32 match ip dst $IP_OR4 flowid 1:6
    tc filter add dev eth0 protocol ip parent 1:0 prio 8 u32 match ip dst $IP_OR6 flowid 1:7

elif [ "$ORACLE_ID" = "6" ]; then
    echo "[NETWORK] Node location: Toronto"
    #tc qdisc add dev eth0 parent 1:1 handle 10: netem delay 80ms  # To Chain/IPFS
    tc qdisc add dev eth0 parent 1:2 handle 20: netem delay 110ms  # To Milano
    tc qdisc add dev eth0 parent 1:3 handle 30: netem delay 17ms # To NewYork
    tc qdisc add dev eth0 parent 1:4 handle 40: netem delay 240ms # To Mumbai
    tc qdisc add dev eth0 parent 1:5 handle 50: netem delay 230ms # To Johannesburg
    tc qdisc add dev eth0 parent 1:6 handle 60: netem delay 120ms # To Lisbon
    tc qdisc add dev eth0 parent 1:7 handle 70: netem delay 140ms # To Moscow

    #tc filter add dev eth0 protocol ip parent 1:0 prio 1 u32 match ip dst $IP_CHAIN flowid 1:1
    #tc filter add dev eth0 protocol ip parent 1:0 prio 2 u32 match ip dst $IP_IPFS flowid 1:1
    tc filter add dev eth0 protocol ip parent 1:0 prio 3 u32 match ip dst $IP_OR0 flowid 1:2
    tc filter add dev eth0 protocol ip parent 1:0 prio 4 u32 match ip dst $IP_OR1 flowid 1:3
    tc filter add dev eth0 protocol ip parent 1:0 prio 5 u32 match ip dst $IP_OR2 flowid 1:4
    tc filter add dev eth0 protocol ip parent 1:0 prio 6 u32 match ip dst $IP_OR3 flowid 1:5
    tc filter add dev eth0 protocol ip parent 1:0 prio 7 u32 match ip dst $IP_OR4 flowid 1:6
    tc filter add dev eth0 protocol ip parent 1:0 prio 8 u32 match ip dst $IP_OR5 flowid 1:7
else
    echo "[NETWORK] No latency specified for this container"
fi 

echo "[NETWORK] Kernel routing rules applied succesfully"
echo "[NETWORK] Pass the control to wait-for-deploy.sh..."

exec /usr/local/bin/wait-for-deploy.sh "$@"