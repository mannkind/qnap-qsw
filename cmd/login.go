package cmd

import (
	"fmt"

	"github.com/mannkind/qnap-qsw/qnap"
	"github.com/spf13/cobra"
)

// Represents the ability to login to the QNAP QSW switch
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the QNAP QSW 2K switch",
	Run: func(cmd *cobra.Command, args []string) {
		q := qnap.New(loginOpts.Host)
		token, err := q.Login(loginOpts.Password)
		if err != nil {
			fmt.Printf("Error logging into %s; %s\n", loginOpts.Host, err)
			return
		}

		// Print the access token for use later
		fmt.Print(token)
	},
}

var loginOpts = loginCmdOptions{}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVar(&loginOpts.Host, "host", "", "The host/ip")
	loginCmd.Flags().StringVar(&loginOpts.Password, "password", "", "The password of the admin user")
	loginCmd.MarkFlagRequired("host")
	loginCmd.MarkFlagRequired("password")
}
