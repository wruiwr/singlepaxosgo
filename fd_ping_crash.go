// +build crash

package singlepaxos

import (
	"log"

	"golang.org/x/net/context"
)

func (fd *failureDetector) ping(node *Node) {
	crashedNodeMutex.Lock()
	defer crashedNodeMutex.Unlock()

	for _, cn := range crashedNodes {
		if node.ID() == cn.ID() {
			fd.logln("crashed", cn)
			// emulating node crash by not pinging and thus not updating the alive set.
			return
		}
	}

	if node.SinglePaxosClient == nil {
		// this happends if nodes aren't initialized properly, or used for testing
		return
	}

	fd.logln("pinging", node)
	ctx, cancel := context.WithTimeout(context.Background(), fd.delay)
	defer cancel()
	_, err := node.SinglePaxosClient.Ping(ctx, &Heartbeat{})
	if err != nil {
		log.Printf("failed to ping: %v\n", err)
		return
	}

	// the node returned the heartbeat without error in a timely manner.
	fd.m.Lock()
	fd.alive[node] = true
	fd.m.Unlock()
}
