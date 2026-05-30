/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : prompt.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Interactive Survey prompts that collect installation
 *                parameters including ports, panel credentials, and keygen.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package installer

import (
	"fmt"
	"strings"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/keygen"
	"abdal-4iproto-cli/core/network"
	"abdal-4iproto-cli/core/ui"
)

// PromptTarget asks the operator which install scope they want and
// returns Options with only Target populated. A "Back" selection bubbles
// up as ui.ErrUserBack so the parent menu can return quietly without an
// extra confirmation step.
func PromptTarget(base Options) (Options, error) {
	targetChoice, err := ui.AskSelect("Installation scope", []string{
		"Full stack (Server + Panel + KeyGen)",
		"Server only",
		"Panel only",
		"Back",
	}, "Full stack (Server + Panel + KeyGen)")
	if err != nil {
		return base, err
	}
	if targetChoice == "Back" {
		return base, ui.ErrUserBack
	}
	switch targetChoice {
	case "Server only":
		base.Target = TargetServer
	case "Panel only":
		base.Target = TargetPanel
	default:
		base.Target = TargetAll
	}
	return base, nil
}

// PromptOptions interactively collects installation parameters. When
// the operator picks the "Back" entry the function returns the
// ui.ErrUserBack sentinel so the caller can quietly bubble up to the
// previous menu without asking any extra confirmation.
//
// This is a thin wrapper around PromptTarget + PromptInstallDetails kept
// for backward compatibility with callers that prefer a single entry.
func PromptOptions(base Options) (Options, error) {
	base, err := PromptTarget(base)
	if err != nil {
		return base, err
	}
	return PromptInstallDetails(base)
}

// PromptInstallDetails collects every install parameter except the
// target, which must already be set on base. Splitting this out lets the
// interactive menu confirm a scope-aware fresh install *before* asking
// for ports, panel credentials or keygen options.
func PromptInstallDetails(base Options) (Options, error) {
	if base.Target == TargetAll || base.Target == TargetServer {
		suggested, err := network.SuggestFreePorts(config.ServerSuggestedCount)
		if err != nil {
			return base, err
		}
		defaultPorts := network.FormatPortList(suggested)
		ui.Info("Suggested free server ports: " + defaultPorts)

		for {
			raw, err := ui.AskString("Server ports (comma-separated)", defaultPorts, true)
			if err != nil {
				return base, err
			}
			ports, err := network.ParsePortList(raw)
			if err != nil {
				ui.Warning(err.Error())
				continue
			}
			if err := network.ValidatePorts(ports); err != nil {
				if pce, ok := err.(*network.PortCheckError); ok {
					ui.Warning(fmt.Sprintf("Port %d is reserved/in use. Choose different ports.", pce.Port))
					continue
				}
				ui.Warning(err.Error())
				continue
			}
			base.ServerPorts = ports
			break
		}

		// Keygen options
		keyType, err := ui.AskSelect("SSH key type", []string{"ed25519", "rsa", "ecdsa"}, config.DefaultKeygenType)
		if err != nil {
			return base, err
		}
		bits, err := ui.AskInt("SSH key size (bits)", config.DefaultKeygenBits, 256, 8192)
		if err != nil {
			return base, err
		}
		force, err := ui.AskConfirm("Overwrite existing key files (-force)?", true)
		if err != nil {
			return base, err
		}
		base.Keygen = keygen.Options{Type: keyType, Bits: bits, Force: force, OutputFile: config.DefaultKeyBaseName}
	}

	if base.Target == TargetAll || base.Target == TargetPanel {
		for {
			port, err := ui.AskPort("Panel HTTP port", config.DefaultPanelPort)
			if err != nil {
				return base, err
			}
			if !network.IsPortAvailable(port) {
				ui.Warning(fmt.Sprintf("Port %d is reserved/in use. Pick another port.", port))
				continue
			}
			base.PanelPort = port
			break
		}
		user, err := ui.AskString("Panel username", config.DefaultPanelUsername, true)
		if err != nil {
			return base, err
		}
		pass, err := ui.AskPassword("Panel password", true)
		if err != nil {
			return base, err
		}
		base.PanelUsername = strings.TrimSpace(user)
		base.PanelPassword = pass
	}

	services, err := ui.AskConfirm("Register systemd/sc services after install?", true)
	if err != nil {
		return base, err
	}
	base.InstallServices = services
	return base, nil
}
