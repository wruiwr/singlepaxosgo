package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// PrepareQFTest is an exported type for PrepareQF tests
type PrepareQFTest struct {
	XMLName    xml.Name              `xml:"Test"`
	Name       string                `xml:"TestName,attr"`
	SystemSize int                   `xml:"SystemSize,attr"`
	QuorumSize int                   `xml:"QuorumSize,attr"`
	TestCases  []*PrepareQFTestCases `xml:"TestCase"`
}

type PrepareQFTestCases struct {
	XMLName     xml.Name              `xml:"TestCase"`
	CaseID      string                `xml:"CaseID,attr"`
	TestValues  *PrepareQFTestValues  `xml:"TestValues"`
	TestOracles *PrepareQFTestOracles `xml:"TestOracles"`
}

type PrepareQFTestValues struct {
	XMLName     xml.Name      `xml:"TestValues"`
	PromiseMsgs []*PromiseMsg `xml:"PromiseMsg"`
}

type PrepareQFTestOracles struct {
	XMLName      xml.Name    `xml:"TestOracles"`
	ExpectQuorum bool        `xml:"ExpectQuorum"`
	ExpectResult *PromiseMsg `xml:"ExpectPromiseMsg"`
}

type PromiseMsg struct {
	Rnd  uint32 `xml:"Rnd"`
	Vrnd uint32 `xml:"Vrnd"`
	Vval string `xml:"Vval"`
}

// AcceptQFTest is an exported type for AcceptQF tests
type AcceptQFTest struct {
	XMLName    xml.Name             `xml:"Test"`
	Name       string               `xml:"TestName,attr"`
	SystemSize int                  `xml:"SystemSize,attr"`
	QuorumSize int                  `xml:"QuorumSize,attr"`
	TestCases  []*AcceptQFTestCases `xml:"TestCase"`
}

type AcceptQFTestCases struct {
	XMLName     xml.Name             `xml:"TestCase"`
	CaseID      string               `xml:"CaseID,attr"`
	TestValues  *AcceptQFTestValues  `xml:"TestValues"`
	TestOracles *AcceptQFTestOracles `xml:"TestOracles"`
}

type AcceptQFTestValues struct {
	XMLName   xml.Name    `xml:"TestValues"`
	LearnMsgs []*LearnMsg `xml:"LearnMsg"`
}

type AcceptQFTestOracles struct {
	XMLName      xml.Name  `xml:"TestOracles"`
	ExpectQuorum bool      `xml:"ExpectQuorum"`
	ExpectResult *LearnMsg `xml:"ExpectLearnMsg"`
}

type LearnMsg struct {
	Rnd uint32 `xml:"Rnd"`
	Val string `xml:"Val"`
}

// prepareQFTestcases contain test cases for testing prepareQF.
var prepareQFTestcases = PrepareQFTest{
	Name:       "PrepareQFTest",
	SystemSize: 3,
	QuorumSize: 2,
	TestCases: []*PrepareQFTestCases{
		{
			CaseID: "1",
			TestValues: &PrepareQFTestValues{

				PromiseMsgs: []*PromiseMsg{
					{
						Rnd:  8080,
						Vrnd: 0,
						Vval: "",
					},
				},
			},
			TestOracles: &PrepareQFTestOracles{
				ExpectQuorum: false,
				ExpectResult: nil,
			},
		},
		{
			CaseID: "2",
			TestValues: &PrepareQFTestValues{
				PromiseMsgs: []*PromiseMsg{
					{
						Rnd:  8080,
						Vrnd: 0,
						Vval: "",
					},
					{
						Rnd:  8080,
						Vrnd: 0,
						Vval: "",
					},
				},
			},
			TestOracles: &PrepareQFTestOracles{
				ExpectQuorum: true,
				ExpectResult: &PromiseMsg{
					Rnd:  8080,
					Vrnd: 0,
					Vval: "",
				},
			},
		},
		{
			CaseID: "3",
			TestValues: &PrepareQFTestValues{

				PromiseMsgs: []*PromiseMsg{
					{
						Rnd:  8081,
						Vrnd: 0,
						Vval: "",
					},
					{
						Rnd:  8081,
						Vrnd: 8080,
						Vval: "M1",
					},
				},
			},
			TestOracles: &PrepareQFTestOracles{
				ExpectQuorum: true,
				ExpectResult: &PromiseMsg{
					Rnd:  8081,
					Vrnd: 8080,
					Vval: "M1",
				},
			},
		},
	},
}

