/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : config.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Cobra configuration commands for server_config.json
 *                and abdal-4iproto-panel.json.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package cli

import (
	"github.com/spf13/cobra"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/configmgr"
	"abdal-4iproto-cli/core/network"
)

func newConfigCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "config",
		Short: "Update server or panel JSON configuration",
	}

	var serverPorts string
	server := &cobra.Command{
		Use:   "server",
		Short: "Update server_config.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			ports, err := network.ParsePortList(serverPorts)
			if err != nil {
				return err
			}
			return configmgr.UpdateServerPorts(ports)
		},
	}
	server.Flags().StringVar(&serverPorts, "ports", "", "Comma-separated listener ports (required)")
	_ = server.MarkFlagRequired("ports")

	var (
		panelPort  int
		panelUser  string
		panelPass  string
	)
	panel := &cobra.Command{
		Use:   "panel",
		Short: "Update abdal-4iproto-panel.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			return configmgr.UpdatePanelConfig(panelPort, panelUser, panelPass, nil)
		},
	}
	panel.Flags().IntVar(&panelPort, "port", config.DefaultPanelPort, "Panel HTTP port")
	panel.Flags().StringVar(&panelUser, "username", "", "Panel username (optional)")
	panel.Flags().StringVar(&panelPass, "password", "", "Panel password (optional)")

	root.AddCommand(server, panel)
	return root
}
