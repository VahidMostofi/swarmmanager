package autoconfigstrategies

import (
	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/spf13/cobra"
)

var cuiCPUUtilizationStatistics string
var cuiCPUUtilizationThreshold float64

// CUICmd represents the cui command. CPU Usage Increase
var CUICmd = &cobra.Command{
	Use:   "cui",
	Short: "CPU Usage Increase",
	Long:  `Add one more instance to every service which has more than a specific CPU utilization`,
	Run: func(cmd *cobra.Command, args []string) {

		strategy := &strategies.CPUUsageIncrease{
			ValueToConsider: cuiCPUUtilizationStatistics,
			Threshold:       cuiCPUUtilizationThreshold,
		}

		initializer.StartAutoconfig(strategy, "cui")
	},
}

func init() {
	CUICmd.Flags().StringVar(&cuiCPUUtilizationStatistics, "cpuStat", "", "Which statistics of the CPU we should consider for scaling. CPUUsageMean,CPUUsage90Percentile 70,75,...,95,99")
	CUICmd.Flags().Float64Var(&cuiCPUUtilizationThreshold, "cpuThreshold", 0, "What is the threshold for scaling up service")

	cobra.MarkFlagRequired(CUICmd.Flags(), "cpuStat")
	cobra.MarkFlagRequired(CUICmd.Flags(), "cpuThreshold")
}