// acceptQFTestcases contain test cases for testing AcceptQF.
var acceptQFTestcases = AcceptQFTest{
	Name:       "PrepareQFTest",
	SystemSize: 3,
	QuorumSize: 2,
	TestCases: []*AcceptQFTestCases{
		{
			CaseID: "1",
			TestValues: &AcceptQFTestValues{

				LearnMsgs: []*LearnMsg{
					{
						Rnd: 8080,
						Val: "M1",
					},
				},
			},
			TestOracles: &AcceptQFTestOracles{
				ExpectQuorum: false,
				ExpectResult: nil,
			},
		},
		{
			CaseID: "2",
			TestValues: &AcceptQFTestValues{
				LearnMsgs: []*LearnMsg{
					{
						Rnd: 8080,
						Val: "M1",
					},
					{
						Rnd: 8080,
						Val: "M1",
					},
				},
			},
			TestOracles: &AcceptQFTestOracles{
				ExpectQuorum: true,
				ExpectResult: &LearnMsg{
					Rnd: 8080,
					Val: "M1",
				},
			},
		},
		{
			CaseID: "3",
			TestValues: &AcceptQFTestValues{
				LearnMsgs: []*LearnMsg{
					{
						Rnd: 8080,
						Val: "",
					},
					{
						Rnd: 8081,
						Val: "M1",
					},
					{
						Rnd: 8081,
						Val: "M1",
					},
				},
			},
			TestOracles: &AcceptQFTestOracles{
				ExpectQuorum: true,
				ExpectResult: &LearnMsg{
					Rnd: 8081,
					Val: "M1",
				},
			},
		},
	},
}

// LeaderChangeTest is an exported type for leader change tests
type LeaderChangeTest struct {
	XMLName    xml.Name                 `xml:"Test"`
	Name       string                   `xml:"TestName,attr"`
	SystemSize int                      `xml:"SystemSize,attr"`
	QuorumSize int                      `xml:"QuorumSize,attr"`
	TestCases  []*LeaderChangeTestcases `xml:"TestCase"`
}

type LeaderChangeTestcases struct {
	XMLName     xml.Name                 `xml:"TestCase"`
	CaseID      string                   `xml:"CaseID,attr"`
	TestValues  *LeaderChangeTestValues  `xml:"TestValue"`
	TestOracles *LeaderChangeTestOracles `xml:"TestOracles"`
}

type LeaderChangeTestValues struct {
	XMLName   xml.Name `xml:"TestValue"`
	ServerIDs *[]int   `xml:"ServerIDs>ID"`
}

type LeaderChangeTestOracles struct {
	XMLName       xml.Name `xml:"TestOracles"`
	ExpectLeaders *[]int   `xml:"ExpectLeaders>Leader"`
}

// leaderChangeTestcases contain test cases for testing leader changes.
var leaderChangeTestcases = LeaderChangeTest{
	Name:       "LeaderChangeTest",
	SystemSize: 3,
	QuorumSize: 2,
	TestCases: []*LeaderChangeTestcases{
		{
			CaseID: "1",
			TestValues: &LeaderChangeTestValues{
				ServerIDs: &[]int{8080, 8081, 8082},
			},
			TestOracles: &LeaderChangeTestOracles{
				ExpectLeaders: &[]int{8080, 8081},
			},
		},
	},
}

// TimeoutTest is an exported type for timeout tests
type TimeoutTest struct {
	XMLName    xml.Name            `xml:"Test"`
	Name       string              `xml:"TestName,attr"`
	SystemSize int                 `xml:"SystemSize,attr"`
	QuorumSize int                 `xml:"QuorumSize,attr"`
	TestCases  []*TimeoutTestcases `xml:"TestCase"`
}

type TimeoutTestcases struct {
	XMLName     xml.Name            `xml:"TestCase"`
	CaseID      string              `xml:"CaseID,attr"`
	TestValues  *TimeoutTestValues  `xml:"TestValues"`
	TestOracles *TimeoutTestOracles `xml:"TestOracles"`
}

type TimeoutTestValues struct {
	XMLName   xml.Name `xml:"TestValues"`
	ServerIDs *[]int   `xml:"ServerIDs>ID"`
	Alives    *[]bool  `xml:"Alives>Alive"`
	Suspects  *[]bool  `xml:"Suspects>Suspected"`
}

type TimeoutTestOracles struct {
	XMLName        xml.Name `xml:"TestOracles"`
	ExpectSuspects *[]bool  `xml:"ExpectSuspects>ExpectSuspect"`
}

