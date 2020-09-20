package jaeger

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/pkg/errors"
)

// import ("github.com/PaesslerAG/gval")

func getValue(token string, t *trace) (float64, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 3 {
		fmt.Println(parts, token)
	}
	serviceName := parts[0]
	operationName := parts[1]
	indicator := parts[2]

	var matchingSpan *span
	for _, s := range t.Spans {
		if s.ServiceName == serviceName && s.OperationName == operationName {
			matchingSpan = s
			break
		}
	}
	if matchingSpan == nil {
		if operationName == "call" {
			return 0, nil
		}
		return 0, fmt.Errorf("matchingSpan is nil for %s %s", serviceName, operationName)
	}
	if indicator == "StartTime" {
		return matchingSpan.StartTime / 1000, nil
	} else if indicator == "EndTime" {
		return matchingSpan.EndTime / 1000, nil
	} else if indicator == "Duration" {
		return matchingSpan.Duration / 1000, nil
	} else {
		return 0, fmt.Errorf("no indicator for %s", indicator)
	}
}

func evaluateJaegerFormula(formula string, t *trace) (float64, error) {
	// serviceName.operationName.[StartTime|Duration|EndTime]

	for strings.Contains(formula, "  ") {
		formula = strings.ReplaceAll(formula, "  ", " ")
	}
	formulaWithNumbers := ""

	for _, token := range strings.Split(formula, " ") {
		if token == "(" || token == ")" || token == "+" || token == "-" {
			formulaWithNumbers += token
		} else {
			value, err := getValue(token, t)
			if err != nil {
				return 0, err
			}
			formulaWithNumbers += strconv.FormatFloat(value, 'f', 6, 64)
		}
	}

	res, err := gval.Evaluate(formulaWithNumbers, map[string]interface{}{})
	if err != nil {
		return 0, errors.Wrap(err, "couldn't evaluate expression: "+formulaWithNumbers)
	}
	resStr := fmt.Sprintf("%v", res)
	resValue, err := strconv.ParseFloat(resStr, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "Couldn't parse result as float: %s", resStr)
	}
	return resValue, nil
}
