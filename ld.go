package singlepaxos

import (
	"sync"
	"time"
)

// LeaderDetector
type LeaderDetector struct {
	m           sync.Mutex
	config      Configuration  // the configuration holds the PI set (all nodes)
	leader      *Node          // current leader (initially not set)
	suspected   map[*Node]bool // map of node ids considered suspected
	subscribers []chan *Node   // subscribers interested in notifications about leader changes
}

func NewLeaderDetector(c Configuration, myNode *Node, delay time.Duration) *LeaderDetector {
	ld := &LeaderDetector{
		config:      c,
		suspected:   make(map[*Node]bool),
		subscribers: make([]chan *Node, 0),
	}
	// set the initial leader
	ld.leader = ld.minRank()
	// create failure detector to be associated with this leader detector
	NewFailureDetector(c, ld, myNode, delay)
	return ld
}

// Suspect notifies the LeaderDetector that the given
// node is suspected to have failued.
func (ld *LeaderDetector) Suspect(node *Node) {
	ld.m.Lock()
	defer ld.m.Unlock()
	ld.suspected[node] = true
	if ld.leader == node {
		ld.leader = ld.minRank()
		// if leader is not nil, publish the new leader
		if ld.leader != nil {
			ld.publish()
		}
	}
}

// Restore notifies the LeaderDetector that the given
// previously suspected node, is no longer suspected.
func (ld *LeaderDetector) Restore(node *Node) {
	ld.m.Lock()
	defer ld.m.Unlock()
	delete(ld.suspected, node)
}

// Subscribe returns a buffered channel on which leader changes are published.
func (ld *LeaderDetector) Subscribe() <-chan *Node {
	ld.m.Lock()
	defer ld.m.Unlock()
	ch := make(chan *Node, 8)
	ld.subscribers = append(ld.subscribers, ch)
	// TODO I don't like to use Sleep here, but it seems that publish() too soon is problematic
	// for the Proposer, if it hasn't started yet. Not totally sure why this is a problem.
	time.Sleep(30 * time.Millisecond)
	// ensure that the new subscriber gets notified of the current leader right away;
	// we don't want to wait for the first failure.

	// if leader is not nil, publish the new leader
	if ld.leader != nil {
		ld.publish()
	}
	return ch
}

// publish sends the the leader id onto every subscriber channel.
// Should only be called when holding the LeaderDetector mutex.
func (ld *LeaderDetector) publish() {
	for _, subscriber := range ld.subscribers {
		select {
		case subscriber <- ld.leader:
			// Send success.
		default:
			// Drop publication.
			// Receviers buffer is full.
		}
	}
}

// minRank returns the lowest ranking unsuspected node in the configuration.
// Should only be called when holding the LeaderDetector mutex.
func (ld *LeaderDetector) minRank() *Node {

	// check if there are no nodes, then return nil.
	if len(ld.config.nodes) == 0 {
		return nil
	}

	leaderCandidates := make([]*Node, 0)
	for _, id := range ld.config.Nodes() {
		if suspected := ld.suspected[id]; suspected {
			continue
		}
		leaderCandidates = append(leaderCandidates, id)
	}

	if len(leaderCandidates) == 0 {
		return nil
	}
	// sort the list of the leader candidates
	OrderedBy(Port).Sort(leaderCandidates)
	// PS: If there are no live servers (or broken connection), this will cause index out of bounds panic.
	return leaderCandidates[0]
}
