package autoconfigstrategies

import (
	"fmt"
	"log"
	"strconv"

	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/VahidMostofi/swarmmanager/internal/swarm"
	"github.com/spf13/cobra"
)

// SingleCmd represents the cui command. CPU Usage Increase
var SingleCmd = &cobra.Command{
	Use:   "single",
	Short: "Single run",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config := make(map[string]swarm.SimpleSpecs)
		serviceName := ""
		for i, arg := range args {
			if i%2 == 0 {
				serviceName = arg
			} else {
				cpu, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					log.Panic(fmt.Errorf("cant parse %s for %s", arg, serviceName))
				}
				config[serviceName] = swarm.SimpleSpecs{CPU: cpu, Replica: 1, Worker: 1}
				log.Println(serviceName, cpu)
			}
		}

		strategy := &strategies.SingleRun{Config: config}

		initializer.StartAutoconfig(strategy, "single")
	},
}
