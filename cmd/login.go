package cmd

import (
	"fmt"
	"os"

	"github.com/mannkind/qnap-qsw/qnap"
	"github.com/spf13/cobra"
)

// Represents the ability to login to the QNAP QSW switch
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the QNAP QSW switch",
	Run: func(cmd *cobra.Command, args []string) {
		q := qnap.New(rootCmdOpts.Host)
		token, err := q.Login(rootCmdOpts.Password)
		if err != nil {
			fmt.Printf("Error logging into %s; %s\n", rootCmdOpts.Host, err)
			os.Exit(1)
			return
		}

		// Print the access token for use later
		fmt.Print(token)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
