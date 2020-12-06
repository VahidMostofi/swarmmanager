package autoconfigstrategies

import (
	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/spf13/cobra"
)

var bnv2SlaAgreementPropertyName string
var bnv2SlaAgreementPropertyValue float64
var bnv2StepSize float64
var bnv2Indicator string
var bnv2IsMultiContainer bool
var bnv2DemandsFile string
var bnv2MinimumStepSize float64

// BNV2Cmd represents the BNV2 command (Bottle Neck Version 2)
var BNV2Cmd = &cobra.Command{
	Use:   "bnv2",
	Short: "Bottle Neck Version 2",
	Long:  `Bottle Neck Version 2 description should be here. Initial Values are EU. At each step, if the path doesnt meet SLA, adds EU * StepSize to the bottleneck only`,
	Run: func(cmd *cobra.Command, args []string) {

		strategy := &strategies.BottleNeckOnlyVersion2{
			Agreements: []strategies.Agreement{
				{
					PropertyToConsider: bnv2SlaAgreementPropertyName,
					Value:              bnv2SlaAgreementPropertyValue,
				},
			},
			StepSize:          bnv2StepSize,
			MultiContainer:    bnv2IsMultiContainer,
			DemandsFilePath:   bnv2DemandsFile,
			ConstantInit:      true,
			ConstantInitValue: 0.5,
			MinimumStepSize:   bnv2MinimumStepSize,
		}

		strategy.Init()
		initializer.StartAutoconfig(strategy, "bnv2")
	},
}

func init() {

	BNV2Cmd.Flags().StringVar(&bnv2SlaAgreementPropertyName, "property", "", "Which property of a run to consider for SLA? ResponseTimesMean, ResponseTimes90Percentile (95,99), ResponseTimes90Percentile")
	BNV2Cmd.Flags().Float64Var(&bnv2SlaAgreementPropertyValue, "value", 0, "The desired value related to SLA")
	BNV2Cmd.Flags().Float64Var(&bnv2StepSize, "stepsize", -1, "how much core to add at each step")
	BNV2Cmd.Flags().BoolVar(&bnv2IsMultiContainer, "mc", true, "run it with multiple containers or one fat container")
	BNV2Cmd.Flags().StringVar(&bnv2DemandsFile, "demands", "", "demand of each request on each service")
	BNV2Cmd.Flags().Float64Var(&bnv2MinimumStepSize, "minstepsize", 0.2, "minimum step size")

	cobra.MarkFlagRequired(BNV2Cmd.Flags(), "property")
	cobra.MarkFlagRequired(BNV2Cmd.Flags(), "value")
	cobra.MarkFlagRequired(BNV2Cmd.Flags(), "stepsize")
	cobra.MarkFlagRequired(BNV2Cmd.Flags(), "mc")
	cobra.MarkFlagRequired(BNV2Cmd.Flags(), "demands")
}
