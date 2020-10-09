package autoconfigstrategies

import (
	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/spf13/cobra"
)

var ppauSlaAgreementPropertyName string
var ppauSlaAgreementPropertyValue float64
var ppauStepSize float64
var ppauIndicator string
var ppauIsMultiContainer bool
var ppauDemandsFile string
var ppauIsConstantInit bool
var ppauConstantInitValue float64
var ppauDynamicFactor float64
var ppauMinStepSize float64

// PPAUCmd represents the PPAU command (Per Path Estimated Utilization)
var PPAUCmd = &cobra.Command{
	Use:   "ppau",
	Short: "Per Path Actual Utilization",
	Long:  `Per Path Actual Utilization description should be here. This approach has initial configuration too.`,
	Run: func(cmd *cobra.Command, args []string) {

		strategy := &strategies.PerPathActualUtilization{
			Agreements: []strategies.Agreement{
				{
					PropertyToConsider: ppauSlaAgreementPropertyName,
					Value:              ppauSlaAgreementPropertyValue,
				},
			},
			StepSize:          ppauStepSize,
			MultiContainer:    ppauIsMultiContainer,
			DemandsFilePath:   ppauDemandsFile,
			ConstantInit:      ppauConstantInitValue > 0.0,
			ConstantInitValue: ppauConstantInitValue,
		}

		strategy.Init()
		initializer.StartAutoconfig(strategy, "ppau")
	},
}

func init() {

	PPAUCmd.Flags().StringVar(&ppauSlaAgreementPropertyName, "property", "", "Which property of a run to consider for SLA? ResponseTimesMean, ResponseTimes90Percentile (95,99), ResponseTimes90Percentile")
	PPAUCmd.Flags().Float64Var(&ppauSlaAgreementPropertyValue, "value", 0, "The desired value related to SLA")
	PPAUCmd.Flags().Float64Var(&ppauStepSize, "stepsize", -1, "how much core to add at each step")
	PPAUCmd.Flags().BoolVar(&ppauIsMultiContainer, "mc", true, "run it with multiple containers or one fat container")
	PPAUCmd.Flags().StringVar(&ppauDemandsFile, "demands", "", "demand of each request on each service")
	PPAUCmd.Flags().Float64Var(&ppauConstantInitValue, "constantinit", -1, "if it's constant value for init, how much?")
	PPAUCmd.Flags().Float64Var(&ppauDynamicFactor, "dynamicfactor", 1, "How to change step size at each step")
	PPAUCmd.Flags().Float64Var(&ppauMinStepSize, "minstepsize", 0.1, "What should be the minimum step size")

	cobra.MarkFlagRequired(PPAUCmd.Flags(), "property")
	cobra.MarkFlagRequired(PPAUCmd.Flags(), "value")
	cobra.MarkFlagRequired(PPAUCmd.Flags(), "stepsize")
	cobra.MarkFlagRequired(PPAUCmd.Flags(), "mc")
	cobra.MarkFlagRequired(PPAUCmd.Flags(), "demands")
}
