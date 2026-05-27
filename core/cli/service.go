/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : service.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Cobra service management commands (status, restart,
 *                diagnostics) for server and panel components.
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

	"abdal-4iproto-cli/core/paths"
	"abdal-4iproto-cli/core/service"
)

func newServiceCmd() *cobra.Command {
	var component string

	root := &cobra.Command{
		Use:   "service",
		Short: "Manage Abdal 4iProto systemd/sc services",
	}

	status := &cobra.Command{
		Use:   "status",
		Short: "Show service status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.Status(parseComponent(component))
		},
	}
	status.Flags().StringVar(&component, "component", "server", "server|panel")

	restart := &cobra.Command{
		Use:   "restart",
		Short: "Restart a service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return service.Restart(parseComponent(component))
		},
	}
	restart.Flags().StringVar(&component, "component", "server", "server|panel")

	diagnostics := &cobra.Command{
		Use:   "diagnostics",
		Short: "Run troubleshooting diagnostics",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := paths.InstallDir()
			if err != nil {
				return err
			}
			return service.Diagnostics(dir)
		},
	}

	root.AddCommand(status, restart, diagnostics)
	return root
}

func parseComponent(raw string) service.Component {
	switch raw {
	case "panel":
		return service.ComponentPanel
	default:
		return service.ComponentServer
	}
}
