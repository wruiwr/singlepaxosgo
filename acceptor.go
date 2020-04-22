package singlepaxos

var ZeroValue = &Value{}

// Acceptor represents an acceptor
type Acceptor struct {
	rnd  uint32
	vrnd uint32
	vval *Value
}

// NewAcceptor returns a new single-decree Paxos acceptor.
func NewAcceptor() *Acceptor {
	return &Acceptor{vval: ZeroValue}
}

// handlePrepare processes prepare message.
func (a *Acceptor) handlePrepare(prepare *PrepareMsg) (*PromiseMsg, bool) {
	if prepare.Rnd >= a.rnd {
		a.rnd = prepare.Rnd
		// send a promise massage back
		return &PromiseMsg{a.rnd, a.vrnd, a.vval}, true
	}

	// if need to ignore, just send back a NACK, which is a promise message with Ignore as a round number
	// this message will be checked from the propoer
	return &PromiseMsg{Ignore, Ignore, ZeroValue}, false
}

// handleAccept processes accept message.
func (a *Acceptor) handleAccept(accept *AcceptMsg) (*LearnMsg, bool) {
	if accept.Rnd >= a.rnd {
		// update the acceptor state for this replica and vote by issuing a LearnMsg.
		a.rnd = accept.Rnd
		a.vrnd = accept.Rnd
		a.vval = accept.Val
		// voting for vrnd, vval
		return &LearnMsg{a.vrnd, a.vval}, true
	}

	// if need to ignore, just send back a NACK, which is a learn message with Ignore as a round number
	// this message will be checked from the proposer
	return &LearnMsg{Ignore, ZeroValue}, false
}
