package autoconfigstrategies

import (
	"log"
	"os"

	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var adfcSlaAgreementPropertyName string
var adfcSlaAgreementPropertyValue float64
var adfcStepSize float64
var adfcIndicator string
var adfcIsMultiContainer bool

// ADFCCmd represents the ADFC command
var ADFCCmd = &cobra.Command{
	Use:   "adfc",
	Short: "Add Diffrent Fractional CPU",
	Long:  `Add Diffrent Fractional CPU proportional to Estimated CPU Utilizations which are normalized over the whole system`,
	Run: func(cmd *cobra.Command, args []string) {
		values, maxIncrease, err := strategies.GetFractionalCPUIncreaseValues(viper.GetString("workloadStr"), adfcIndicator, adfcStepSize)
		if err != nil {
			log.Panic(err)
			os.Exit(1)
		}
		// log.Println("values for Fractional Increase:", values)
		// log.Println("values for Max Increase:", maxIncrease)
		// log.Println("Configuring AddFractionalCPUcores with Value:", *adfccThreshold, "and property of", *adfccValueName, " and core amount of", *adfccAmount, "and indicator=", *adfccIndicator)
		strategy := &strategies.AddDifferentFractionalCPUcores{
			ServiceToAmount:   values,
			MaxServiceIncease: maxIncrease,
			Agreements: []strategies.Agreement{
				{
					PropertyToConsider: adfcSlaAgreementPropertyName,
					Value:              adfcSlaAgreementPropertyValue,
				},
			},
			MultiContainer: adfcIsMultiContainer,
		}

		initializer.StartAutoconfig(strategy, "adfc")
	},
}

func init() {

	ADFCCmd.Flags().StringVar(&adfcSlaAgreementPropertyName, "property", "", "Which property of a run to consider for SLA? ResponseTimesMean, ResponseTimes90Percentile (95,99), ResponseTimes90Percentile")
	ADFCCmd.Flags().Float64Var(&adfcSlaAgreementPropertyValue, "value", 0, "The desired value related to SLA")
	ADFCCmd.Flags().Float64Var(&adfcStepSize, "stepsize", -1, "how much core to add at each step")
	adfcIndicator = "Utilization" // previously this could also be demand
	ADFCCmd.Flags().BoolVar(&adfcIsMultiContainer, "mc", true, "run it with multiple containers or one fat container")
	cobra.MarkFlagRequired(ADFCCmd.Flags(), "property")
	cobra.MarkFlagRequired(ADFCCmd.Flags(), "value")
	cobra.MarkFlagRequired(ADFCCmd.Flags(), "stepsize")
	cobra.MarkFlagRequired(ADFCCmd.Flags(), "mc")

}
