package statutils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/senseyeio/roger"
)

// ComputeToleranceIntervalNormalDist computes the tolerance interval and returns the the lower and upper bound
// todo add the p and alpha as parameters
func ComputeToleranceIntervalNormalDist(data []float64) (float64, float64, error) {
	valuesStr := make([]string, len(data)+1)
	valuesStr[0] = "/Users/vahid/Desktop/projects/swarm-manager/internal/statutils/tolerance_interval.py" //TODO
	for i, v := range data {
		valuesStr[i+1] = strconv.FormatFloat(v, 'f', 6, 64)
	}
	cmd := exec.Command("python", valuesStr...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 0, err
	}
	outStr := strings.Trim(string(out), "\n")
	outStr = strings.Trim(outStr, " ")
	tempStrs := strings.Split(outStr, ",")
	count, err := strconv.Atoi(tempStrs[0])
	if err != nil {
		return 0, 0, err
	}
	if count != len(data) {
		return 0, 0, fmt.Errorf("count returned by python code is not equal to len(data)")
	}
	lower, err := strconv.ParseFloat(tempStrs[1], 64)
	if err != nil {
		return 0, 0, err
	}
	upper, err := strconv.ParseFloat(tempStrs[2], 64)
	if err != nil {
		return 0, 0, err
	}
	return lower, upper, nil
}

// ComputeToleranceIntervalNonParametric returns both lower and upper bound
func ComputeToleranceIntervalNonParametric(data []float64, confidence, portionOfPopulation float64) (float64, float64, error) {
	rClient, err := roger.NewRClient("127.0.0.1", 6311)
	if err != nil {
		return 0, 0, fmt.Errorf("error connecting to R server: %w", err)
	}
	numbers := ""
	for i, v := range data {
		numbers += strconv.FormatFloat(v, 'f', 8, 64)
		if i != len(data)-1 {
			numbers += ", "
		}
	}
	alpha := strconv.FormatFloat(1-float64(confidence), 'f', 4, 64)
	p := strconv.FormatFloat(portionOfPopulation, 'f', 4, 64)
	// log.Println("Statutils: alpha:", alpha, "P:", p)
	command := "library(\"tolerance\"); nptol.int(x = c(" + numbers + "), alpha = " + alpha + ", P = " + p + ", side = 1,  method = \"WALD\", upper = NULL, lower = NULL)"

	value, err := rClient.Eval(command)
	if err != nil {
		return 0, 0, fmt.Errorf("error executing command in R server: %w", err)
	}
	upper := value.(map[string]interface{})["1-sided.upper"].(float64)
	lower := value.(map[string]interface{})["1-sided.lower"].(float64)
	return lower, upper, nil
}
