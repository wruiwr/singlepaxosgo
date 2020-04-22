package singlepaxos

// PaxosQSpec is a quorum specification object for Paxos.
// It only holds the quorum size.
type PaxosQSpec struct {
	qSize int
}

// NewPaxosQSpec returns a quorum specification object for Paxos
// for the given quorum size.
func NewPaxosQSpec(quorumSize int) QuorumSpec {
	return &PaxosQSpec{
		qSize: quorumSize,
	}
}

// PrepareQF is the quorum function for the Prepare
// quorum call method.
// This is where the Proposer handle PromiseMsgs returned by the Acceptors
// to determine if a value from a previous round needs to be considered.
func (qs PaxosQSpec) PrepareQF(prepare *PrepareMsg, replies []*PromiseMsg) (*PromiseMsg, bool) {
	if len(replies) < qs.qSize {
		return nil, false
	}
	reply := &PromiseMsg{Rnd: prepare.GetRnd()}
	for _, r := range replies {
		if r.GetRnd() != reply.GetRnd() {
			// fail the sanity check
			return nil, false
		}
		if r.GetVrnd() >= reply.GetVrnd() {
			reply.Vrnd = r.GetVrnd()
			reply.Vval = r.GetVval()
		}
	}
	return reply, true
}

// AcceptQF is the quorum function for the Accept
// quorum call method.
// This is where the Proposer handle LearnMsgs to determine if a
// value has been decided by the Acceptors.
// The quorum function returns true if a value has been decided,
// and the corresponding LearnMsg holds the round number and value
// that was decided. If false is returned, no value was decided.
func (qs PaxosQSpec) AcceptQF(replies []*LearnMsg) (*LearnMsg, bool) {
	if len(replies) < qs.qSize {
		return nil, false
	}
	validLearns := len(replies)
	var highest *LearnMsg
	// find a learn with highest round
	for _, reply := range replies {
		if highest != nil && reply.GetRnd() <= highest.GetRnd() {
			continue
		}
		highest = reply
	}
	// discount invalid learns
	for _, reply := range replies {
		if reply.GetRnd() < highest.GetRnd() || !reply.Equal(highest) {
			validLearns--
		}
	}
	// check if we have a quorum of valid learns
	if validLearns >= qs.qSize {
		return highest, true
	}
	return nil, false
}

// CommitQF is the quorum function for the Commit
// quorum call method.
// This function just waits for a quorum of empty replies,
// indicating that at least a quorum of Learners have committed
// the value decided by the Acceptors.
func (qs PaxosQSpec) CommitQF(replies []*Empty) (*Empty, bool) {
	if len(replies) < qs.qSize {
		return nil, false
	}
	return replies[0], true
}

func (qs PaxosQSpec) ClientHandleQF(replies []*Response) (*Response, bool) {
	if len(replies) < qs.qSize {
		return nil, false
	}

	return replies[0], true
}
