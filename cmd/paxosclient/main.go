package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	qc "github.com/selabhvl/singlepaxos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	/* get initial values from inputs */
	var (
		saddrs        = flag.String("addrs", "", "server addresses separated by ','")
		clientRequest = flag.String("clientRequest", "", "client request")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	addrs := strings.Split(*saddrs, ",")
	if len(addrs) == 0 {
		log.Fatalln("no server addresses provided")
	}

	// start a initial proposer
	ClientStart(addrs, clientRequest)
}

func ClientStart(addrs []string, clientRequest *string) {
	fmt.Printf("Connecting to %d Paxos replicas: %v\n", len(addrs), addrs)

	// setup the client connection to servers
	config, mgr := setupClientConn(addrs, 2) // TODO avoid hardcoded quorum size
	defer mgr.Close()

	req := qc.Value{ClientRequest: *clientRequest}
	resp := doSendRequest(config, req)
	fmt.Println("response:", resp.ClientResponse, "for the client request:", req)
}

// Internal: doSendRequest can send requests to paxos servers by quorum call
func doSendRequest(config *qc.Configuration, value qc.Value) *qc.Response {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := config.ClientHandle(ctx, &value)
	if err != nil {
		log.Fatalf("ClientHandle quorum call error: %v", err)
	}
	return resp
}

// setupClientConn creates a manager and configuration and
// a Paxos quorum specification based on the provided quorumSize.
func setupClientConn(addrs []string, quorumSize int) (*qc.Configuration, *qc.Manager) {
	grpcOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
		grpc.WithInsecure(),
	}
	mgrOpts := []qc.ManagerOption{qc.WithGrpcDialOptions(grpcOpts...), qc.WithTracing()}
	mgr, err := qc.NewManager(addrs, mgrOpts...)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	// get all available node ids from the manager
	ids := mgr.NodeIDs()
	// create new quorum specification for Paxos
	qspec := qc.NewPaxosQSpec(quorumSize)
	// create new configuration with all node ids and the quorum specification
	config, err := mgr.NewConfiguration(ids, qspec)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return config, mgr
}
