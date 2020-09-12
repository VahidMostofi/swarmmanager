package autoconfigstrategies

import (
	"log"
	"os"

	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var valueName string
var valueTreshold float64
var stepSize float64
var indicator string
var isMultiContainer bool

// ADFCCmd represents the ADFC command
var ADFCCmd = &cobra.Command{
	Use:   "adfc",
	Short: "Add Diffrent Fractional CPU",
	Long:  `Add Diffrent Fractional CPU proportional to Estimated CPU Utilizations which are normalized over the whole system`,
	Run: func(cmd *cobra.Command, args []string) {
		values, maxIncrease, err := strategies.GetFractionalCPUIncreaseValues(viper.GetString("workloadStr"), indicator, stepSize)
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
					PropertyToConsider: valueName,
					Value:              valueTreshold,
				},
			},
			MultiContainer: isMultiContainer,
		}

		initializer.StartAutoconfig(strategy, "adfc")
	},
}

func init() {

	ADFCCmd.Flags().StringVar(&valueName, "property", "", "Which property of a run to consider? CPUUsageMean,CPUUsage90Percentile 70-95, 99")
	ADFCCmd.Flags().Float64Var(&valueTreshold, "value", 0, "Which property of a run to consider? CPUUsageMean,CPUUsage90Percentile 70-95, 99")
	ADFCCmd.Flags().Float64Var(&stepSize, "stepsize", -1, "how much core to add at each step")
	indicator = "Utilization" // previously this could also be demand
	ADFCCmd.Flags().BoolVar(&isMultiContainer, "mc", true, "run it with multiple containers or one fat container")
	cobra.MarkFlagRequired(ADFCCmd.Flags(), "property")
	cobra.MarkFlagRequired(ADFCCmd.Flags(), "value")
	cobra.MarkFlagRequired(ADFCCmd.Flags(), "stepsize")
	cobra.MarkFlagRequired(ADFCCmd.Flags(), "mc")

}
