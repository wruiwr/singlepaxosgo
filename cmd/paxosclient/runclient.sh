#! /bin/bash
set -e

go build

./paxosclient -addrs="localhost:8080,localhost:8081,localhost:8082" -clientRequest="M1"


read && killall paxosclient