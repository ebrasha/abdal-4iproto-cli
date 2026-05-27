/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : menu.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Primary interactive menu loop for operators who run
 *                the CLI without subcommands or flags.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package interactive

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/configmgr"
	"abdal-4iproto-cli/core/filesgen"
	"abdal-4iproto-cli/core/installer"
	"abdal-4iproto-cli/core/network"
	"abdal-4iproto-cli/core/paths"
	"abdal-4iproto-cli/core/selfinstall"
	"abdal-4iproto-cli/core/service"
	"abdal-4iproto-cli/core/ui"
	"abdal-4iproto-cli/core/uninstaller"
	"abdal-4iproto-cli/core/users"
)

// Run launches the interactive main menu (default when no flags are passed).
func Run() error {
	ui.PrintBanner()

	for {
		choice, err := ui.AskSelect("Main Menu", []string{
			"Install Abdal 4iProto",
			"Uninstall Abdal 4iProto",
			"Manage Services",
			"Manage Users",
			"Edit Server Configuration",
			"Edit Panel Configuration",
			"Install CLI Command (abdal-4iproto-cli)",
			"Diagnostics",
			"Help",
			"Exit",
		}, "Install Abdal 4iProto")
		if err != nil {
			if err == terminal.InterruptErr {
				return nil
			}
			return err
		}

		switch choice {
		case "Install Abdal 4iProto":
			if err := handleInstall(); err != nil {
				ui.ErrorBox("Install Failed", err.Error())
			}
		case "Uninstall Abdal 4iProto":
			if err := handleUninstall(); err != nil {
				ui.ErrorBox("Uninstall Failed", err.Error())
			}
		case "Manage Services":
			if err := handleServices(); err != nil {
				ui.ErrorBox("Service Error", err.Error())
			}
		case "Manage Users":
			if err := handleManageUsers(); err != nil {
				ui.ErrorBox("User Management Failed", err.Error())
			}
		case "Edit Server Configuration":
			if err := handleEditServerConfig(); err != nil {
				ui.ErrorBox("Server Config Failed", err.Error())
			}
		case "Edit Panel Configuration":
			if err := handleEditPanelConfig(); err != nil {
				ui.ErrorBox("Panel Config Failed", err.Error())
			}
		case "Install CLI Command (abdal-4iproto-cli)":
			if err := selfinstall.Install(); err != nil {
				ui.ErrorBox("Self-Install Failed", err.Error())
			}
		case "Diagnostics":
			dir, _ := paths.InstallDir()
			_ = service.Diagnostics(dir)
		case "Help":
			printInteractiveHelp()
		case "Exit":
			ui.Success("Goodbye.")
			return nil
		}
		fmt.Println()
	}
}

func handleInstall() error {
	opts := installer.DefaultOptions()
	var err error
	opts, err = installer.PromptOptions(opts)
	if err != nil {
		return err
	}
	return installer.Run(opts)
}

func handleUninstall() error {
	scope, err := ui.AskSelect("Uninstall scope", []string{
		"Full stack",
		"Server only",
		"Panel only",
	}, "Full stack")
	if err != nil {
		return err
	}
	removeFiles, err := ui.AskConfirm("Delete installation directory and files?", true)
	if err != nil {
		return err
	}
	var target uninstaller.Target
	switch scope {
	case "Server only":
		target = uninstaller.TargetServer
	case "Panel only":
		target = uninstaller.TargetPanel
	default:
		target = uninstaller.TargetAll
	}
	return uninstaller.Run(target, removeFiles)
}

func handleServices() error {
	action, err := ui.AskSelect("Service action", []string{
		"Status (Server)",
		"Status (Panel)",
		"Restart Server",
		"Restart Panel",
		"Diagnostics bundle",
	}, "Status (Server)")
	if err != nil {
		return err
	}
	switch action {
	case "Status (Panel)":
		return service.Status(service.ComponentPanel)
	case "Restart Server":
		return service.Restart(service.ComponentServer)
	case "Restart Panel":
		return service.Restart(service.ComponentPanel)
	case "Diagnostics bundle":
		dir, _ := paths.InstallDir()
		return service.Diagnostics(dir)
	default:
		return service.Status(service.ComponentServer)
	}
}

// handleManageUsers renders the user-management submenu.
func handleManageUsers() error {
	for {
		choice, err := ui.AskSelect("Manage Users", []string{
			"List & View Users",
			"Add User",
			"Remove User",
			"Back",
		}, "List & View Users")
		if err != nil {
			if err == terminal.InterruptErr {
				return nil
			}
			return err
		}
		switch choice {
		case "List & View Users":
			if err := handleListUsers(); err != nil {
				ui.ErrorBox("List Users Failed", err.Error())
			}
		case "Add User":
			if err := handleAddUser(); err != nil {
				ui.ErrorBox("Add User Failed", err.Error())
			}
		case "Remove User":
			if err := handleRemoveUser(); err != nil {
				ui.ErrorBox("Remove User Failed", err.Error())
			}
		case "Back":
			return nil
		}
		fmt.Println()
	}
}

