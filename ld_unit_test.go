package singlepaxos

import (
	"strconv"
	"testing"
	"time"
)

func TestLeaderChange(t *testing.T) {
	t.Logf("%v (system size=%d, quorum size=%d)", leaderChangeTest.Name, leaderChangeTest.SystemSize, leaderChangeTest.QuorumSize)

	var testNodes = make([]*Node, leaderChangeTest.SystemSize)

	for i, test := range leaderChangeTest.TestCases {
		t.Logf("test case=%v", test.CaseID)

		for j, serverID := range *test.TestValues.ServerIDs {
			testNodes[j] = &Node{id: uint32(serverID), addr: "localhost:" + strconv.Itoa(serverID)}
		}

		config := &Configuration{id: testNodes[0].ID(), nodes: testNodes, n: leaderChangeTest.SystemSize}

		// initialize leader detector
		ld := NewLeaderDetector(*config, nil, 10*time.Second)

		for n, leader := range *test.TestOracles.ExpectLeaders {
			switch n {
			case 0:
				// check the first leader
				if ld.leader.ID() != uint32(leader) {
					t.Errorf("TestLeaderChange: Test %d\ngot first leader %v, want first leader %v",
						i+1, ld.leader.ID(), leader)
				}
				// suspect first leader
				ld.Suspect(ld.leader)
			case 1:
				// check the second leader
				if ld.leader.ID() != uint32(leader) {
					t.Errorf("TestLeaderChange: Test %d\ngot second leader %v, want second leader %v",
						i+1, ld.leader.ID(), leader)
				}
				// suspect second leader
				ld.Suspect(ld.leader)
			case 2:
				// check the third leader
				if ld.leader.ID() != uint32(leader) {
					t.Errorf("TestLeaderChange: Test %d\ngot third leader %v, want third leader %v",
						i+1, ld.leader.ID(), leader)
				}
			}
		}
	}
}
