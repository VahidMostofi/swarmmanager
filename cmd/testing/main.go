package main

import (
	"fmt"

	utils "github.com/VahidMostofi/swarmmanager/internal/statutils"
)

func main() {
	lower, upper, err := utils.ComputeToleranceInterval([]float64{1.3, 2.4, 5.4})
	if err != nil {
		panic(err)
	}
	fmt.Println(lower, upper)
}