// handleListUsers shows the usernames and lets the admin pick one to view
// its full record in a colored, structured box.
func handleListUsers() error {
	names, err := users.ListUsernames()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		ui.WarningBox("No Users", "users.json is empty. Add a user first.")
		return nil
	}

	options := append([]string{}, names...)
	options = append(options, "Back")

	for {
		choice, err := ui.AskSelect("Select a user to view details", options, names[0])
		if err != nil {
			if err == terminal.InterruptErr {
				return nil
			}
			return err
		}
		if choice == "Back" {
			return nil
		}
		account, err := users.GetUser(choice)
		if err != nil {
			ui.ErrorBox("Lookup Failed", err.Error())
			continue
		}
		renderUserDetails(account)

		showPass, err := ui.AskConfirm("Reveal full password for this account?", false)
		if err == nil && showPass {
			ui.Box("Password (clear text)", account.Password)
		}
	}
}

// renderUserDetails prints a single account record using a colored key/value
// box so the admin can read every field at a glance.
func renderUserDetails(u *filesgen.UserAccount) {
	ui.KeyValueBox("User: "+u.Username, [][2]string{
		{"Username", u.Username},
		{"Password", maskPassword(u.Password)},
		{"Role", u.Role},
		{"Log", u.Log},
		{"Max Sessions", fmt.Sprintf("%d", u.MaxSessions)},
		{"Session TTL (sec)", fmt.Sprintf("%d", u.SessionTTL)},
		{"Max Speed (Kbps)", fmt.Sprintf("%d", u.MaxSpeedKbps)},
		{"Max Total (MB)", quotaLabel(u.MaxTotalMB)},
		{"Blocked Domains", joinOrDash(u.BlockedDomains)},
		{"Blocked IPs", joinOrDash(u.BlockedIPs)},
	})
}

// maskPassword hides the middle portion of the password to keep secrets out
// of casual screenshots while still showing the first and last character.
func maskPassword(p string) string {
	if len(p) == 0 {
		return "—"
	}
	if len(p) <= 2 {
		return "••"
	}
	return string(p[0]) + strings.Repeat("•", len(p)-2) + string(p[len(p)-1])
}

// joinOrDash joins a string slice with commas and returns an em dash when
// the slice is empty so the box never shows a blank value.
func joinOrDash(xs []string) string {
	if len(xs) == 0 {
		return "—"
	}
	return strings.Join(xs, ", ")
}

// quotaLabel renders the total-transfer quota with a helpful unlimited tag.
func quotaLabel(mb int) string {
	if mb <= 0 {
		return "0 (unlimited)"
	}
	return fmt.Sprintf("%d", mb)
}

func handleAddUser() error {
	username, err := ui.AskString("Username", "", true)
	if err != nil {
		return err
	}
	password, err := ui.AskPassword("Password", true)
	if err != nil {
		return err
	}
	role, err := ui.AskSelect("Role", []string{config.UserRoleAdmin, config.UserRoleUser}, config.UserRoleUser)
	if err != nil {
		return err
	}
	if role == config.UserRoleAdmin {
		ui.WarningBox("Security Notice", "Admin accounts are highly privileged. Proceed only if you understand the risk.")
		ok, err := ui.AskConfirm("Continue creating an admin account?", false)
		if err != nil || !ok {
			return fmt.Errorf("admin creation cancelled")
		}
	}
	maxSessions, err := ui.AskInt("Max concurrent sessions (1-10000)", 1, config.MinSessions, config.MaxSessions)
	if err != nil {
		return err
	}
	maxSpeed, err := ui.AskInt("Max speed (Kbps)", 512, 0, config.MaxSpeedKbps)
	if err != nil {
		return err
	}
	maxTotal, err := ui.AskInt("Max total transfer (MB, 0=unlimited)", 0, 0, 1<<30)
	if err != nil {
		return err
	}
	return users.AddUser(users.AddInput{
		Username: username, Password: password, Role: role,
		MaxSessions: maxSessions, MaxSpeedKbps: maxSpeed, MaxTotalMB: maxTotal,
	})
}

