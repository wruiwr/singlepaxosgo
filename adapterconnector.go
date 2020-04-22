package singlepaxos

import (
	"fmt"
	"sync"
)

const adapterPrefix = "[TestAdapter                   ]"

var (
	crashedNodes     = make([]*Node, 0) // note: it needs to reset if multiple test cases in one test file
	crashedNodeMutex sync.Mutex
)

var (
	LeaderChan          = make(chan *Node, 8)
	FailurePhaseOneChan = make(chan bool, 8)
	FailurePhaseTwoChan = make(chan bool, 8)
	SuccessPhaseChan    = make(chan int, 8)
)

type AdapterConnector struct {
	leaderChanOut    chan *Node
	successChanOut   chan int
	failureOneChanIn chan bool
	failureTwoChanIn chan bool
}

func NewAdapterConnector(leaderChanOut chan *Node, successChanOut chan int) *AdapterConnector {
	return &AdapterConnector{
		leaderChanOut:    leaderChanOut,
		successChanOut:   successChanOut,
		failureOneChanIn: make(chan bool, 8),
		failureTwoChanIn: make(chan bool, 8),
	}
}

func (ac *AdapterConnector) start() {
	go func() {
		for {
			select {
			case leader := <-LeaderChan:
				fmt.Printf("%s Got leader=%v\n", adapterPrefix, leader)
				ac.sendLeader(leader)
			case success := <-SuccessPhaseChan:
				fmt.Printf("%s Got success=%v\n", adapterPrefix, success)
				ac.sendSuccessInfor(success)
			case failurePhaseOneInfor := <-ac.failureOneChanIn:
				fmt.Printf("%s Got failure=%v in Phase One\n", adapterPrefix, failurePhaseOneInfor)
				FailurePhaseOneChan <- failurePhaseOneInfor
			case failurePhaseTwoInfor := <-ac.failureTwoChanIn:
				fmt.Printf("%s Got failure=%v in Phase Two\n", adapterPrefix, failurePhaseTwoInfor)
				FailurePhaseTwoChan <- failurePhaseTwoInfor
			}
		}
	}()
}

func (ac *AdapterConnector) sendLeader(leader *Node) {
	ac.leaderChanOut <- leader
}

func (ac *AdapterConnector) sendSuccessInfor(success int) {
	ac.successChanOut <- success
}

func (ac *AdapterConnector) deliverFailurePhaseOneInfor(failureInfor bool) {
	ac.failureOneChanIn <- failureInfor
}

func (ac *AdapterConnector) deliverFailurePhaseTwoInfor(failureInfor bool) {
	ac.failureTwoChanIn <- failureInfor
}

func notifyTestAdapter(err error, phase int, failurePhaseChan chan bool) error {
	// notify test adapter for phase
	if err == nil {
		SuccessPhaseChan <- phase
		fmt.Printf("%s Phase %d success sent to test adapter\n", adapterPrefix, phase)
	}
	// waiting to obtain decision from test adapter connector for current phase
	if failurePhase := <-failurePhaseChan; failurePhase {
		return fmt.Errorf("test adapter decided to fail phase %d", phase)
	}
	// test adapter decided not to fail this phase; return original err, if any
	return err
}
