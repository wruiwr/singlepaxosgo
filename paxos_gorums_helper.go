package singlepaxos

import (
	"log"
	time "time"

	grpc "google.golang.org/grpc"
)

// newPaxosConfig creates a Gorums manager and configuration and
// a Paxos quorum specification based on the provided quorumSize.
// The returned manager should be closed when no longer needed.
func newPaxosConfig(addrs []string, quorumSize int) (*Configuration, *Manager) {
	grpcOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second),
		grpc.WithInsecure(),
	}
	mgrOpts := []ManagerOption{WithGrpcDialOptions(grpcOpts...), WithTracing()}
	mgr, err := NewManager(addrs, mgrOpts...)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	// get all available node ids from the manager
	ids := mgr.NodeIDs()
	// create new quorum specification for Paxos
	qspec := NewPaxosQSpec(quorumSize)
	// create new configuration with all node ids and the quorum specification
	config, err := mgr.NewConfiguration(ids, qspec)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return config, mgr
}
