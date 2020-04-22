package reader

import (
	"encoding/xml"
	"io/ioutil"
)

// PrepareQFTest is an exported type for PrepareQF tests
type PrepareQFTest struct {
	XMLName xml.Name `xml:"Test"`
	Name    string   `xml:"TestName,attr"`
	//SystemSize int                   `xml:"SystemSize,attr"`
	//QuorumSize int                   `xml:"QuorumSize,attr"`
	Configuration *Configuration        `xml:"Configuration"`
	TestCases     []*PrepareQFTestCases `xml:"TestCase"`
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
	ExpectQuorum bool        `xml:"Quorum"`
	ExpectResult *PromiseMsg `xml:"PromiseMsg"`
}

type PromiseMsg struct {
	Rnd  uint32 `xml:"Rnd"`
	Vrnd uint32 `xml:"Vrnd"`
	Vval string `xml:"Vval"`
}

// AcceptQFTest is an exported type for AcceptQF tests
type AcceptQFTest struct {
	XMLName xml.Name `xml:"Test"`
	Name    string   `xml:"TestName,attr"`
	//SystemSize int                  `xml:"SystemSize,attr"`
	//QuorumSize int                  `xml:"QuorumSize,attr"`
	Configuration *Configuration       `xml:"Configuration"`
	TestCases     []*AcceptQFTestCases `xml:"TestCase"`
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
	ExpectQuorum bool      `xml:"Quorum"`
	ExpectResult *LearnMsg `xml:"LearnMsg"`
}

type LearnMsg struct {
	Rnd uint32 `xml:"Rnd"`
	Val string `xml:"Val"`
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
	ClientRequests *[]string `xml:"ClientPropose"`
	P1Failures     *[]string `xml:"P1Failure"`
	P2Failures     *[]string `xml:"P2Failure"`
}

type SystemTestOracles struct {
	XMLName                xml.Name  `xml:"TestOracles"`
	ExpectedLegalResponses *[]string `xml:"Response"`
	ExpectLeaders          *[]string `xml:"Leader"`
}

func ParseXMLTestCase(file string, xmlTestCaseType interface{}) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return xml.Unmarshal(b, &xmlTestCaseType)
}
