package cmd

import (
	"fmt"
	"os"

	strategiesCmd "github.com/VahidMostofi/swarmmanager/cmd/autoconfigstrategies"
	"github.com/VahidMostofi/swarmmanager/configs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var appName string
var workloadStr string
var testName string

// autoconfigCmd represents the autoconfig command
var autoconfigCmd = &cobra.Command{
	Use:   "autoconfig",
	Short: "Automatically configure an application which is CPU-intensive while running multiple tests agains a workload",
	Long:  `Automatically configure an application which is CPU-intensive while running multiple tests agains a workload`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(autoconfigCmd)
	cobra.OnInitialize(initConfig)

	autoconfigCmd.PersistentFlags().StringVar(&appName, "appname", "", "the name of the application we want to auto-configure")
	viper.BindPFlag("appname", autoconfigCmd.PersistentFlags().Lookup("appname"))
	cobra.MarkFlagRequired(autoconfigCmd.PersistentFlags(), "appname")

	autoconfigCmd.PersistentFlags().StringVar(&workloadStr, "workload", "", "the key specifying workload for load-generator")
	viper.BindPFlag("workloadStr", autoconfigCmd.PersistentFlags().Lookup("workload"))
	cobra.MarkFlagRequired(autoconfigCmd.PersistentFlags(), "workloadStr")

	autoconfigCmd.PersistentFlags().StringVar(&testName, "testName", "", "name of the test")
	viper.BindPFlag("testName", autoconfigCmd.PersistentFlags().Lookup("testName"))
	cobra.MarkFlagRequired(autoconfigCmd.PersistentFlags(), "testName")

	autoconfigCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/configs/{appname}.yaml)")

	// Every new strategy must be added here
	autoconfigCmd.AddCommand(strategiesCmd.ADFCCmd)
	autoconfigCmd.AddCommand(strategiesCmd.CUICmd)
	autoconfigCmd.AddCommand(strategiesCmd.DemandsCmd)
	autoconfigCmd.AddCommand(strategiesCmd.PPEUCmd)
	autoconfigCmd.AddCommand(strategiesCmd.PPAUCmd)
	autoconfigCmd.AddCommand(strategiesCmd.PPECmd)
	autoconfigCmd.AddCommand(strategiesCmd.BNV1Cmd)
	autoconfigCmd.AddCommand(strategiesCmd.BNV2Cmd)
	autoconfigCmd.AddCommand(strategiesCmd.BruteCMD)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if appName == "" {
		fmt.Println("you must specify the appname, run: swarmmanager autoconfig --help")
		os.Exit(1)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		mydir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(mydir + "/configurations/")
		viper.SetConfigName(appName)
	}

	// viper.AutomaticEnv() // uncomment read in environment variables that match

	if workloadStr == "" {
		fmt.Println("you must specify the workload, run: swarmmanager autoconfig --help")
		os.Exit(1)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
	configs.Initialize()
}
