package autoconfigstrategies

import (
	"github.com/VahidMostofi/swarmmanager/internal/initializer"
	"github.com/VahidMostofi/swarmmanager/internal/strategies"
	"github.com/spf13/cobra"
)

var moboPythonPath string
var moboPythonScriptPath string

// MOBOCMD represents the BNV1 command (Bottle Neck Version 1)
var MOBOCMD = &cobra.Command{
	Use:   "mobo",
	Short: "Multi Objective Bayesian Optimization",
	Long:  `Under the name Multi Objective Bayesian Optimization, but basically any python code could be used. Now only response times are reported`, // TODO pass the whole info
	Run: func(cmd *cobra.Command, args []string) {

		strategy := &strategies.MultiObjectiveBayesianOptimization{
			PythonPath:       moboPythonPath,
			PythonScriptPath: moboPythonScriptPath,
			InitialConfig:    nil,
		}

		initializer.StartAutoconfig(strategy, "mobo")
	},
}

func init() {

	MOBOCMD.Flags().StringVar(&moboPythonPath, "python", "", "python path")
	MOBOCMD.Flags().StringVar(&moboPythonScriptPath, "script", "", "script path")

	cobra.MarkFlagRequired(MOBOCMD.Flags(), "python")
	cobra.MarkFlagRequired(MOBOCMD.Flags(), "script")
}
