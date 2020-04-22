package singlepaxos

import "testing"

func TestPrepareQF(t *testing.T) {
	qspec := NewPaxosQSpec(2)
	for i, test := range prepareQFTestcases {
		var replies []*PromiseMsg
		// emulate the a prepare was sent with the promise.Rnd of the first entry.
		prepare := &PrepareMsg{Rnd: test.actions[0].promise.GetRnd()}
		for j, action := range test.actions {
			promise := action.promise
			replies = append(replies, &promise)
			prm, quorum := qspec.PrepareQF(prepare, replies)
			switch {
			case !action.wantOutput && quorum:
				t.Errorf("test: %d, action: %d\ndescription: %s\nwant no quorum\ngot: %v",
					i+1, j+1, test.desc, prm)
			case action.wantOutput && !quorum:
				t.Errorf("test: %d, action: %d\ndescription: %s\nwant: %v\ngot no quorm",
					i+1, j+1, test.desc, action.wantVal)
			case action.wantOutput && quorum:
				if prm.Rnd != action.wantVal.Rnd {
					t.Errorf("test: %d, action: %d\ndescription: %s\nwant promise rnd: %v\ngot promise rnd: %v",
						i+1, j+1, test.desc, action.wantVal.Rnd, prm.Rnd)
				}
				if prm.Vrnd != action.wantVal.Vrnd {
					t.Errorf("test: %d, action: %d\ndescription: %s\nwant promise vrnd: %v\ngot promise vrnd: %v",
						i+1, j+1, test.desc, action.wantVal.Vrnd, prm.Vrnd)
				}
				if *prm.Vval != *action.wantVal.Vval {
					t.Errorf("test: %d, action: %d\ndescription: %s\nwant promise vval: %v\ngot promise vval: %v",
						i+1, j+1, test.desc, action.wantVal.Vval, prm.Vval)
				}
			}
		}
	}

}

func TestAcceptQF(t *testing.T) {
	// TODO This test is reusing the table driven tests from Tormod's original code;
	// TODO We should redesign the tests to better match the QSpec approach.
	qspec := NewPaxosQSpec(2)
	for i, test := range acceptQFTestcases {
		var replies []*LearnMsg
		for j, action := range test.actions {
			learn := action.learn
			replies = append(replies, &learn)
			lrn, quorum := qspec.AcceptQF(replies)
			gotVal := lrn.GetVal()
			switch {
			case !action.wantOutput && quorum:
				t.Errorf("test: %d, action: %d\ndescription: %s\nwant no quorum\ngot: %v",
					i+1, j+1, test.desc, gotVal)
			case action.wantOutput && !quorum:
				t.Errorf("test: %d, action: %d\ndescription: %s\nwant: %v\ngot no quorm",
					i+1, j+1, test.desc, action.wantVal)
			case action.wantOutput && quorum:
				if gotVal != action.wantVal {
					t.Errorf("test: %d, action: %d\ndescription: %s\nwant: %v\ngot: %v",
						i+1, j+1, test.desc, action.wantVal, gotVal)
				}
			}
		}
	}
}

// the client values for the tests
var (
	lamport = &Value{"Lamport"}
	leslie  = &Value{"Leslie"}
)

// prepareQFTests includes test cases for testing PrepareQF
type prepareQFTests struct {
	desc    string
	actions []prmAction
}

type prmAction struct {
	promise    PromiseMsg
	wantOutput bool
	wantVal    *PromiseMsg
}

// test cases for testing PrepareQF
var prepareQFTestcases = []prepareQFTests{
	{
		"no quorum -> no output",
		[]prmAction{
			{
				PromiseMsg{
					Rnd:  1,
					Vrnd: NoRound,
					Vval: ZeroValue,
				},
				false,
				nil,
			},
		},
	},
	{
		"valid quorum and no value reported",
		[]prmAction{
			{
				PromiseMsg{
					Rnd:  2,
					Vrnd: NoRound,
					Vval: ZeroValue,
				},
				false,
				nil,
			},
			{
				PromiseMsg{
					Rnd:  2,
					Vrnd: NoRound,
					Vval: ZeroValue,
				},
				true,
				&PromiseMsg{
					Rnd:  2,
					Vrnd: NoRound,
					Vval: ZeroValue,
				},
			},
		},
	},
	{
		"valid quorum and a value reported",
		[]prmAction{
			{
				PromiseMsg{
					Rnd:  2,
					Vrnd: NoRound,
					Vval: ZeroValue,
				},
				false,
				nil,
			},
			{
				PromiseMsg{
					Rnd:  2,
					Vrnd: 1,
					Vval: lamport,
				},
				true,
				&PromiseMsg{
					Rnd:  2,
					Vrnd: 1,
					Vval: lamport,
				},
			},
		},
	},
	{
		"three promises, different rounds -> ignore all promises",
		[]prmAction{
			{
				PromiseMsg{
					Rnd:  1,
					Vrnd: NoRound,
					Vval: ZeroValue,
				},
				false,
				nil,
			},
			{
				PromiseMsg{
					Rnd:  6,
					Vrnd: NoRound,
					Vval: ZeroValue,
				},
				false,
				nil,
			},
			{
				PromiseMsg{
					Rnd:  4,
					Vrnd: 1,
					Vval: lamport,
				},
				false,
				nil,
			},
		},
	},
	{
		"valid quorum and two different values reported -> propose correct value (highest vrnd) in accept",
		[]prmAction{
			{
				PromiseMsg{
					Rnd:  2,
					Vrnd: 1,
					Vval: lamport,
				},
				false,
				nil,
			},
			{
				PromiseMsg{
					Rnd:  2,
					Vrnd: 0,
					Vval: leslie,
				},
				true,
				&PromiseMsg{
					Rnd:  2,
					Vrnd: 1,
					Vval: lamport,
				},
			},
		},
	},
	{
		"valid quorum and two different values reported -> propose correct value (highest vrnd) in accept",
		[]prmAction{
			{
				PromiseMsg{
					Rnd:  2,
					Vrnd: 0,
					Vval: lamport,
				},
				false,
				nil,
			},
			{
				PromiseMsg{
					Rnd:  2,
					Vrnd: 1,
					Vval: leslie,
				},
				true,
				&PromiseMsg{
					Rnd:  2,
					Vrnd: 1,
					Vval: leslie,
				},
			},
		},
	},
}

