package autoconfigstrategies

import (
	"log"
	"strings"

	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/spf13/cobra"
)

var demandsDuration int
var demandsResultPath string

// DemandsCmd represents the demands command, it finds the demand of each service for each request
var DemandsCmd = &cobra.Command{
	Use:   "demands",
	Short: "Finds Demands",
	Long:  `Finds Demands for each service for each request`,
	Run: func(cmd *cobra.Command, args []string) {
		if !(strings.HasSuffix(demandsResultPath, ".yml") || strings.HasSuffix(demandsResultPath, ".yaml")) {
			log.Panic("error: the resultpath must be .yaml or .yml")
		}

		configs.GetConfig().Test.Duration = demandsDuration

		strategy := &strategies.DemandsFinder{
			ResultPath: demandsResultPath,
		}

		initializer.StartAutoconfig(strategy, "demands")
	},
}

func init() {

	DemandsCmd.Flags().StringVar(&demandsResultPath, "resultpath", "", "Where to store demands? must be .yml or .yaml")
	DemandsCmd.Flags().IntVar(&demandsDuration, "duration", 0, "Duration of simulation for finding the demands")
	cobra.MarkFlagRequired(DemandsCmd.Flags(), "resultpath")
	cobra.MarkFlagRequired(DemandsCmd.Flags(), "duration")

}
