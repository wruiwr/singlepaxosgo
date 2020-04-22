package singlepaxos

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/selabhvl/singlepaxos/reader"
	"golang.org/x/net/context"
)

const clientPrefix = "[Client(TestSystem)            ]"

func TestSystem(t *testing.T) {

	t.Logf("%v (system size=%d, quorum size=%d)", systemTest.Name, systemTest.Configuration.SystemSize, systemTest.Configuration.QuorumSize)

	// create server addresses
	var addrs []string
	for _, id := range *systemTest.Configuration.ServerIDs {
		addrs = append(addrs, "localhost:"+strconv.Itoa(id))
	}
	quorumSize := (len(addrs)-1)/2 + 1
	fmt.Printf("%s Server addresses (quorum size=%d): %v\n", clientPrefix, quorumSize, addrs)

	//****************************************
	// initialize and start the adapter connector
	leaderChan := make(chan *Node, 8)
	successChan := make(chan int, 8)
	adapterConnector := NewAdapterConnector(leaderChan, successChan)
	adapterConnector.start()

	//****************************************
	// start Paxos system

	for _, id := range *systemTest.Configuration.ServerIDs {
		go ServerStart(id, addrs, quorumSize)
	}

	//****************************************
	// start testing
	testResponseChan := make(chan *Response, 3)

	for _, test := range systemTest.TestCases {

		go func(test *reader.SystemTestcases) {

			fmt.Printf("%s test case=%s\n", clientPrefix, test.CaseID)

			var phaseOneFailLeaders []string
			var phaseTwoFailLeaders []string

			if test.TestValues.P1Failures != nil {
				fmt.Printf("%s failed leaders=%s in phase one\n", clientPrefix, *test.TestValues.P1Failures)
				phaseOneFailLeaders = make([]string, len(*test.TestValues.P1Failures))
				for n, phaseOneFailLeader := range *test.TestValues.P1Failures {
					phaseOneFailLeaders[n] = phaseOneFailLeader
				}
			}

			if test.TestValues.P2Failures != nil {
				fmt.Printf("%s failed leader=%s in phase two\n", clientPrefix, *test.TestValues.P2Failures)
				phaseTwoFailLeaders = make([]string, len(*test.TestValues.P2Failures))
				for n, phaseTwoFailLeader := range *test.TestValues.P2Failures {
					phaseTwoFailLeaders[n] = phaseTwoFailLeader
				}
			}

			// construct a slice for expected leaders by using the leaders in test oracle
			expectedLeaders := make([]string, len(*test.TestOracles.ExpectLeaders))
			for n, expectedLeader := range *test.TestOracles.ExpectLeaders {
				expectedLeaders[n] = expectedLeader
			}

			var oldLeader *Node
			n := 0
			for {
				select {
				case leader := <-leaderChan:
					fmt.Printf("%s new leader Node=%s\n", clientPrefix, leader)
					t.Logf("Test adapter got a new leader=%s from paxos system\n", leader)

					if oldLeader == nil {
						if leader.Port() != expectedLeaders[n] {
							t.Errorf("test: %v\nwant expected fitst leader: %v\ngot leader: %v",
								systemTest.Name, expectedLeaders[n], leader)
						}
						n++
					} else if len(expectedLeaders) > 1 && oldLeader.Port() == expectedLeaders[n-1] {
						if leader.Port() != expectedLeaders[n] {
							t.Errorf("test: %v\nwant expected leader: %v\ngot leader: %v",
								systemTest.Name, expectedLeaders[n], leader)
						}
						n++
					}

					oldLeader = leader

					crashedNodeMutex.Lock()
					if len(crashedNodes) < len(*test.TestOracles.ExpectLeaders)-1 {
						crashedNodes = append(crashedNodes, oldLeader)
					}
					crashedNodeMutex.Unlock()

					fmt.Printf("%s crashedNodes=%s\n", clientPrefix, crashedNodes)
					fmt.Printf("%s oldleaderNode=%s\n", clientPrefix, oldLeader)
				case successIn := <-successChan:
					fmt.Printf("%s new sucess message=%d\n", clientPrefix, successIn)
					t.Logf("Test adapter got a new message=%v from paxos system\n", successIn)

					switch successIn {
					case 1:
						switch len(*test.TestOracles.ExpectLeaders) {
						case 1:
							adapterConnector.deliverFailurePhaseOneInfor(false)
						case 2:
							if phaseOneFailLeaders != nil && phaseOneFailLeaders[0] == oldLeader.Port() {
								adapterConnector.deliverFailurePhaseOneInfor(true)
							} else {
								adapterConnector.deliverFailurePhaseOneInfor(false)
							}
						case 3:
							switch len(phaseOneFailLeaders) {
							case 0:
								adapterConnector.deliverFailurePhaseOneInfor(false)
							case 1:
								if phaseOneFailLeaders[0] == oldLeader.Port() {
									adapterConnector.deliverFailurePhaseOneInfor(true)
								} else {
									adapterConnector.deliverFailurePhaseOneInfor(false)
								}
							case 2:
								if phaseOneFailLeaders[0] == oldLeader.Port() {
									adapterConnector.deliverFailurePhaseOneInfor(true)
								} else if phaseOneFailLeaders[1] == oldLeader.Port() {
									adapterConnector.deliverFailurePhaseOneInfor(true)
								} else {
									adapterConnector.deliverFailurePhaseOneInfor(false)
								}
							}
						}
					case 2:
						switch len(*test.TestOracles.ExpectLeaders) {
						case 1:
							adapterConnector.deliverFailurePhaseTwoInfor(false)
						case 2:
							if phaseTwoFailLeaders != nil && phaseTwoFailLeaders[0] == oldLeader.Port() {
								adapterConnector.deliverFailurePhaseTwoInfor(true)
							} else {
								adapterConnector.deliverFailurePhaseTwoInfor(false)
							}
						case 3:
							switch len(phaseTwoFailLeaders) {
							case 0:
								adapterConnector.deliverFailurePhaseTwoInfor(false)
							case 1:
								if phaseTwoFailLeaders[0] == oldLeader.Port() {
									adapterConnector.deliverFailurePhaseTwoInfor(true)
								} else {
									adapterConnector.deliverFailurePhaseTwoInfor(false)
								}
							case 2:
								if phaseTwoFailLeaders[0] == oldLeader.Port() {
									adapterConnector.deliverFailurePhaseTwoInfor(true)
								} else if phaseTwoFailLeaders[1] == oldLeader.Port() {
									adapterConnector.deliverFailurePhaseTwoInfor(true)
								} else {
									adapterConnector.deliverFailurePhaseTwoInfor(false)
								}
							}
						}
					}
				}
			}
		}(test)

		// collect request from test cases
		requests := make([]string, len(*test.TestValues.ClientRequests))

		for n, req := range *test.TestValues.ClientRequests {
			requests[n] = req
		}

		if len(requests) > 1 {
			// the execution of the client two
			go func() {
				resp := testExecutor(t, &addrs, requests[1], test.TestOracles.ExpectedLegalResponses)
				fmt.Println("responseTwo:", resp)
				testResponseChan <- resp
			}()
		}

		// the execution of the client one
		responseOne := testExecutor(t, &addrs, requests[0], test.TestOracles.ExpectedLegalResponses)
		fmt.Println("responseOne:", responseOne)

		if len(requests) > 1 {
			// test if the responses are the same for two clients
			if responseTwo := <-testResponseChan; responseTwo.GetClientResponse() != responseOne.GetClientResponse() {
				t.Errorf("test: %v, client request: %v, got response: %v, but client request: %v, got response: %v",
					systemTest.Name, requests[0], responseOne.GetClientResponse(), requests[1],
					responseTwo.GetClientResponse())
			}
		}

	}
}

