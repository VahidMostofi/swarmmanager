package jaeger

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/PaesslerAG/gval"
	"github.com/pkg/errors"
)

// import ("github.com/PaesslerAG/gval")

func evaluateJaegerFormula(formula string, spans map[string]*span) (float64, error) {
	// gval.Evaluable
	operationNames := make([]string, 0)
	for operationName := range spans {
		operationNames = append(operationNames, operationName)
	}

	sort.Slice(operationNames, func(i, j int) bool {
		return len(operationNames[i]) > len(operationNames[j])
	})
	for strings.Contains(formula, " ") {
		formula = strings.Replace(formula, " ", "", 1)
	}
	token := ""
	formulaWithNumbers := ""
	for i := 0; i < len(formula); i++ {
		if formula[i] == '+' || formula[i] == '-' || formula[i] == '(' || formula[i] == ')' {
			if len(token) > 0 {
				opName := strings.Split(token, ".")[0]
				field := strings.Split(token, ".")[1]
				var value float64
				if field == "StartTime" {
					value = spans[opName].StartTime
				} else if field == "EndTime" {
					value = spans[opName].EndTime
				} else if field == "Duration" {
					value = spans[opName].Duration
				} else {
					return 0, fmt.Errorf("the filed %s is undefined", field)
				}
				formulaWithNumbers += strconv.FormatFloat(value, 'f', 6, 64)
			}
			token = ""
			formulaWithNumbers += string(formula[i])
		} else {
			token += string(formula[i])
		}
	}
	if len(token) > 0 {
		opName := strings.Split(token, ".")[0]
		field := strings.Split(token, ".")[1]
		var value float64
		if field == "StartTime" {
			value = spans[opName].StartTime
		} else if field == "EndTime" {
			value = spans[opName].EndTime
		} else if field == "Duration" {
			value = spans[opName].Duration
		} else {
			return 0, fmt.Errorf("the filed %s is undefined", field)
		}
		formulaWithNumbers += strconv.FormatFloat(value, 'f', 6, 64)
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
