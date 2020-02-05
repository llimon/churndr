// Copyright © 2019 Luis Limon
//

package cmd

import (
	"fmt"
	"os"

	"github.com/llimon/churndr/common"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "churndr",
	Short: "Monitors POD churn over a number of namespaces",
	Long: `Monitors and alerts when PODs are misbeheavings on specified namespaces. Generates alerts and deailed reports
    Provides "noise reduction" functionability
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return CheckRequiredFlags(cmd.Flags())
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.churndr.yaml)")
	rootCmd.PersistentFlags().BoolVar(&common.Config.Development, "development", false, "Enable development mode")
	rootCmd.PersistentFlags().BoolVar(&common.Config.NoAPIServer, "no-api-server", false, "Disable Rest API server")
	rootCmd.PersistentFlags().BoolVar(&common.Config.DissableEmailNotifications, "no-email-notifications", false, "Disable Email notifications")
	rootCmd.PersistentFlags().StringSliceVarP(&common.Config.Namespaces, "namespace", "n", []string{}, "")

	// Email parameters GO Here

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Settign the default deploymentMode to deployment_mode
	// Used to initialize the loggers

	/*
		deploymentMode := flag.Int("deployment_mode", 0, "deployment_mode")
		if *deploymentMode == 0 {
			log.Println("development")
		} else {
			log.Println("production")
		}
	*/

}