// testExecutor starts the client to send request to Paxos replicas and obtain response,
// it can also test if the obtained response belongs to expected legal responses.
func testExecutor(t *testing.T, addrs *[]string, clientRequest string, expectedLegalResponses *[]string) *Response {
	resp := clientStart(*addrs, &clientRequest, systemTest.Configuration.QuorumSize)

	// test if the response belongs to expected legal responses
	isOK := contain(resp.GetClientResponse(), *expectedLegalResponses)
	if !isOK {
		t.Errorf("test: %v, client request: %v, got response: %v, not belong to expected responses: %v",
			systemTest.Name, clientRequest, resp.GetClientResponse(), *expectedLegalResponses)
	} else {
		t.Logf("test: %v, client request: %v, got response: %v, that belong to expected responses: %v",
			systemTest.Name, clientRequest, resp.GetClientResponse(), *expectedLegalResponses)
	}

	return resp
}

// Internal: contain checks if the obtained response is within the expected response list.
func contain(v string, expectVal []string) bool {
	for _, val := range expectVal {
		if v == val {
			return true
		}
	}
	return len(expectVal) == 0 && len(v) == 0
}

// Internal: clientStart starts a client to send request to Paxos replicas by invoking a quorum call.
func clientStart(addrs []string, clientRequest *string, quorumSize int) *Response {
	fmt.Printf("%s Connecting to %d Paxos replicas: %v\n", clientPrefix, len(addrs), addrs)

	// setup the client connection to the Paxos replicas
	config, mgr := newPaxosConfig(addrs, quorumSize)
	defer mgr.Close()

	req := Value{ClientRequest: *clientRequest}
	resp := doSendRequest(config, req)
	fmt.Println(clientPrefix, "Response:", resp.ClientResponse, "for the client request:", req)

	return resp
}

// Internal: doSendRequest can send requests to paxos servers by quorum call
func doSendRequest(config *Configuration, value Value) *Response {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	defer cancel()
	resp, err := config.ClientHandle(ctx, &value)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	return resp
}
