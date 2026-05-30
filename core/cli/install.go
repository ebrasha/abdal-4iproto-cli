/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : install.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Cobra install subcommand with flags for ports, panel
 *                credentials, keygen parameters, and component scope.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/installer"
	"abdal-4iproto-cli/core/keygen"
	"abdal-4iproto-cli/core/network"
)

func newInstallCmd() *cobra.Command {
	var (
		serverOnly    bool
		panelOnly     bool
		serverPorts   string
		panelPort     int
		panelUser     string
		panelPass     string
		keyType       string
		keyBits       int
		keyForce      bool
		keyFile       string
		noServices    bool
		forceFresh    bool
	)

	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install Abdal 4iProto server stack components",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := installer.DefaultOptions()
			opts.PanelPort = panelPort
			opts.PanelUsername = panelUser
			opts.PanelPassword = panelPass
			opts.InstallServices = !noServices
			opts.Force = forceFresh
			opts.Keygen = keygen.Options{
				Type: keyType, Bits: keyBits, Force: keyForce, OutputFile: keyFile,
			}

			if serverOnly && panelOnly {
				return fmt.Errorf("cannot use --server-only and --panel-only together")
			}
			switch {
			case serverOnly:
				opts.Target = installer.TargetServer
			case panelOnly:
				opts.Target = installer.TargetPanel
			default:
				opts.Target = installer.TargetAll
			}

			if serverPorts != "" {
				ports, err := network.ParsePortList(serverPorts)
				if err != nil {
					return err
				}
				if err := network.ValidatePorts(ports); err != nil {
					return fmt.Errorf("server ports validation failed: %w", err)
				}
				opts.ServerPorts = ports
			}

			if opts.Target == installer.TargetAll || opts.Target == installer.TargetPanel {
				if !network.IsPortAvailable(opts.PanelPort) {
					return &network.PortCheckError{Port: opts.PanelPort}
				}
			}

			if err := installer.Run(opts); err != nil {
				if errors.Is(err, installer.ErrAlreadyInstalled) {
					return fmt.Errorf(
						"%w for the requested scope. Re-run with --force to wipe only the selected component(s) and reinstall from scratch",
						err,
					)
				}
				return err
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&serverOnly, "server-only", false, "Install server component only")
	cmd.Flags().BoolVar(&panelOnly, "panel-only", false, "Install panel component only")
	cmd.Flags().StringVar(&serverPorts, "server-ports", "", "Comma-separated server listener ports")
	cmd.Flags().IntVar(&panelPort, "panel-port", config.DefaultPanelPort, "Panel HTTP port")
	cmd.Flags().StringVar(&panelUser, "panel-username", config.DefaultPanelUsername, "Panel login username")
	cmd.Flags().StringVar(&panelPass, "panel-password", config.DefaultPanelPassword, "Panel login password")
	cmd.Flags().StringVar(&keyType, "key-type", config.DefaultKeygenType, "SSH key type: rsa|ed25519|ecdsa")
	cmd.Flags().IntVar(&keyBits, "key-bits", config.DefaultKeygenBits, "SSH key size in bits")
	cmd.Flags().BoolVar(&keyForce, "key-force", config.DefaultKeygenForce, "Overwrite existing SSH key files")
	cmd.Flags().StringVar(&keyFile, "key-file", config.DefaultKeyBaseName, "SSH private key output filename")
	cmd.Flags().BoolVar(&noServices, "no-services", false, "Skip systemd/sc service registration")
	cmd.Flags().BoolVar(&forceFresh, "force", false, "Wipe an existing installation and perform a fresh install")

	return cmd
}
