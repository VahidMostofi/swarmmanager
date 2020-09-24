package autoconfigstrategies

import (
	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/spf13/cobra"
)

var ppeSlaAgreementPropertyName string
var ppeSlaAgreementPropertyValue float64
var ppeStepSize float64
var ppeIndicator string
var ppeIsMultiContainer bool
var ppeDemandsFile string
var ppeIsConstantInit bool
var ppeConstantInitValue float64

// PPECmd represents the PPE command (Per Path Estimated Utilization)
var PPECmd = &cobra.Command{
	Use:   "ppe",
	Short: "Per Path Equal",
	Long:  `Per Path Equal. Equally shre stepSize CPU between all services in the reqeust path`,
	Run: func(cmd *cobra.Command, args []string) {

		strategy := &strategies.PerPathEqual{
			Agreements: []strategies.Agreement{
				{
					PropertyToConsider: ppeSlaAgreementPropertyName,
					Value:              ppeSlaAgreementPropertyValue,
				},
			},
			StepSize:        ppeStepSize,
			MultiContainer:  ppeIsMultiContainer,
			DemandsFilePath: ppeDemandsFile,
			ConstantInit: ppeConstantInitValue > 0.0,
			ConstantInitValue: ppeConstantInitValue,
		}

		strategy.Init()
		initializer.StartAutoconfig(strategy, "ppe")
	},
}

func init() {

	PPECmd.Flags().StringVar(&ppeSlaAgreementPropertyName, "property", "", "Which property of a run to consider for SLA? ResponseTimesMean, ResponseTimes90Percentile (95,99), ResponseTimes90Percentile")
	PPECmd.Flags().Float64Var(&ppeSlaAgreementPropertyValue, "value", 0, "The desired value related to SLA")
	PPECmd.Flags().Float64Var(&ppeStepSize, "stepsize", -1, "how much core to add at each step")
	PPECmd.Flags().BoolVar(&ppeIsMultiContainer, "mc", true, "run it with multiple containers or one fat container")
	PPECmd.Flags().StringVar(&ppeDemandsFile, "demands", "", "demand of each request on each service")
	PPECmd.Flags().Float64Var(&ppeConstantInitValue, "constantinit", -1, "if it's constant value for init, how much?")

	cobra.MarkFlagRequired(PPECmd.Flags(), "property")
	cobra.MarkFlagRequired(PPECmd.Flags(), "value")
	cobra.MarkFlagRequired(PPECmd.Flags(), "stepsize")
	cobra.MarkFlagRequired(PPECmd.Flags(), "mc")
	cobra.MarkFlagRequired(PPECmd.Flags(), "demands")
}
