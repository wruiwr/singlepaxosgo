package singlepaxos

import "testing"

// TestPrepareQFUnitTest is the unit test for PrepareQF.
func TestPrepareQFUnitTest(t *testing.T) {
	t.Logf("%v (system size=%d, quorum size=%d)", prepareQFTest.Name, prepareQFTest.Configuration.SystemSize, prepareQFTest.Configuration.QuorumSize)

	qspec := NewPaxosQSpec(prepareQFTest.Configuration.QuorumSize)
	for i, test := range prepareQFTest.TestCases {

		t.Logf("test case ID=%v", test.CaseID)
		var replies []*PromiseMsg
		// emulate the a prepare was sent with the promise.Rnd of the first entry.
		prepare := &PrepareMsg{Rnd: test.TestValues.PromiseMsgs[0].Rnd}

		for _, prmMsg := range test.TestValues.PromiseMsgs {
			// create promise message according to the value of the test case.
			vval := Value{prmMsg.Vval}
			promise := PromiseMsg{Rnd: prmMsg.Rnd, Vrnd: prmMsg.Vrnd, Vval: &vval}
			// collect replies according to the value of the test case.
			replies = append(replies, &promise)
		}

		// invoke PrepareQF
		prm, quorum := qspec.PrepareQF(prepare, replies)
		// compare results with oracles
		switch {
		case !test.TestOracles.ExpectQuorum && quorum:
			t.Errorf("test: %d\nwant no quorum\ngot: %v",
				i+1, prm)
		case test.TestOracles.ExpectQuorum && !quorum:
			t.Errorf("test: %d\nwant: %v\ngot no quorm",
				i+1, test.TestOracles.ExpectResult)
		case test.TestOracles.ExpectQuorum && quorum:
			if prm.Rnd != test.TestOracles.ExpectResult.Rnd {
				t.Errorf("test: %d\nwant promise rnd: %v\ngot promise rnd: %v",
					i+1, test.TestOracles.ExpectResult.Rnd, prm.Rnd)
			}
			if prm.Vrnd != test.TestOracles.ExpectResult.Vrnd {
				t.Errorf("test: %d\nwant promise vrnd: %v\ngot promise vrnd: %v",
					i+1, test.TestOracles.ExpectResult.Vrnd, prm.Vrnd)
			}
			// get expected vval according to the value of the test case.
			expectVval := Value{test.TestOracles.ExpectResult.Vval}
			if *prm.Vval != expectVval {
				t.Errorf("test: %d\nwant promise vval: %v\ngot promise vval: %v",
					i+1, expectVval, prm.Vval)
			}
		}
	}
}

// TestAcceptQFUnitTest is the unit test for AcceptQF.
func TestAcceptQFUnitTest(t *testing.T) {
	t.Logf("%v (system size=%d, quorum size=%d)", acceptQFTest.Name, acceptQFTest.Configuration.SystemSize, acceptQFTest.Configuration.QuorumSize)

	qspec := NewPaxosQSpec(acceptQFTest.Configuration.QuorumSize)
	for i, test := range acceptQFTest.TestCases {

		t.Logf("test case ID=%v", test.CaseID)
		var replies []*LearnMsg

		for _, lrnMsg := range test.TestValues.LearnMsgs {
			// create learn message according to the value of test case.
			val := Value{lrnMsg.Val}
			learn := LearnMsg{Rnd: lrnMsg.Rnd, Val: &val}
			// collect replies according to the value of the test case.
			replies = append(replies, &learn)
		}

		// invoke AcceptQF
		lrn, quorum := qspec.AcceptQF(replies)
		gotVal := lrn.GetVal()
		// compare results with oracles
		switch {
		case !test.TestOracles.ExpectQuorum && quorum:
			t.Errorf("test: %d\nwant no quorum\ngot: %v",
				i+1, gotVal)
		case test.TestOracles.ExpectQuorum && !quorum:
			t.Errorf("test: %d\nwant: %v\ngot no quorm",
				i+1, test.TestOracles.ExpectResult.Val)
		case test.TestOracles.ExpectQuorum && quorum:
			expectVal := Value{test.TestOracles.ExpectResult.Val}
			if gotVal.GetClientRequest() != expectVal.GetClientRequest() {
				t.Errorf("test: %d\nwant: %v\ngot: %v",
					i+1, expectVal.GetClientRequest(), gotVal.GetClientRequest())
			}
		}
	}
}
