package cmd

import (
	"fmt"
	"sync"

	"github.com/mannkind/qnap-qsw/qnap"
	"github.com/spf13/cobra"
)

// Represents the ability to change the POE interface mode
var poeModeCmd = &cobra.Command{
	Use:   "poeMode",
	Short: "Change the POE port mode",
	Run: func(cmd *cobra.Command, args []string) {
		q := qnap.NewWithToken(poeModeOpts.Host, poeModeOpts.Token)
		_, err := q.Login(poeModeOpts.Password)
		if err != nil {
			fmt.Printf("Error logging in to qnap qsw; %s\n", err)
			return
		}

		// Get the current POE interface statuses
		knownInterafces, err := q.POEInterfaces()
		if err != nil {
			fmt.Printf("Error fetching known interfaces from qnap qsw; %s\n", err)
			return
		}

		// Set the POE interface statuses
		wg := sync.WaitGroup{}
		ch := make(chan error, len(poeModeOpts.Ports))
		for _, port := range poeModeOpts.Ports {
			properties, ok := knownInterafces[port]
			if !ok {
				continue
			}

			properties.Mode = poeModeOpts.Mode
			wg.Add(1)
			go func(portIdx string) {
				ch <- q.UpdatePOEInterfaces(&wg, portIdx, properties)
			}(port)
		}

		wg.Wait()
		close(ch)

		errors := false
		for err := range ch {
			if err == nil {
				continue
			}

			fmt.Printf("Error updating POE interface on qnap qsw; %s\n", err)
		}

		if !errors {
			fmt.Print("OK")
		}
	},
}

var poeModeOpts = poeModeCmdOptions{}

func init() {
	rootCmd.AddCommand(poeModeCmd)

	poeModeCmd.Flags().StringVar(&poeModeOpts.Host, "host", "", "The host/ip")
	poeModeCmd.Flags().StringVar(&poeModeOpts.Password, "password", "", "The password of the admin user; use this or token")
	poeModeCmd.Flags().StringVar(&poeModeOpts.Token, "token", "", "The token representing the admin user; use this or password")
	poeModeCmd.Flags().StringSliceVar(&poeModeOpts.Ports, "ports", []string{}, "The ports to modify")
	poeModeCmd.Flags().StringVar(&poeModeOpts.Mode, "mode", "disable", "The mode of the ports to modify")
	poeModeCmd.MarkFlagRequired("host")
}