// timeoutTestcases contain test cases for testing timeout of fd.
var timeoutTestcases = TimeoutTest{
	Name:       "TimeoutTest",
	SystemSize: 3,
	QuorumSize: 2,
	TestCases: []*TimeoutTestcases{
		{
			CaseID: "1",
			TestValues: &TimeoutTestValues{
				ServerIDs: &[]int{8080, 8081, 8082},
				Alives:    &[]bool{true, true, true},
				Suspects:  &[]bool{},
			},
			TestOracles: &TimeoutTestOracles{
				ExpectSuspects: &[]bool{},
			},
		},
		{
			CaseID: "2",
			TestValues: &TimeoutTestValues{
				ServerIDs: &[]int{8080, 8081, 8082},
				Alives:    &[]bool{true, false, true},
				Suspects:  &[]bool{false, true, false},
			},
			TestOracles: &TimeoutTestOracles{
				ExpectSuspects: &[]bool{false, true, false},
			},
		},
		{
			CaseID: "3",
			TestValues: &TimeoutTestValues{
				ServerIDs: &[]int{8080, 8081, 8082},
				Alives:    &[]bool{false, false, false},
				Suspects:  &[]bool{false, false, false},
			},
			TestOracles: &TimeoutTestOracles{
				ExpectSuspects: &[]bool{true, true, true},
			},
		},
		{
			CaseID: "4",
			TestValues: &TimeoutTestValues{
				ServerIDs: &[]int{8080, 8081, 8082},
				Alives:    &[]bool{false, false, true},
				Suspects:  &[]bool{false, false, false},
			},
			TestOracles: &TimeoutTestOracles{
				ExpectSuspects: &[]bool{true, true, false},
			},
		},
		{
			CaseID: "5",
			TestValues: &TimeoutTestValues{
				ServerIDs: &[]int{8080, 8081, 8082},
				Alives:    &[]bool{false, false, false},
				Suspects:  &[]bool{true, true, true},
			},
			TestOracles: &TimeoutTestOracles{
				ExpectSuspects: &[]bool{true, true, true},
			},
		},
	},
}

// AcceptorTest is an exported type for testing Acceptor
type AcceptorTest struct {
	XMLName    xml.Name             `xml:"Test"`
	Name       string               `xml:"TestName,attr"`
	SystemSize int                  `xml:"SystemSize,attr"`
	QuorumSize int                  `xml:"QuorumSize,attr"`
	TestCases  []*AcceptorTestcases `xml:"TestCase"`
}

type AcceptorTestcases struct {
	XMLName    xml.Name              `xml:"TestCase"`
	CaseID     string                `xml:"CaseID,attr"`
	TestValues []*AcceptorTestValues `xml:"TestValue"`
}

type AcceptorTestValues struct {
	MsgType                 string      `xml:"MsgType"`
	PrepareMsg              *PrepareMsg `xml:"PrepareMsg"`
	AcceptMsg               *AcceptMsg  `xml:"AcceptMsg"`
	ExpectOutput            bool        `xml:"ExpectOutput"`
	ExpectPromiseMsgResults *PromiseMsg `xml:"ExpectPromise"`
	ExpectLearnMsgResults   *LearnMsg   `xml:"ExpectLearn"`
}

type PrepareMsg struct {
	Rnd uint32 `xml:"Rnd"`
}

type AcceptMsg struct {
	Rnd uint32 `xml:"Rnd"`
	Val string `xml:"Val"`
}

