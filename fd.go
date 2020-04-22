package singlepaxos

import (
	"sync"
	"time"

	"fmt"
)

// Suspecter is the interface that wraps the Suspect method. Suspect indicates
// that the node with identifier id should be considered suspected.
type Suspecter interface {
	Suspect(node *Node)
}

// Restorer is the interface that wraps the Restore method. Restore indicates
// that the node with identifier id should be considered restored.
type Restorer interface {
	Restore(node *Node)
}

// SuspectRestorer is the interface that groups the Suspect and Restore
// methods.
type SuspectRestorer interface {
	Suspecter
	Restorer
}

// failureDetector holds the state of the failure detector.
type failureDetector struct {
	m sync.Mutex

	// failure detector fields
	config    Configuration   // the configuration holds the PI set (all nodes)
	alive     map[*Node]bool  // map of nodes considered alive
	suspected map[*Node]bool  // map of nodes considered suspected
	sr        SuspectRestorer // SuspectRestorer to notify about suspects and restores
	node      *Node           // my local node

	delay         time.Duration // the current delay for the timeout procedure
	delta         time.Duration // the delta value to be used when increasing delay
	timeoutSignal *time.Ticker  // the timeout procedure ticker
	hbChan        chan *Node    // the heartbeat channel of nodes to ping)
	stop          chan struct{}
}

func NewFailureDetector(c Configuration, sr SuspectRestorer, myNode *Node, delta time.Duration) *failureDetector {
	alive := make(map[*Node]bool)
	for _, id := range c.Nodes() {
		alive[id] = true
	}
	fd := &failureDetector{
		config:        c,
		alive:         alive,
		suspected:     make(map[*Node]bool),
		sr:            sr,
		node:          myNode,
		delay:         delta,
		delta:         delta,
		timeoutSignal: time.NewTicker(delta),
		hbChan:        make(chan *Node, 10*c.Size()),
	}
	fd.start()
	return fd
}

func (fd *failureDetector) logf(str string, a ...interface{}) {
	fmt.Printf("[FD.%s                       ] "+str, fd.node.Port(), a)
}

func (fd *failureDetector) logln(a ...interface{}) {
	fmt.Printf("[FD.%s                       ] %v\n", fd.node.Port(), a)
}

// Start handles timeout procedure when failure happens
func (fd *failureDetector) start() {
	go func() {
		for {
			select {
			case node := <-fd.hbChan:
				go fd.ping(node)
			case <-fd.timeoutSignal.C:
				fd.timeout()
				// reset timer
				fd.timeoutSignal = time.NewTicker(fd.delay)
			case <-fd.stop:
				fmt.Println("stoppted fd.")
				return
			}
		}
	}()
}

// Stop stops ld's main run loop.
func (fd *failureDetector) Stop() {
	fd.stop <- struct{}{}
}

// timeout checks the liveness status of the replicas and updates the suspect-restore listener.
func (fd *failureDetector) timeout() {
	fd.logln("timeout")
	fd.m.Lock()
	defer fd.m.Unlock()
	if !fd.aliveSuspectedIntersectionEmpty() {
		fd.delay = fd.delay + fd.delta
		fd.logf("new delay %d\n", fd.delay)
		fd.timeoutSignal = time.NewTicker(fd.delay)
	}
	for _, node := range fd.config.Nodes() {
		if !fd.alive[node] && !fd.suspected[node] {
			fd.suspected[node] = true
			fd.logf("suspect %v\n", node)
			fd.sr.Suspect(node)
		} else if fd.alive[node] && fd.suspected[node] {
			delete(fd.suspected, node)
			fd.logf("restore %v\n", node)
			fd.sr.Restore(node)
		}

		fd.hbChan <- node
	}
	fd.logln("fd.alive", fd.alive)
	fd.alive = make(map[*Node]bool)
}

// Internal helper method.
func (fd *failureDetector) aliveSuspectedIntersectionEmpty() bool {
	for node := range fd.suspected {
		if fd.alive[node] {
			return false
		}
	}
	return true
}
