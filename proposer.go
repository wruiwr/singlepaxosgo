package singlepaxos

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/context"
)

// Proposer represents the state of a single-decree Paxos proposer.
type Proposer struct {
	m         sync.RWMutex   // lock to prevent concurrent modification
	node      *Node          // the node running this Proposer instance
	leader    *Node          // the leader node
	config    *Configuration // the configuration holds the PI set (all nodes)
	crnd      uint32         // current round number
	cval      *Value         // current value, if set
	trustMsgs <-chan *Node   // channel on which to receive leader changes
	cvalIn    chan bool      // channel used to trigger consensus when received a new cval from client
	stop      chan struct{}

	electedLeader *Node
}

// NewProposer returns a new single-decree Paxos proposer.
func NewProposer(c *Configuration, port int) *Proposer {
	var myNode *Node
	for _, node := range c.Nodes() {
		if node.Port() == strconv.Itoa(port) {
			myNode = node
			break
		}
	}
	if myNode == nil {
		panic("my node not found in configuration")
	}

	// initialize leader detector with a time delay for failure detector.
	ld := NewLeaderDetector(*c, myNode, 10*time.Second)

	return &Proposer{
		node:      myNode,
		config:    c,
		crnd:      uint32(port - c.Size()), // initialize round number
		trustMsgs: ld.Subscribe(),
		cvalIn:    make(chan bool, 1),
		stop:      make(chan struct{}),
	}
}

func (p *Proposer) logf(f string, a ...interface{}) {
	format := "[P.%s, crnd=%d, leader=%s] " + f
	leaderPort := p.leader.Port()
	if leaderPort == "" {
		leaderPort = "____"
	}
	fmt.Printf(format, p.node.Port(), p.crnd, leaderPort, a)
}

func (p *Proposer) Start() {
	p.logf("%s\n", "Starting Proposer")
	go func() {
		for {
			select {
			case <-p.cvalIn:
				p.logf("Received client request %v in for loop of Start method\n", p.cval)
				p.checkAndRunPaxos()
			case p.leader = <-p.trustMsgs:
				p.logf("%s : %v\n ", "Received new leader", p.leader)
				p.checkAndRunPaxos()
			case <-p.stop:
				p.logf("%s\n", "Stopped proposer of replica")
				return
			}
		}
	}()
}
func (p *Proposer) isLeaderAndHaveValue() bool {
	p.m.RLock()
	defer p.m.RUnlock()
	return p.leader == p.node && p.cval != nil && p.leader != p.electedLeader
}

func (p *Proposer) checkAndRunPaxos() {
	if p.isLeaderAndHaveValue() {
		p.electedLeader = p.leader // set the received leader as the elected leader
		// only the proposer who is leader should send the 'leader node' to the test adapter;
		// other proposers should not send the leader node.
		LeaderChan <- p.leader // send leader to test adapter connector

		p.increaseCrnd()
		p.logf("I'm the current leader and I've got a value %v\n", p.cval)
		err := p.runPaxosPhases()
		if err != nil {
			// It is unclear what can cause this execution path to be taken,
			// except that a majority of the remaining nodes have failed, e.g.
			// if the leader (this node) is alone in a 3-member configuration.
			p.logf("leader failure %v\n", err)
		}
	}
}

// runPaxosPhases executes three Paxos phases
func (p *Proposer) runPaxosPhases() error {
	// access proposer state in mutual exclusion for use below; avoid holding lock during quorum calls
	p.m.RLock()
	crnd, cval := p.crnd, p.cval
	p.m.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// ******************************************************
	//  PHASE ONE: send Prepare to obtain quorum of Promises
	preMsg := &PrepareMsg{Rnd: crnd}
	p.logf("Sending  Phase 1a msg: %v\n", preMsg)
	prmMsg, err := p.config.Prepare(ctx, preMsg)
	err = notifyTestAdapter(err, 1, FailurePhaseOneChan)
	if err != nil {
		return err
	}
	p.logf("Received Phase 1b msg: %v\n", prmMsg)

	// ******************************************************
	//  PHASE TWO: send Accept to obtain quorum of Learns
	if prmMsg.GetVrnd() != Ignore {
		// promise msg has a locked-in value
		cval = prmMsg.GetVval()
		// update proposer state in mutual exclusion
		p.m.Lock()
		p.cval = cval
		p.m.Unlock()
	}
	// use local proposer's cval or locked-in value from promise msg, if any.
	accMsg := &AcceptMsg{Rnd: crnd, Val: cval}
	p.logf("Sending  Phase 2a msg: %v\n", accMsg)
	lrnMsg, err := p.config.Accept(ctx, accMsg)
	err = notifyTestAdapter(err, 2, FailurePhaseTwoChan)
	if err != nil {
		return err
	}
	p.logf("Received Phase 2b msg: %v\n", lrnMsg)

	// ******************************************************
	//  PHASE THREE: send Commit to obtain a quorum of Acks
	p.logf("Sending  Phase 3a msg: %v\n", lrnMsg)
	ackMsg, err := p.config.Commit(ctx, lrnMsg)
	if err != nil {
		return err
	}
	p.logf("Received Phase 3b msg: %v\n", ackMsg)

	return nil
}

// Stop stops p's main run loop.
func (p *Proposer) Stop() {
	p.stop <- struct{}{}
}

// DeliverClientValue delivers client value cval to proposer.
func (p *Proposer) ProposeClientValue(cval *Value) {
	p.logf("Received client request %v\n", cval)
	p.m.Lock()

	if p.cval == nil {
		p.cval = cval
		p.logf("Updated proposer cval: %v\n", p.cval)
	}
	p.m.Unlock()
	// notify event loop that a client value exists
	p.cvalIn <- true
}

// increaseCrnd increases proposer p's crnd field by the total number
// of Paxos nodes.
func (p *Proposer) increaseCrnd() {
	p.m.Lock()
	p.crnd = p.crnd + uint32(p.config.Size())
	p.m.Unlock()
}
