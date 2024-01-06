package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// Represents the base command when called without any subcommands
var rootCmdOpts = rootCommandOptions{}
var rootCmd = &cobra.Command{
	Use:   "qnap-qsw",
	Short: "The QNAP QSW tool",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootCmdOpts.Host, "host", "switch.lan", "The host/ip")
	rootCmd.PersistentFlags().StringVar(&rootCmdOpts.Password, "password", os.Getenv("QNAP_QSW_PASSWORD"), "The password of the admin user (default: $QNAP_QSW_PASSWORD)")
}
