package singlepaxos

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	t.Logf("%v (system size=%d, quorum size=%d)", timeoutTest.Name, timeoutTest.SystemSize, timeoutTest.QuorumSize)

	for i, test := range timeoutTest.TestCases {
		t.Logf("test case=%v", test.CaseID)

		var testNodes = make([]*Node, leaderChangeTest.SystemSize)
		var expectSuspects = make(map[*Node]bool, leaderChangeTest.SystemSize)

		for j, serverID := range *test.TestValues.ServerIDs {
			testNodes[j] = &Node{id: uint32(serverID), addr: "localhost:" + strconv.Itoa(serverID)}
		}

		t.Logf("test nodes: %v", testNodes)

		if test.TestOracles.ExpectSuspects != nil {
			for n, expectSuspected := range *test.TestOracles.ExpectSuspects {
				if expectSuspected {
					expectSuspects[testNodes[n]] = expectSuspected
				} else {
					delete(expectSuspects, testNodes[n])
				}
			}
		}

		config := &Configuration{id: testNodes[0].ID(), nodes: testNodes, n: leaderChangeTest.SystemSize}

		ld := NewLeaderDetector(*config, nil, 10*time.Second)
		fd := NewFailureDetector(*config, ld, nil, 10*time.Second)

		if test.TestValues.Alives != nil {
			for n, alive := range *test.TestValues.Alives {
				if alive {
					fd.alive[testNodes[n]] = alive
				} else {
					delete(fd.alive, testNodes[n])
				}

			}
		}

		if test.TestValues.Suspects != nil {
			for n, suspected := range *test.TestValues.Suspects {
				if suspected {
					fd.suspected[testNodes[n]] = suspected
				} else {
					delete(fd.suspected, testNodes[n])
				}
			}
		}

		// Trigger timeout procedure
		fd.timeout()

		// Alive set should always be empty
		if len(fd.alive) > 0 {
			t.Errorf("TestTimeoutProcedure %d: Alive set should always be empty after timeout procedure completes, has length %d", i, len(fd.alive))
		}

		if !reflect.DeepEqual(expectSuspects, fd.suspected) {
			t.Errorf("TestTimeoutProcedure %d: suspected set post timeout procedure differs", i)
			printSuspectedDiff(t, fd.suspected, expectSuspects)
		}
	}
}