// acceptorTestcases contain test cases for testing acceptor
var acceptorTestcases = AcceptorTest{
	Name:       "TimeoutTest",
	SystemSize: 3,
	QuorumSize: 2,
	TestCases: []*AcceptorTestcases{
		{
			CaseID: "1",
			TestValues: []*AcceptorTestValues{
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8080,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8080,
						Vrnd: 0,
						Vval: "",
					},
				},
			},
		},
		{
			CaseID: "2",
			TestValues: []*AcceptorTestValues{
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8080,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8080,
						Vrnd: 0,
						Vval: "",
					},
				},
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8081,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8081,
						Vrnd: 0,
						Vval: "",
					},
				},
			},
		},
		{
			CaseID: "3",
			TestValues: []*AcceptorTestValues{
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8080,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8080,
						Vrnd: 0,
						Vval: "",
					},
				},
				{
					MsgType: "accept",
					AcceptMsg: &AcceptMsg{
						Rnd: 8080,
						Val: "M1",
					},
					ExpectOutput: true,
					ExpectLearnMsgResults: &LearnMsg{
						Rnd: 8080,
						Val: "M1",
					},
				},
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8081,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8081,
						Vrnd: 8080,
						Vval: "M1",
					},
				},
			},
		},
		{
			CaseID: "4",
			TestValues: []*AcceptorTestValues{
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8081,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8081,
						Vrnd: 0,
						Vval: "",
					},
				},
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8080,
					},
					ExpectOutput: false,
				},
			},
		},
		{
			CaseID: "5",
			TestValues: []*AcceptorTestValues{
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8081,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8081,
						Vrnd: 0,
						Vval: "",
					},
				},
				{
					MsgType: "accept",
					AcceptMsg: &AcceptMsg{
						Rnd: 8080,
						Val: "M1",
					},
					ExpectOutput: false,
				},
			},
		},
		{
			CaseID: "6",
			TestValues: []*AcceptorTestValues{
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8080,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8080,
						Vrnd: 0,
						Vval: "",
					},
				},
			},
		},
		{
			CaseID: "7",
			TestValues: []*AcceptorTestValues{
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8080,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8080,
						Vrnd: 0,
						Vval: "",
					},
				},
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8081,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8081,
						Vrnd: 0,
						Vval: "",
					},
				},
			},
		},
		{
			CaseID: "8",
			TestValues: []*AcceptorTestValues{
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8080,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8080,
						Vrnd: 0,
						Vval: "",
					},
				},
				{
					MsgType: "accept",
					AcceptMsg: &AcceptMsg{
						Rnd: 8080,
						Val: "M1",
					},
					ExpectOutput: true,
					ExpectLearnMsgResults: &LearnMsg{
						Rnd: 8080,
						Val: "M1",
					},
				},
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8081,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8081,
						Vrnd: 8080,
						Vval: "M1",
					},
				},
			},
		},
		{
			CaseID: "9",
			TestValues: []*AcceptorTestValues{
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8081,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8081,
						Vrnd: 0,
						Vval: "",
					},
				},
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8080,
					},
					ExpectOutput: false,
				},
			},
		},
		{
			CaseID: "10",
			TestValues: []*AcceptorTestValues{
				{
					MsgType: "prepare",
					PrepareMsg: &PrepareMsg{
						Rnd: 8081,
					},
					ExpectOutput: true,
					ExpectPromiseMsgResults: &PromiseMsg{
						Rnd:  8081,
						Vrnd: 0,
						Vval: "",
					},
				},
				{
					MsgType: "accept",
					AcceptMsg: &AcceptMsg{
						Rnd: 8080,
						Val: "M1",
					},
					ExpectOutput: false,
				},
			},
		},
	},
}

// SystemTest is an exported type for system tests
type SystemTest struct {
	XMLName       xml.Name           `xml:"Test"`
	Name          string             `xml:"TestName,attr"`
	Configuration *Configuration     `xml:"Configuration"`
	TestCases     []*SystemTestcases `xml:"TestCase"`
}

type Configuration struct {
	XMLName         xml.Name `xml:"Configuration"`
	SystemSize      int      `xml:"SystemSize"`
	QuorumSize      int      `xml:"QuorumSize"`
	ServerIDs       *[]int   `xml:"ServerIDs>ID"`
	FailurePhaseOne int      `xml:"FailurePhaseOne"`
	FailurePhaseTwo int      `xml:"FailurePhaseTwo"`
}

type SystemTestcases struct {
	XMLName     xml.Name           `xml:"TestCase"`
	CaseID      string             `xml:"CaseID,attr"`
	TestValues  *SystemTestValues  `xml:"TestValues"`
	TestOracles *SystemTestOracles `xml:"TestOracles"`
}

type SystemTestValues struct {
	XMLName        xml.Name  `xml:"TestValues"`
	ClientRequests *[]string `xml:"ClientRequests>Requests"`
}

type SystemTestOracles struct {
	XMLName                xml.Name  `xml:"TestOracles"`
	ExpectedLegalResponses *[]string `xml:"Response"`
	ExpectLeaders          *[]string `xml:"Leader"`
}

