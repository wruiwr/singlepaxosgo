package singlepaxos

import (
	"fmt"
	"log"
	"net"

	"sync"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Initial values in each replica.
const (
	NoRound uint32 = 0
	Ignore  uint32 = 0
)

// PaxosReplica is the structure composing the Proposer, Acceptor, and Learner.
type PaxosReplica struct {
	SinglePaxosServer
	*Acceptor
	*Proposer
	config *Configuration

	decidedValue    *Response
	decidedValueOut chan Value
	m               sync.RWMutex
}

// NewPaxosReplica returns a new Paxos replica with a configuration as provided
// by the input addrs. This replica will run on the given port.
func NewPaxosReplica(port int, config *Configuration) *PaxosReplica {
	return &PaxosReplica{
		config:          config,
		Acceptor:        NewAcceptor(),
		Proposer:        NewProposer(config, port),
		decidedValueOut: make(chan Value, 16),
	}
}

// ServerStart setups the server connection and starts it
func ServerStart(port int, addrs []string, quorumSize int) {
	l, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	// setup connection between Paxos replicas based on the provided addrs and quorum size
	config, mgr := newPaxosConfig(addrs, quorumSize)
	defer mgr.Close()
	// create new Paxos replica instance running on the given port
	replica := NewPaxosReplica(port, config)
	replica.Start()

	// create new grpc server instance
	grpcServer := grpc.NewServer([]grpc.ServerOption{}...)
	// register Paxos replica with Gorums runtime
	RegisterSinglePaxosServer(grpcServer, replica)

	replica.logf("Starting server on port %d\n", port)
	log.Fatal(grpcServer.Serve(l))
}

// Ping just replies with the same heartbeat message.
func (r *PaxosReplica) Ping(ctx context.Context, hb *Heartbeat) (*Heartbeat, error) {
	return hb, nil
}

// Prepare handles the prepare quorum calls from the proposer by passing the received messages to its acceptor.
// It receives prepare massages and pass them to handlePrepare method of acceptor.
// It returns promise messages back to the proposer by its acceptor.
func (r *PaxosReplica) Prepare(ctx context.Context, prepMsg *PrepareMsg) (*PromiseMsg, error) {
	r.logf("Acceptor.Prepare(%v) received\n", prepMsg)
	prm, _ := r.handlePrepare(prepMsg)
	return prm, nil
}

// Accept handles the accept quorum calls from the proposer by passing the received messages to its acceptor.
// It receives Accept massages and pass them to handleAccept method of acceptor.
// It returns learn massages back to the proposer by its acceptor
func (r *PaxosReplica) Accept(ctx context.Context, accMsg *AcceptMsg) (*LearnMsg, error) {
	r.logf("Acceptor.Accept(%v) received\n", accMsg)
	lrn, _ := r.handleAccept(accMsg)
	return lrn, nil
}

// Commit handles the commit quorum calls from the proposer.
// It receives a learn massage from proposer and deliver its decided value to .
// It returns a empty massage back.
func (r *PaxosReplica) Commit(ctx context.Context, lrnMsg *LearnMsg) (*Empty, error) {
	r.logf("Learner.Commit(%v) received\n", lrnMsg)
	// deliver decided value to ClientHandle
	r.decidedValueOut <- *lrnMsg.Val
	return &Empty{}, nil
}

func (r *PaxosReplica) ClientHandle(ctx context.Context, req *Value) (rsp *Response, err error) {
	r.logf("Replica.ClientHandle(%v) received\n", req)
	r.m.Lock()
	defer r.m.Unlock()
	if r.decidedValue != nil {
		r.logf("Replica.ClientHandle() returning already decided value: %v\n", r.decidedValue)
		return r.decidedValue, nil
	}
	r.ProposeClientValue(req)

	// blocks until we have a decided value
	dv := <-r.decidedValueOut
	// create a response according to decided value
	r.decidedValue = &Response{dv.GetClientRequest()}
	r.logf("Replica.ClientHandle() got decided value: %v\n", r.decidedValue)
	// we have decided, we can stop the proposer
	r.Stop()
	return r.decidedValue, nil
}
