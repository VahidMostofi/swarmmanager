package main

import (
	"fmt"

	utils "github.com/VahidMostofi/swarmmanager/internal/statutils"
)

func main() {
	lower, upper, err := utils.ComputeToleranceIntervalNonParametric([]float64{1.3, 2.4, 5.4}, 0.90, 0.90)
	if err != nil {
		panic(err)
	}
	fmt.Println(lower, upper)
}