// acceptQFTests includes test cases for testing AcceptQF
type acceptQFTests struct {
	desc    string
	actions []lrnAction
}

type lrnAction struct {
	learn      LearnMsg
	wantOutput bool
	wantVal    *Value
}

// test cases for testing AcceptQF
var acceptQFTestcases = []acceptQFTests{
	{

		"single learn, 3 nodes, no quorum -> no output",
		[]lrnAction{
			{
				LearnMsg{
					Rnd: 1,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
		},
	},
	{

		"two learns, 3 nodes, same round and value, unique senders = quorum -> report output and value",
		[]lrnAction{
			{
				LearnMsg{
					Rnd: 1,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
			{
				LearnMsg{
					Rnd: 1,
					Val: lamport,
				},
				true,
				lamport,
			},
		},
	},
	{
		"two learns, 3 nodes, different rounds, unique senders = no quorum -> no output",
		[]lrnAction{
			{
				LearnMsg{
					Rnd: 1,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
			{
				LearnMsg{
					Rnd: 2,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
		},
	},
	{
		"two learns, 3 nodes, second learn should be ignored due to lower round -> no output",
		[]lrnAction{
			{
				LearnMsg{
					Rnd: 2,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
			{
				LearnMsg{
					Rnd: 1,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
		},
	},
	{
		"3 nodes, single learn with rnd 2, then two learns with rnd 4 (quorum) -> report output and value of quorum",
		[]lrnAction{
			{
				LearnMsg{
					Rnd: 2,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
			{
				LearnMsg{
					Rnd: 4,
					Val: leslie,
				},
				false,
				ZeroValue,
			},
			{
				LearnMsg{
					Rnd: 4,
					Val: leslie,
				},
				true,
				leslie,
			},
		},
	},
	{
		"single learn, 3 nodes, no quorum -> no value",
		[]lrnAction{
			{
				LearnMsg{
					Rnd: 1,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
		},
	},
	{
		"two learns, 3 nodes, same round and value, unique senders = quorum -> report output and value",
		[]lrnAction{
			{
				LearnMsg{
					Rnd: 1,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
			{
				LearnMsg{
					Rnd: 1,
					Val: lamport,
				},
				true,
				lamport,
			},
		},
	},
	{
		"two learns, 3 nodes, different rounds, unique senders = no quorum -> no output",
		[]lrnAction{
			{
				LearnMsg{
					Rnd: 1,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
			{
				LearnMsg{
					Rnd: 2,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
		},
	},
	{
		"two learns, 3 nodes, second learn should be ignored due to lower round -> no output",
		[]lrnAction{
			{
				LearnMsg{
					Rnd: 2,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
			{
				LearnMsg{
					Rnd: 1,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
		},
	},
	{
		"(sanity check) two learns, 3 nodes, second learn should be ignored due to different value in same round -> no output",
		[]lrnAction{
			{
				LearnMsg{
					Rnd: 1,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
			{
				LearnMsg{
					Rnd: 1,
					Val: leslie,
				},
				false,
				ZeroValue,
			},
		},
	},
	{
		"3 nodes, single learn with rnd 2, then two learns with rnd 4 (quorum) -> report output and value of quorum",
		[]lrnAction{
			{
				LearnMsg{
					Rnd: 2,
					Val: lamport,
				},
				false,
				ZeroValue,
			},
			{
				LearnMsg{
					Rnd: 4,
					Val: leslie,
				},
				false,
				ZeroValue,
			},
			{
				LearnMsg{
					Rnd: 4,
					Val: leslie,
				},
				true,
				leslie,
			},
		},
	},
}
