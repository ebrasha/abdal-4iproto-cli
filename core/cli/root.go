/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : root.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Cobra root command wiring, banner/countdown bootstrap,
 *                and dispatch between interactive and flag-driven modes.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package cli

import (
	"os"

	"github.com/spf13/cobra"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/interactive"
	"abdal-4iproto-cli/core/ui"
	"abdal-4iproto-cli/core/updatecheck"
)

// Execute is the public entry used by main.go.
func Execute() error {
	if err := NewRoot().Execute(); err != nil {
		ui.ErrorBox("Error", err.Error())
		return err
	}
	return nil
}

// NewRoot builds the Cobra command tree.
func NewRoot() *cobra.Command {
	root := &cobra.Command{
		Use:   config.AppCommandName,
		Short: config.AppName + " installer and manager",
		Long:  config.AppName + " – install, configure, and manage the Abdal 4iProto server stack.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// No subcommand → interactive menu.
			return interactive.Run()
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// Bare binary with no subcommand → interactive menu prints its own banner.
		if cmd.Parent() == nil && cmd.Name() == config.AppCommandName && len(args) == 0 {
			return
		}
		ui.PrintBanner()
		updatecheck.Notify()
		if ui.RunCountdown(config.CountdownSeconds) {
			os.Exit(0)
		}
	}

	root.AddCommand(
		newInstallCmd(),
		newUninstallCmd(),
		newUserCmd(),
		newServiceCmd(),
		newConfigCmd(),
		newSelfInstallCmd(),
		newSelfUpdateCmd(),
		newHelpCmd(),
	)

	// Built-in Cobra help is customized via our help command as well.
	root.SetHelpCommand(newHelpCmd())
	return root
}
