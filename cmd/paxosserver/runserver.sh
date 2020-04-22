#! /bin/bash
set -e

go build

./paxosserver -port=8080 -addrs="localhost:8080,localhost:8081,localhost:8082" &
./paxosserver -port=8081 -addrs="localhost:8080,localhost:8081,localhost:8082" &
./paxosserver -port=8082 -addrs="localhost:8080,localhost:8081,localhost:8082" &

echo "running, enter to stop"

read && killall paxosserver