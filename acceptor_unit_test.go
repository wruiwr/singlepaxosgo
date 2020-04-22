package singlepaxos

import (
	"testing"
)

func TestAcceptor(t *testing.T) {
	t.Logf("%v (system size=%d, quorum size=%d)", acceptorTest.Name, acceptorTest.SystemSize, acceptorTest.QuorumSize)

	for i, test := range acceptorTest.TestCases {
		t.Logf("test case ID=%v", test.CaseID)
		// initialize an acceptor
		acceptor := NewAcceptor()
		for j, value := range test.TestValues {
			switch value.MsgType {
			case "prepare":
				// create a prepare message from test cases
				prepareMsg := PrepareMsg{Rnd: value.PrepareMsg.Rnd}

				var vval Value
				var expectPromise PromiseMsg
				if value.ExpectPromiseMsgResults != nil {
					// create a expected Value from test cases to construct an expected promise message
					vval = Value{value.ExpectPromiseMsgResults.Vval}
					// create an expected promise message from test cases
					expectPromise = PromiseMsg{Rnd: value.ExpectPromiseMsgResults.Rnd, Vrnd: value.ExpectPromiseMsgResults.Vrnd, Vval: &vval}
				}
				// invoke handlePrepare of the acceptor
				gotPrm, gotOutput := acceptor.handlePrepare(&prepareMsg)
				switch {
				case !value.ExpectOutput && gotOutput:
					t.Errorf("test nr:%d\naction nr: %d\nwant no output\ngot %v",
						i+1, j+1, gotPrm)
				case value.ExpectOutput && !gotOutput:
					t.Errorf("test nr:%d\naction nr: %d\nwant %v\ngot no output",
						i+1, j+1, value.ExpectPromiseMsgResults)
				case value.ExpectOutput && gotOutput:
					if gotPrm.Rnd != expectPromise.Rnd || gotPrm.Vrnd != expectPromise.Vrnd || *gotPrm.Vval != *expectPromise.Vval {
						t.Errorf("test nr:%d\naction nr: %d\nwant: %v\ngot: %v",
							i+1, j+1, expectPromise, *gotPrm)
					}
				}
			case "accept":
				// create a test value from test cases to construct an accept message for testing
				val := Value{value.AcceptMsg.Val}
				// create an accept message from test cases
				acceptMsg := AcceptMsg{Rnd: value.AcceptMsg.Rnd, Val: &val}

				var expectval Value
				var expectedLearn LearnMsg
				if value.ExpectLearnMsgResults != nil {
					// create an expected value from test cases to construct an expected learn message.
					expectval = Value{value.ExpectLearnMsgResults.Val}
					// create an expected learn message from test cases
					expectedLearn = LearnMsg{Rnd: value.ExpectLearnMsgResults.Rnd, Val: &expectval}
				}
				// invoke handleAccept of the acceptor
				gotLrn, gotOutput := acceptor.handleAccept(&acceptMsg)
				switch {
				case !value.ExpectOutput && gotOutput:
					t.Errorf("test nr:%d\naction nr: %d\nwant no output\ngot %v",
						i+1, j+1, gotLrn)
				case value.ExpectOutput && !gotOutput:
					t.Errorf("test nr:%d\naction nr: %d\nwant %v\ngot no output",
						i+1, j+1, expectedLearn)
				case value.ExpectOutput && gotOutput:
					if gotLrn.Rnd != expectedLearn.Rnd || *gotLrn.Val != *expectedLearn.Val {
						t.Errorf("test nr:%d\naction nr: %d\nwant: %v\ngot: %v",
							i+1, j+1, expectedLearn, gotLrn)
					}
				}
			default:
				t.Fatal("assertion failed: unkown messages type for acceptor")
			}
		}
	}
}
