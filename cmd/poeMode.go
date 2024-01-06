package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/mannkind/qnap-qsw/qnap"
	"github.com/spf13/cobra"
)

// Represents the ability to change the POE interface mode
var poeModeOpts = poeModeCmdOptions{}
var poeModeCmd = &cobra.Command{
	Use:   "poeMode",
	Short: "Change the POE port mode",
	Run: func(cmd *cobra.Command, args []string) {
		q := qnap.NewWithToken(rootCmdOpts.Host, poeModeOpts.Token)
		_, err := q.Login(rootCmdOpts.Password)
		if err != nil {
			fmt.Printf("Error logging in to qnap qsw; %s\n", err)
			os.Exit(1)
			return
		}

		// Get the current POE interface statuses
		knownInterfaces, err := q.POEInterfaces()
		if err != nil {
			fmt.Printf("Error fetching known interfaces from qnap qsw; %s\n", err)
			os.Exit(1)
			return
		}

		// Set the POE interface statuses
		wg := sync.WaitGroup{}
		ch := make(chan error, len(poeModeOpts.Ports))
		for _, port := range poeModeOpts.Ports {
			portIdx := port
			properties, ok := knownInterfaces[portIdx]
			if !ok {
				continue
			}

			properties.Mode = qnap.POEModes.Unknown.FromString(poeModeOpts.Mode)
			wg.Add(1)
			go func() {
				ch <- q.UpdatePOEInterfaces(&wg, portIdx, properties)
			}()
		}

		wg.Wait()
		close(ch)

		errors := false
		for err := range ch {
			if err == nil {
				continue
			}

			fmt.Printf("Error updating POE interface on qnap qsw; %s\n", err)
			errors = true
		}

		if errors {
			os.Exit(1)
			return
		}

		fmt.Print("OK")
	},
}

func init() {
	rootCmd.AddCommand(poeModeCmd)

	poeModeCmd.Flags().StringVar(&poeModeOpts.Token, "token", "", "The token representing the admin user; use this or password")
	poeModeCmd.Flags().StringSliceVar(&poeModeOpts.Ports, "ports", []string{}, "The ports to modify")
	poeModeCmd.Flags().StringVar(&poeModeOpts.Mode, "mode", "disable", "The mode of the ports to modify")
}
