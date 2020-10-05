package autoconfigstrategies

import (
	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/spf13/cobra"
)

// BruteCMD represents the cui command. CPU Usage Increase
var BruteCMD = &cobra.Command{
	Use:   "brute",
	Short: "Brute Force",
	Long:  `Brute Force`,
	Run: func(cmd *cobra.Command, args []string) {
		startBruteForce()
	},
}

func init() {

}

func startBruteForce() {
	initializer.StartBruteForce()
}
