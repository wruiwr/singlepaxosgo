package singlepaxos

import "testing"

var (
	nodes = []*Node{{id: 8080, addr: "localhost:8080"}, {id: 8081, addr: "localhost:8081"}, {id: 8082, addr: "localhost:8082"}}
)

// the tests below are for testing proposer.increaseCrnd()
var incCrndTests = []struct {
	id    uint32
	port  int
	crnds []uint32
}{
	{
		8080,
		8080,
		[]uint32{8080, 8083, 8086, 8089, 8092, 8095, 8098, 8101},
	},
	{
		8081,
		8081,
		[]uint32{8081, 8084, 8087, 8090, 8093, 8096, 8099, 8102},
	},
	{
		8082,
		8082,
		[]uint32{8082, 8085, 8088, 8091, 8094, 8097, 8100, 8103},
	},
}

func TestIncreaseCrnd(t *testing.T) {
	for i, test := range incCrndTests {
		config := &Configuration{id: test.id, nodes: nodes, n: len(nodes)}
		proposer := NewProposer(config, test.port)
		for j, wantCrnd := range test.crnds {
			proposer.increaseCrnd()
			if proposer.crnd != wantCrnd {
				t.Errorf("TestIncreaseCrnd %d, inc nr %d: proposer has current crnd %d, should have %d",
					i, j, proposer.crnd, wantCrnd)
			}
		}
	}
}