func handleRemoveUser() error {
	names, err := users.ListUsernames()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		ui.WarningBox("No Users", "users.json is empty. There is nothing to remove.")
		return nil
	}

	options := append([]string{}, names...)
	options = append(options, "Cancel")

	choice, err := ui.AskSelect("Select a user to remove", options, names[0])
	if err != nil {
		if err == terminal.InterruptErr {
			return nil
		}
		return err
	}
	if choice == "Cancel" {
		ui.Info("Removal cancelled.")
		return nil
	}

	// Show full record so the admin can verify before deletion.
	if account, lookupErr := users.GetUser(choice); lookupErr == nil {
		renderUserDetails(account)
	}

	ui.WarningBox("Confirm Deletion", fmt.Sprintf("You are about to permanently remove user '%s'.\nThe Abdal 4iProto Server service will be restarted afterwards.", choice))
	ok, err := ui.AskConfirm(fmt.Sprintf("Are you sure you want to remove '%s'?", choice), false)
	if err != nil {
		if err == terminal.InterruptErr {
			return nil
		}
		return err
	}
	if !ok {
		ui.Info("Removal cancelled.")
		return nil
	}

	// users.RemoveUser updates users.json and restarts the server service.
	return users.RemoveUser(choice)
}

func handleEditServerConfig() error {
	field, err := ui.AskSelect("Field to update", []string{"ports", "max_auth_attempts", "server_version"}, "ports")
	if err != nil {
		return err
	}
	switch field {
	case "ports":
		for {
			raw, err := ui.AskString("New ports (comma-separated)", "", true)
			if err != nil {
				return err
			}
			ports, err := network.ParsePortList(raw)
			if err != nil {
				ui.Warning(err.Error())
				continue
			}
			if err := network.ValidatePorts(ports); err != nil {
				if pce, ok := err.(*network.PortCheckError); ok {
					ui.Warning(fmt.Sprintf("Port %d is in use. Try again.", pce.Port))
					continue
				}
				return err
			}
			return configmgr.UpdateServerPorts(ports)
		}
	case "max_auth_attempts":
		v, err := ui.AskInt("max_auth_attempts", 30, 1, 1000)
		if err != nil {
			return err
		}
		return configmgr.UpdateServerField("max_auth_attempts", v)
	default:
		val, err := ui.AskString("server_version", config.DefaultServerBanner, true)
		if err != nil {
			return err
		}
		return configmgr.UpdateServerField("server_version", val)
	}
}

func handleEditPanelConfig() error {
	for {
		port, err := ui.AskPort("Panel port", config.DefaultPanelPort)
		if err != nil {
			return err
		}
		if !network.IsPortAvailable(port) {
			ui.Warning(fmt.Sprintf("Port %d is in use. Choose another.", port))
			continue
		}
		user, err := ui.AskString("Panel username", config.DefaultPanelUsername, true)
		if err != nil {
			return err
		}
		pass, err := ui.AskPassword("Panel password", true)
		if err != nil {
			return err
		}
		return configmgr.UpdatePanelConfig(port, strings.TrimSpace(user), pass, nil)
	}
}

func printInteractiveHelp() {
	ui.SectionHeader("Help & Reference")
	ui.Box("Programmer", config.ProgrammerName+"\n"+config.ProgrammerMail+"\n"+config.ProgrammerTG)
	ui.Box("Interactive Menu", strings.TrimSpace(`
Main Menu
  Install Abdal 4iProto
  Uninstall Abdal 4iProto
  Manage Services
  Manage Users  ──► List & View Users / Add User / Remove User
  Edit Server Configuration
  Edit Panel Configuration
  Install CLI Command (abdal-4iproto-cli)
  Diagnostics
  Help / Exit
`))
	ui.Box("CLI Usage (non-interactive)", strings.TrimSpace(`
abdal-4iproto-cli install [--server-only|--panel-only] [--server-ports 64235,64236] [--panel-port 52202]
abdal-4iproto-cli uninstall [--server-only|--panel-only] [--keep-files]
abdal-4iproto-cli user add --username X --password Y --role user --max-sessions 2 --max-speed-kbps 512 --max-total-mb 0
abdal-4iproto-cli user remove --username X
abdal-4iproto-cli service status|restart --component server|panel
abdal-4iproto-cli config server --ports 64235,64236
abdal-4iproto-cli config panel --port 52202 --username ebrasha --password secret
abdal-4iproto-cli self-install
abdal-4iproto-cli help
`))
	ui.Box("Service Names", strings.TrimSpace(`
Linux systemd units : abdal-4iproto-server, abdal-4iproto-panel
Windows services    : Abdal4iProtoServer, Abdal4iProtoPanel
`))
	ui.Box("Diagnostics", "Use the Diagnostics menu entry or: abdal-4iproto-cli service diagnostics")
}
