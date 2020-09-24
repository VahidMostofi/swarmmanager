package autoconfigstrategies

import (
	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/spf13/cobra"
)

var bnv1SlaAgreementPropertyName string
var bnv1SlaAgreementPropertyValue float64
var bnv1StepSize float64
var bnv1Indicator string
var bnv1IsMultiContainer bool
var bnv1DemandsFile string

// BNV1Cmd represents the BNV1 command (Bottle Neck Version 1)
var BNV1Cmd = &cobra.Command{
	Use:   "bnv1",
	Short: "Bottle Neck Version 1",
	Long:  `Bottle Neck Version 1 description should be here. Initial Values are EU. At each step, if the path doesnt meet SLA, adds EU * StepSize to the bottleneck only`,
	Run: func(cmd *cobra.Command, args []string) {

		strategy := &strategies.BottleNeckOnlyVersion1{
			Agreements: []strategies.Agreement{
				{
					PropertyToConsider: bnv1SlaAgreementPropertyName,
					Value:              bnv1SlaAgreementPropertyValue,
				},
			},
			StepSize:        bnv1StepSize,
			MultiContainer:  bnv1IsMultiContainer,
			DemandsFilePath: bnv1DemandsFile,
		}

		strategy.Init()
		initializer.StartAutoconfig(strategy, "bnv1")
	},
}

func init() {

	BNV1Cmd.Flags().StringVar(&bnv1SlaAgreementPropertyName, "property", "", "Which property of a run to consider for SLA? ResponseTimesMean, ResponseTimes90Percentile (95,99), ResponseTimes90Percentile")
	BNV1Cmd.Flags().Float64Var(&bnv1SlaAgreementPropertyValue, "value", 0, "The desired value related to SLA")
	BNV1Cmd.Flags().Float64Var(&bnv1StepSize, "stepsize", -1, "how much core to add at each step")
	BNV1Cmd.Flags().BoolVar(&bnv1IsMultiContainer, "mc", true, "run it with multiple containers or one fat container")
	BNV1Cmd.Flags().StringVar(&bnv1DemandsFile, "demands", "", "demand of each request on each service")

	cobra.MarkFlagRequired(BNV1Cmd.Flags(), "property")
	cobra.MarkFlagRequired(BNV1Cmd.Flags(), "value")
	cobra.MarkFlagRequired(BNV1Cmd.Flags(), "stepsize")
	cobra.MarkFlagRequired(BNV1Cmd.Flags(), "mc")
	cobra.MarkFlagRequired(BNV1Cmd.Flags(), "demands")
}
