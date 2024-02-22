package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/mannkind/qnap-qsw/logging"
	"github.com/mannkind/qnap-qsw/qnap"
	"github.com/spf13/cobra"
)

// Represents the ability to change the POE interface mode
var poeModeOpts = poeModeCmdOptions{}
var poeModeCmd = &cobra.Command{
	Use:   "poeMode",
	Short: "Change the POE port mode",
	Run: func(cmd *cobra.Command, args []string) {
		log := logging.New(rootCmdOpts.Verbosity)
		q := qnap.NewWithToken(rootCmdOpts.Host, poeModeOpts.Token)
		_, err := q.Login(rootCmdOpts.Password)
		if err != nil {
			log.Error(err, "Error logging into switch", "host", rootCmdOpts.Host)
			os.Exit(1)
			return
		}

		// Get the current POE interface statuses
		knownInterfaces, err := q.POEInterfaces()
		if err != nil {
			log.Error(err, "Error fetching known interfaces from switch", "host", rootCmdOpts.Host)
			os.Exit(1)
			return
		}

		// Set the POE interface statuses
		wg := sync.WaitGroup{}
		allPortsCount := len(poeModeOpts.DisablePorts) + len(poeModeOpts.PoePorts) + len(poeModeOpts.PoePlusPorts) + len(poeModeOpts.PoePlusPlusPorts)
		ch := make(chan error, allPortsCount)
		portsAndModes := []struct {
			ports []string
			mode  string
		}{
			{ports: poeModeOpts.DisablePorts, mode: "disabled"},
			{ports: poeModeOpts.PoePorts, mode: "poe"},
			{ports: poeModeOpts.PoePlusPorts, mode: "poe+"},
			{ports: poeModeOpts.PoePlusPlusPorts, mode: "poe++"},
		}

		for _, portsAndMode := range portsAndModes {
			for _, port := range portsAndMode.ports {
				// Skip empty ports
				if port == "" {
					continue
				}

				portIdx := port
				properties, ok := knownInterfaces[portIdx]
				if !ok {
					continue
				}

				properties.Mode = qnap.POEModes.Unknown.FromString(portsAndMode.mode)
				wg.Add(1)
				go func() {
					defer wg.Done()
					ch <- q.UpdatePOEInterfaces(portIdx, properties)
				}()
			}
		}

		wg.Wait()
		close(ch)

		errors := false
		for err := range ch {
			if err == nil {
				continue
			}

			log.Error(err, "Error updating POE interface on switch", "host", rootCmdOpts.Host)
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
	poeModeCmd.Flags().StringSliceVar(&poeModeOpts.DisablePorts, "disable-ports", []string{}, "The ports to disable")
	poeModeCmd.Flags().StringSliceVar(&poeModeOpts.PoePorts, "poe-ports", []string{}, "The ports to enable poe")
	poeModeCmd.Flags().StringSliceVar(&poeModeOpts.PoePlusPorts, "poeplus-ports", []string{}, "The ports to enable poe+")
	poeModeCmd.Flags().StringSliceVar(&poeModeOpts.PoePlusPlusPorts, "poeplusplus-ports", []string{}, "The ports to poe++")
}
