package statutils

import (
	"os/exec"
	"strconv"
	"strings"
)

// ComputeToleranceInterval computes the tolerance interval and returns the the lower and upper bound
// todo add the p and alpha as parameters
func ComputeToleranceInterval(data []float64) (float64, float64, error) {
	valuesStr := make([]string, len(data)+1)
	valuesStr[0] = "/Users/vahid/Desktop/projects/swarm-manager/internal/statutils/tolerance_interval.py"
	for i, v := range data {
		valuesStr[i+1] = strconv.FormatFloat(v, 'f', 6, 64)
	}
	cmd := exec.Command("python", valuesStr...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// 	if (strings.Contains(string(out), "not found") || strings.Contains(string(out), "cannot be used with services.")) && strings.Contains(string(out), "network") && attempt <= 25 {
		// 		var waitTime int64 = 5
		// 		log.Printf("deploying stack, attempt %d failed. Wait %d seconds\n", attempt, waitTime)
		// 		time.Sleep(time.Duration(waitTime) * time.Second)
		// 		return s.DeployStackWithDockerCompose(dockerComposePath, attempt+1)
		// 	}
		// 	return fmt.Errorf("deploying stack with docker compose file failed with error: %w; %s", err, string(out))
	}
	outStr := strings.Trim(string(out), "\n")
	outStr = strings.Trim(outStr, " ")
	tempStrs := strings.Split(outStr, ",")
	lower, err := strconv.ParseFloat(tempStrs[0], 64)
	if err != nil {
		return 0, 0, err
	}
	upper, err := strconv.ParseFloat(tempStrs[1], 64)
	if err != nil {
		return 0, 0, err
	}
	return lower, upper, nil
}
