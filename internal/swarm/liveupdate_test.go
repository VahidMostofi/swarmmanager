package swarm

import (
	"fmt"
	"strconv"
	"testing"
)

func TestJaegerParser(t *testing.T) {
	cpuCount := 0.47
	var replicaCount uint64 = 13
	cpuCountStr := strconv.FormatFloat(cpuCount, 'f', 6, 64)
	replicaCountStr := strconv.FormatUint(replicaCount, 10)
	fmt.Println(cpuCountStr + " " + replicaCountStr)
}
