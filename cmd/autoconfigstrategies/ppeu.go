package autoconfigstrategies

import (
	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/spf13/cobra"
)

var ppeuSlaAgreementPropertyName string
var ppeuSlaAgreementPropertyValue float64
var ppeuStepSize float64
var ppeuIndicator string
var ppeuIsMultiContainer bool
var ppeuDemandsFile string

// PPEUCmd represents the PPEU command (Per Path Estimated Utilization)
var PPEUCmd = &cobra.Command{
	Use:   "ppeu",
	Short: "Per Path Estimated Utilization",
	Long:  `Per Path Estimated Utilization description should be here. This approach has initial configuration too.`,
	Run: func(cmd *cobra.Command, args []string) {

		strategy := &strategies.PerPathEstimatedUtilization{
			Agreements: []strategies.Agreement{
				{
					PropertyToConsider: ppeuSlaAgreementPropertyName,
					Value:              ppeuSlaAgreementPropertyValue,
				},
			},
			StepSize:        ppeuStepSize,
			MultiContainer:  ppeuIsMultiContainer,
			DemandsFilePath: ppeuDemandsFile,
		}

		strategy.Init()
		initializer.StartAutoconfig(strategy, "ppeu")
	},
}

func init() {

	PPEUCmd.Flags().StringVar(&ppeuSlaAgreementPropertyName, "property", "", "Which property of a run to consider for SLA? ResponseTimesMean, ResponseTimes90Percentile (95,99), ResponseTimes90Percentile")
	PPEUCmd.Flags().Float64Var(&ppeuSlaAgreementPropertyValue, "value", 0, "The desired value related to SLA")
	PPEUCmd.Flags().Float64Var(&ppeuStepSize, "stepsize", -1, "how much core to add at each step")
	PPEUCmd.Flags().BoolVar(&ppeuIsMultiContainer, "mc", true, "run it with multiple containers or one fat container")
	PPEUCmd.Flags().StringVar(&ppeuDemandsFile, "demands", "", "demand of each request on each service")

	cobra.MarkFlagRequired(PPEUCmd.Flags(), "property")
	cobra.MarkFlagRequired(PPEUCmd.Flags(), "value")
	cobra.MarkFlagRequired(PPEUCmd.Flags(), "stepsize")
	cobra.MarkFlagRequired(PPEUCmd.Flags(), "mc")
	cobra.MarkFlagRequired(PPEUCmd.Flags(), "demands")
}
