package autoconfigstrategies

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ADFCCmd represents the ADFC command
var ADFCCmd = &cobra.Command{
	Use:   "adfc",
	Short: "Add Diffrent Fractional CPU",
	Long:  `Add Diffrent Fractional CPU proportional to Estimated CPU Utilizations which are normalized over the whole system`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("you want to optimize using ADFC")
		//TODO nex to go into cmd and
	},
}

func init() {
	// adfcCmd.Flags()
}
