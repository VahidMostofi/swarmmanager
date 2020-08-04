package jaeger

import (
	"fmt"
	"testing"
)

func TestFormulaEvaluation(t *testing.T) {
	value, err := evaluateJaegerFormula("auth.EndTime-   (auth_connect.EndTime+auth_connect.EndTime)", map[string]*span{"auth": &span{EndTime: 3}, "auth_connect": &span{EndTime: 1}})
	if err != nil {
		t.Errorf("error evaluateing formula: %w", err)
	}
	fmt.Println(value)
}