var systemTestcaseOneZero = SystemTest{
	Name: "systemtest-ss-s3-1-0",
	Configuration: &Configuration{
		SystemSize:      3,
		QuorumSize:      2,
		ServerIDs:       &[]int{8080, 8081, 8082},
		FailurePhaseOne: 1,
		FailurePhaseTwo: 0,
	},

	TestCases: []*SystemTestcases{
		{
			CaseID: "1",
			TestValues: &SystemTestValues{
				ClientRequests: &[]string{"M1", "M2"},
			},
			TestOracles: &SystemTestOracles{
				ExpectedLegalResponses: &[]string{"M1", "M2"},
				ExpectLeaders:          &[]string{"8080"},
			},
		},
		{
			CaseID: "2",
			TestValues: &SystemTestValues{
				ClientRequests: &[]string{"M1", "M2"},
			},
			TestOracles: &SystemTestOracles{
				ExpectedLegalResponses: &[]string{"M1", "M2"},
				ExpectLeaders:          &[]string{"8080", "8081"},
			},
		},
	},
}

var systemTestcaseZeroOne = SystemTest{
	Name: "systemtest-ss-s3-0-1",
	Configuration: &Configuration{
		SystemSize:      3,
		QuorumSize:      2,
		ServerIDs:       &[]int{8080, 8081, 8082},
		FailurePhaseOne: 0,
		FailurePhaseTwo: 1,
	},

	TestCases: []*SystemTestcases{
		{
			CaseID: "1",
			TestValues: &SystemTestValues{
				ClientRequests: &[]string{"M1", "M2"},
			},
			TestOracles: &SystemTestOracles{
				ExpectedLegalResponses: &[]string{"M1", "M2"},
				ExpectLeaders:          &[]string{"8080"},
			},
		},
		{
			CaseID: "2",
			TestValues: &SystemTestValues{
				ClientRequests: &[]string{"M1", "M2"},
			},
			TestOracles: &SystemTestOracles{
				ExpectedLegalResponses: &[]string{"M1", "M2"},
				ExpectLeaders:          &[]string{"8080", "8081"},
			},
		},
	},
}

var systemTestcaseOneOne = SystemTest{
	Name: "systemtest-ss-s5-1-1",
	Configuration: &Configuration{
		SystemSize:      5,
		QuorumSize:      3,
		ServerIDs:       &[]int{8080, 8081, 8082, 8083, 8084},
		FailurePhaseOne: 1,
		FailurePhaseTwo: 1,
	},

	TestCases: []*SystemTestcases{
		{
			CaseID: "1",
			TestValues: &SystemTestValues{
				ClientRequests: &[]string{"M1", "M2"},
			},
			TestOracles: &SystemTestOracles{
				ExpectedLegalResponses: &[]string{"M1", "M2"},
				ExpectLeaders:          &[]string{"8080"},
			},
		},
		{
			CaseID: "2",
			TestValues: &SystemTestValues{
				ClientRequests: &[]string{"M1", "M2"},
			},
			TestOracles: &SystemTestOracles{
				ExpectedLegalResponses: &[]string{"M1", "M2"},
				ExpectLeaders:          &[]string{"8080", "8081"},
			},
		},
		{
			CaseID: "3",
			TestValues: &SystemTestValues{
				ClientRequests: &[]string{"M1", "M2"},
			},
			TestOracles: &SystemTestOracles{
				ExpectedLegalResponses: &[]string{"M1", "M2"},
				ExpectLeaders:          &[]string{"8080", "8081", "8082"},
			},
		},
	},
}

// xmlWriter can write test cases for ReadQF, WriteQF or System tests into xml files.
func xmlWriter(dir string, testcases interface{}) {
	output, err := xml.MarshalIndent(testcases, " ", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	testFilePath, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create(testFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	_, err = f.Write(output)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Wrote XML output to:\n", testFilePath)
}

func main() {
	xmlWriter("./xml/unittests/prepareqftest.xml", prepareQFTestcases)
	xmlWriter("./xml/unittests/acceptqftest.xml", acceptQFTestcases)
	xmlWriter("./xml/unittests/leaderchangetest.xml", leaderChangeTestcases)
	xmlWriter("./xml/unittests/timeouttest.xml", timeoutTestcases)
	xmlWriter("./xml/unittests/acceptortest.xml", acceptorTestcases)
	xmlWriter("./xml/systemtests/systemtest-ss-1-0.xml", systemTestcaseOneZero)
	xmlWriter("./xml/systemtests/systemtest-ss-0-1.xml", systemTestcaseZeroOne)
	xmlWriter("./xml/systemtests/systemtest-ss-1-1.xml", systemTestcaseOneOne)
}
