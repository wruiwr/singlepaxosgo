package singlepaxos

import (
	"testing"
	"time"
)

// test cases for testing minRank() method
var leaderTests = []struct {
	id         uint32
	nodes      []*Node
	numOfNodes int
	wantLeader *Node
}{
	{8080, []*Node{}, 0, nil},
	{8080, []*Node{}, 3, nil},
	{8080, []*Node{{id: 8080, addr: "localhost:8080"}}, 1, &Node{id: 8080, addr: "localhost:8080"}},
	{8080, []*Node{{id: 8080, addr: "localhost:8080"}, {id: 8081, addr: "localhost:8081"}}, 2, &Node{id: 8080, addr: "localhost:8080"}},
	{8080, []*Node{{id: 8080, addr: "localhost:8080"}, {id: 8081, addr: "localhost:8081"}, {id: 8082, addr: "localhost:8082"}}, 3, &Node{id: 8080, addr: "localhost:8080"}},
	{8080, []*Node{{id: 8081, addr: "localhost:8081"}, {id: 8082, addr: "localhost:8082"}}, 3, &Node{id: 8081, addr: "localhost:8081"}},
	{8080, []*Node{{id: 8082, addr: "localhost:8082"}}, 3, &Node{id: 8082, addr: "localhost:8082"}},
	{8080, []*Node{{id: 8082, addr: "localhost:8082"}, {id: 8080, addr: "localhost:8080"}, {id: 8081, addr: "localhost:8081"}}, 2, &Node{id: 8080, addr: "localhost:8080"}},
}

// TestMinRank tests minRank() method
func TestMinRank(t *testing.T) {
	for i, test := range leaderTests {
		config := &Configuration{id: test.id, nodes: test.nodes, n: test.numOfNodes}
		ld := NewLeaderDetector(*config, nil, 10*time.Second)
		leaderNode := ld.minRank()
		if test.numOfNodes == 0 || len(nodes) == 0 {
			if leaderNode != test.wantLeader {
				t.Errorf("TestGetLeader %v: got leader %v, want leader %v",
					i+1, leaderNode, test.wantLeader)
			}
		} else {
			if leaderNode.ID() != test.wantLeader.ID() {
				t.Errorf("TestGetLeader %v: got leader %v, want leader %v",
					i+1, leaderNode.ID(), test.wantLeader.ID())
			}
		}
	}
}

// test cases for testing Suspect(node *Node) method
var suspectTests = []struct {
	id         uint32
	nodes      []*Node
	numOfNodes int
	wantLeader *Node
}{
	{8080, []*Node{{id: 8080, addr: "localhost:8080"}, {id: 8081, addr: "localhost:8081"}, {id: 8082, addr: "localhost:8082"}}, 3, &Node{id: 8081, addr: "localhost:8081"}},
	{8080, []*Node{{id: 8081, addr: "localhost:8081"}, {id: 8082, addr: "localhost:8082"}}, 3, &Node{id: 8082, addr: "localhost:8082"}},
	{8080, []*Node{{id: 8082, addr: "localhost:8082"}}, 3, nil},
	{8080, []*Node{{id: 8082, addr: "localhost:8082"}, {id: 8080, addr: "localhost:8080"}, {id: 8081, addr: "localhost:8081"}}, 3, &Node{id: 8081, addr: "localhost:8081"}},
}

// TestSuspect tests Suspect(node *Node) method
func TestSuspect(t *testing.T) {
	for i, test := range suspectTests {
		config := &Configuration{id: test.id, nodes: test.nodes, n: test.numOfNodes}
		ld := NewLeaderDetector(*config, nil, 10*time.Second)
		ld.Suspect(ld.leader)
		if test.numOfNodes == 1 || len(nodes) == 1 {
			if ld.leader != test.wantLeader {
				t.Errorf("TestGetLeader %v: got leader %v, want leader %v",
					i+1, ld.leader, test.wantLeader)
			}
		} else {
			if ld.leader.ID() != test.wantLeader.ID() {
				t.Errorf("TestGetLeader %d: got leader %v, want leader %v",
					i+1, ld.leader.ID(), test.wantLeader.ID())
			}
		}

	}
}
