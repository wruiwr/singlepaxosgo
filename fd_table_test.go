package singlepaxos

import (
	"reflect"
	"testing"
	"time"
)

const ourID = 8080

var (
	nodeOne   = Node{id: 8080, addr: "localhost:8080"}
	nodeTwo   = Node{id: 8081, addr: "localhost:8081"}
	nodeThree = Node{id: 8082, addr: "localhost:8082"}
	testNodes = []*Node{&nodeOne, &nodeTwo, &nodeThree}
)

func TestAllNodesShouldBeAlivePreStart(t *testing.T) {

	config := &Configuration{id: ourID, nodes: testNodes, n: len(testNodes)}
	ld := NewLeaderDetector(*config, nil, 10*time.Second)
	fd := NewFailureDetector(*config, ld, nil, 10*time.Second)

	if len(fd.alive) != len(testNodes) {
		t.Errorf("TestAllNodesShouldBeAlivePreStart: alive set contains %d node ids, want %d", len(fd.alive), len(testNodes))
	}

	for _, node := range testNodes {
		alive := fd.alive[node]
		if !alive {
			t.Errorf("TestAllNodesShouldBeAlivePreStart: node %v was not set alive", node)
			continue
		}
	}
}

var timeoutTests = []struct {
	alive             map[*Node]bool
	suspected         map[*Node]bool
	wantPostSuspected map[*Node]bool
	wantDelay         time.Duration
}{
	{
		alive:             map[*Node]bool{&nodeOne: true, &nodeTwo: true, &nodeThree: true},
		suspected:         map[*Node]bool{},
		wantPostSuspected: map[*Node]bool{},
		wantDelay:         10 * time.Second,
	},
	{
		alive:             map[*Node]bool{&nodeOne: true, &nodeThree: true},
		suspected:         map[*Node]bool{&nodeTwo: true},
		wantPostSuspected: map[*Node]bool{&nodeTwo: true},
		wantDelay:         10 * time.Second,
	},
	{
		alive:             map[*Node]bool{},
		suspected:         map[*Node]bool{},
		wantPostSuspected: map[*Node]bool{&nodeOne: true, &nodeTwo: true, &nodeThree: true},
		wantDelay:         10 * time.Second,
	},
	{
		alive:             map[*Node]bool{&nodeThree: true},
		suspected:         map[*Node]bool{},
		wantPostSuspected: map[*Node]bool{&nodeOne: true, &nodeTwo: true},
		wantDelay:         10 * time.Second,
	},
	{
		alive:             map[*Node]bool{},
		suspected:         map[*Node]bool{&nodeOne: true, &nodeTwo: true, &nodeThree: true},
		wantPostSuspected: map[*Node]bool{&nodeOne: true, &nodeTwo: true, &nodeThree: true},
		wantDelay:         10 * time.Second,
	},
}

func TestTimeoutProcedure(t *testing.T) {
	for i, test := range timeoutTests {
		config := &Configuration{id: ourID, nodes: testNodes, n: len(testNodes)}
		ld := NewLeaderDetector(*config, nil, 10*time.Second)
		fd := NewFailureDetector(*config, ld, nil, 10*time.Second)

		// Set our test data
		fd.alive = test.alive
		fd.suspected = test.suspected

		// Trigger timeout procedure
		fd.timeout()

		// Alive set should always be empty
		if len(fd.alive) > 0 {
			t.Errorf("TestTimeoutProcedure %d: Alive set should always be empty after timeout procedure completes, has length %d", i, len(fd.alive))
		}

		if !reflect.DeepEqual(test.wantPostSuspected, fd.suspected) {
			t.Errorf("TestTimeoutProcedure %d: suspected set post timeout procedure differs", i)
			printSuspectedDiff(t, fd.suspected, test.wantPostSuspected)
		}

		// Check delay
		if fd.delay != test.wantDelay {
			t.Errorf("TestTimeoutProcedure %d: want %v delay after timeout procedure, got %v", i, test.wantDelay, fd.delay)
		}

	}
}

func printSuspectedDiff(t *testing.T, got, want map[*Node]bool) {
	t.Errorf("-----------------------------------------------------------------------------")
	t.Errorf("Got:")
	if len(got) == 0 {
		t.Errorf("None")
	}
	for node := range got {
		t.Errorf("suspect %v", node)
	}
	t.Errorf("Want:")
	if len(want) == 0 {
		t.Errorf("None")
	}
	for node := range want {
		t.Errorf("suspect %v", node)
	}
	t.Errorf("-----------------------------------------------------------------------------")
}
