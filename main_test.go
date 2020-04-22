package singlepaxos

import (
	"flag"
	r "github.com/selabhvl/singlepaxos/reader"
	"os"
	"testing"
)

var (
	prepareQFTest    r.PrepareQFTest
	acceptQFTest     r.AcceptQFTest
	acceptorTest     r.AcceptorTest
	leaderChangeTest r.LeaderChangeTest
	timeoutTest      r.TimeoutTest
	systemTest       r.SystemTest
)

func TestMain(m *testing.M) {
	// Flag definitions.
	// ******************************************************
	// directories of xml files for unit tests:
	// dir for prepareQF test cases
	var prepareQFTCsDir = flag.String(
		"prepareQFTCsDir",
		"./xml/examples/prepareqftest.xml",
		"path to unit test file for PrepareQF tests",
	)

	// dir for prepareQF test cases
	var acceptQFTCsDir = flag.String(
		"acceptQFTCsDir",
		"./xml/examples/acceptqftest.xml",
		"path to unit test file for AcceptQF tests",
	)

	// dir for acceptor test cases
	var acceptorTCsDir = flag.String(
		"acceptorTCsDir",
		"./xml/examples/acceptortest.xml",
		"path to unit test file for acceptor tests",
	)

	// dir for leaderChange test cases
	var leaderChangeTCsDir = flag.String(
		"leaderChangeTCsDir",
		"./xml/examples/leaderchangetest.xml",
		"path to unit test file for leader change tests",
	)

	// dir for leaderChange test cases
	var timeoutTCsDir = flag.String(
		"timeoutTCsDir",
		"./xml/examples/timeouttest.xml",
		"path to unit test file for timeout tests of fd",
	)

	// ******************************************************
	// directories of xml files for system test
	var systemTCsDir = flag.String(
		"systemTCsDir",
		"./xml/examples/systemtestexample.xml",
		"path to xml file for system test",
	)

	// Parse and validate flags.
	flag.Parse()

	// ******************************************************
	// Load test cases from XML files for unit tests:
	// Load the PrepareQF test cases from XML file
	r.ParseXMLTestCase(*prepareQFTCsDir, &prepareQFTest)
	// Load the AcceptQF test cases from XML file
	r.ParseXMLTestCase(*acceptQFTCsDir, &acceptQFTest)
	// Load the AcceptQF test cases from XML file
	r.ParseXMLTestCase(*acceptorTCsDir, &acceptorTest)
	// Load the LeaderChange test cases from XML file
	r.ParseXMLTestCase(*leaderChangeTCsDir, &leaderChangeTest)
	// Load the Timeout test cases from XML file
	r.ParseXMLTestCase(*timeoutTCsDir, &timeoutTest)

	// ******************************************************
	// Load the system test cases from XML file
	r.ParseXMLTestCase(*systemTCsDir, &systemTest)

	// Run tests/benchmarks.
	res := m.Run()
	os.Exit(res)
}
