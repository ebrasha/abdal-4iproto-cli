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
	"abdal-4iproto-cli/core/selfupdate"
	"abdal-4iproto-cli/core/service"
	"abdal-4iproto-cli/core/ui"
	"abdal-4iproto-cli/core/uninstaller"
	"abdal-4iproto-cli/core/updatecheck"
	"abdal-4iproto-cli/core/users"
)

// Run launches the interactive main menu (default when no flags are passed).
func Run() error {
	ui.ClearAndBanner()
	updatecheck.Notify()

	for {
		choice, err := ui.AskSelect("Main Menu", []string{
			"Install Abdal 4iProto",
			"Uninstall Abdal 4iProto",
			"Manage Services",
			"Manage Users",
			"Server Configuration",
			"Panel Configuration",
			"Manage CLI Command (abdal-4iproto-cli)",
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

		// Wipe the screen between transitions so each sub-flow starts on
		// a clean canvas; the prior banner is re-printed for context.
		ui.ClearAndBanner()

		var actionErr error
		switch choice {
		case "Install Abdal 4iProto":
			actionErr = handleInstall()
			reportIfNotBack(actionErr, "Install Failed")
		case "Uninstall Abdal 4iProto":
			actionErr = handleUninstall()
			reportIfNotBack(actionErr, "Uninstall Failed")
		case "Manage Services":
			actionErr = handleServices()
			reportIfNotBack(actionErr, "Service Error")
		case "Manage Users":
			actionErr = handleManageUsers()
			reportIfNotBack(actionErr, "User Management Failed")
		case "Server Configuration":
			actionErr = handleServerConfigMenu()
			reportIfNotBack(actionErr, "Server Config Failed")
		case "Panel Configuration":
			actionErr = handlePanelConfigMenu()
			reportIfNotBack(actionErr, "Panel Config Failed")
		case "Manage CLI Command (abdal-4iproto-cli)":
			actionErr = handleManageCli()
			reportIfNotBack(actionErr, "CLI Management Failed")
		case "Diagnostics":
			dir, _ := paths.InstallDir()
			_ = service.Diagnostics(dir)
		case "Help":
			printInteractiveHelp()
		case "Exit":
			ui.Success("Goodbye! If you liked this tool, don't forget to star us on GitHub.")
			ui.Success("Built by Abdal Security Group, Led by Ebrahim Shafiei (EbraSha).")
			return nil
		}

		// Skip the post-action pause when the user explicitly asked to
		// go back: they have already made the decision to leave.
		if !ui.IsBack(actionErr) {
			ui.PressEnter()
		}
		ui.ClearAndBanner()
	}
}

// reportIfNotBack renders an error box only for genuine failures,
// silently swallowing the ErrUserBack sentinel.
func reportIfNotBack(err error, title string) {
	if err == nil || ui.IsBack(err) {
		return
	}
	ui.ErrorBox(title, err.Error())
}

func handleInstall() error {
	// Ask the installation scope first so the existence check below can
	// match the operator's intent precisely. Picking "Back" returns
	// silently to the main menu.
	opts := installer.DefaultOptions()
	opts, err := installer.PromptTarget(opts)
	if err != nil {
		return err
	}

	// Only warn when the artefacts that belong to the chosen scope are
	// already present. A leftover panel must not block a fresh server
	// install (and vice versa).
	report, err := installer.DetectExisting()
	if err != nil {
		return err
	}
	if report.MatchesTarget(opts.Target) {
		ui.WarningBox("Existing Installation Detected", buildScopeWarning(opts.Target, report.InstallDir))
		confirm, err := ui.AskConfirm("Proceed with a fresh install (wipe + reinstall)?", false)
		if err != nil {
			return err
		}
		if !confirm {
			return ui.ErrUserBack
		}
		opts.Force = true
	}

	opts, err = installer.PromptInstallDetails(opts)
	if err != nil {
		return err
	}
	return installer.Run(opts)
}

// buildScopeWarning returns a context-aware message that describes
// exactly which files and services the fresh install will touch so the
// operator can never be surprised about side-effects.
func buildScopeWarning(target installer.Target, dir string) string {
	switch target {
	case installer.TargetServer:
		return fmt.Sprintf(
			"Abdal 4iProto Server is already installed in:\n%s\n\nA fresh install will stop and remove ONLY the server service and its files (server binary, keygen binary, SSH keys, server_config.json, users.json, blocked_ips.json). The panel and its data will be preserved.",
			dir,
		)
	case installer.TargetPanel:
		return fmt.Sprintf(
			"Abdal 4iProto Panel is already installed in:\n%s\n\nA fresh install will stop and remove ONLY the panel service and its files (panel binary, abdal-4iproto-panel.json). The server and its data will be preserved.",
			dir,
		)
	default:
		return fmt.Sprintf(
			"Abdal 4iProto full stack (Server + Panel + KeyGen) is already installed in:\n%s\n\nA fresh install will stop every service and DELETE every file under this directory before re-downloading.",
			dir,
		)
	}
}

func handleUninstall() error {
	scope, err := ui.AskSelect("Uninstall scope", []string{
		"Full stack",
		"Server only",
		"Panel only",
		"Back",
	}, "Full stack")
	if err != nil {
		return err
	}
	if scope == "Back" {
		return ui.ErrUserBack
	}

	var target uninstaller.Target
	var removeFiles bool

	switch scope {
	case "Server only":
		target = uninstaller.TargetServer
		// Partial uninstall: never delete the shared install directory.
		removeFiles = false
		ui.Info("Only the server service will be removed. The installation folder is kept for the panel and configs.")
	case "Panel only":
		target = uninstaller.TargetPanel
		removeFiles = false
		ui.Info("Only the panel service will be removed. The installation folder is kept for the server and configs.")
	default:
		target = uninstaller.TargetAll
		installDir, _ := paths.InstallDir()
		removeFiles, err = ui.AskConfirm(
			fmt.Sprintf("Delete the installation directory and all files?\n%s", installDir),
			true,
		)
		if err != nil {
			return err
		}
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
		"Back",
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
	case "Back":
		return ui.ErrUserBack
	default:
		return service.Status(service.ComponentServer)
	}
}

// handleManageCli renders the submenu that groups CLI lifecycle actions:
// a system-wide install of the current binary and an in-place self-update
// to the latest GitHub release. Picking "Back" returns silently to the
// main menu without the post-action pause.
func handleManageCli() error {
	for {
		ui.ClearAndBanner()
		choice, err := ui.AskSelect("Manage CLI Command (abdal-4iproto-cli)", []string{
			"Install CLI Command (abdal-4iproto-cli)",
			"Update CLI Command (abdal-4iproto-cli)",
			"Back",
		}, "Install CLI Command (abdal-4iproto-cli)")
		if err != nil {
			if err == terminal.InterruptErr {
				return ui.ErrUserBack
			}
			return err
		}
		if choice == "Back" {
			return ui.ErrUserBack
		}

		ui.ClearAndBanner()
		var subErr error
		switch choice {
		case "Install CLI Command (abdal-4iproto-cli)":
			if installErr := selfinstall.Install(); installErr != nil {
				ui.ErrorBox("Self-Install Failed", installErr.Error())
				subErr = installErr
			}
		case "Update CLI Command (abdal-4iproto-cli)":
			if updateErr := selfupdate.Update(); updateErr != nil {
				ui.ErrorBox("Self-Update Failed", updateErr.Error())
				subErr = updateErr
			}
		}
		if !ui.IsBack(subErr) {
			ui.PressEnter()
		}
	}
}

// handleManageUsers renders the user-management submenu.
func handleManageUsers() error {
	for {
		ui.ClearAndBanner()
		choice, err := ui.AskSelect("Manage Users", []string{
			"List & View Users",
			"Add User",
			"Edit User",
			"Remove User",
			"Back",
		}, "List & View Users")
		if err != nil {
			if err == terminal.InterruptErr {
				return ui.ErrUserBack
			}
			return err
		}
		if choice == "Back" {
			return ui.ErrUserBack
		}

		ui.ClearAndBanner()
		var subErr error
		switch choice {
		case "List & View Users":
			subErr = handleListUsers()
			reportIfNotBack(subErr, "List Users Failed")
		case "Add User":
			subErr = handleAddUser()
			reportIfNotBack(subErr, "Add User Failed")
		case "Edit User":
			subErr = handleEditUser()
			reportIfNotBack(subErr, "Edit User Failed")
		case "Remove User":
			subErr = handleRemoveUser()
			reportIfNotBack(subErr, "Remove User Failed")
		}
		if !ui.IsBack(subErr) {
			ui.PressEnter()
		}
	}
}

// handleEditUser lets the admin pick a user and update one or more fields.
// Each field is shown with its current value so the change is intentional.
func handleEditUser() error {
	names, err := users.ListUsernames()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		ui.WarningBox("No Users", "users.json is empty. Add a user first.")
		return nil
	}

	options := append([]string{}, names...)
	options = append(options, "Cancel")

	choice, err := ui.AskSelect("Select a user to edit", options, names[0])
	if err != nil {
		if err == terminal.InterruptErr {
			return ui.ErrUserBack
		}
		return err
	}
	if choice == "Cancel" {
		return ui.ErrUserBack
	}

	account, err := users.GetUser(choice)
	if err != nil {
		return err
	}
	renderUserDetails(account)

	in := users.UpdateInput{Username: account.Username}

	for {
		field, err := ui.AskSelect("Field to edit", []string{
			"Password",
			"Role",
			"Max Sessions",
			"Max Speed (Kbps)",
			"Max Total (MB)",
			"Log",
			"Blocked Domains",
			"Blocked IPs",
			"Save & Exit",
			"Cancel",
		}, "Password")
		if err != nil {
			if err == terminal.InterruptErr {
				return nil
			}
			return err
		}

		switch field {
		case "Password":
			pwd, err := ui.AskPassword("New password", true)
			if err != nil {
				return err
			}
			in.NewPassword = &pwd
		case "Role":
			role, err := ui.AskSelect("New role", []string{config.UserRoleUser, config.UserRoleAdmin}, account.Role)
			if err != nil {
				return err
			}
			if role == config.UserRoleAdmin {
				ui.WarningBox("Security Notice", "Admin accounts are highly privileged. Proceed only if you understand the risk.")
				ok, err := ui.AskConfirm("Promote this account to admin?", false)
				if err != nil || !ok {
					ui.Info("Role change skipped.")
					continue
				}
			}
			in.NewRole = &role
		case "Max Sessions":
			v, err := ui.AskInt("Max concurrent sessions (1-10000)", account.MaxSessions, config.MinSessions, config.MaxSessions)
			if err != nil {
				return err
			}
			in.NewMaxSessions = &v
		case "Max Speed (Kbps)":
			v, err := ui.AskInt("Max speed (Kbps)", account.MaxSpeedKbps, config.MinSpeedKbps, config.MaxSpeedKbps)
			if err != nil {
				return err
			}
			in.NewMaxSpeed = &v
		case "Max Total (MB)":
			v, err := ui.AskInt("Max total transfer (MB, 0=unlimited)", account.MaxTotalMB, 0, 1<<30)
			if err != nil {
				return err
			}
			in.NewMaxTotalMB = &v
		case "Log":
			log, err := ui.AskSelect("Log activity?", []string{"yes", "no"}, account.Log)
			if err != nil {
				return err
			}
			in.NewLog = &log
		case "Blocked Domains":
			raw, err := ui.AskString("Blocked domains (comma-separated)", strings.Join(account.BlockedDomains, ","), false)
			if err != nil {
				return err
			}
			list := splitCSV(raw)
			in.NewBlockedDom = &list
		case "Blocked IPs":
			raw, err := ui.AskString("Blocked IPs (comma-separated)", strings.Join(account.BlockedIPs, ","), false)
			if err != nil {
				return err
			}
			list := splitCSV(raw)
			in.NewBlockedIPs = &list
		case "Save & Exit":
			if !updateInputHasChanges(in) {
				ui.Info("No changes to save.")
				return nil
			}
			ui.WarningBox("Confirm Update", fmt.Sprintf("Save changes for user '%s' and restart the server service?", account.Username))
			ok, err := ui.AskConfirm("Apply changes?", true)
			if err != nil || !ok {
				return ui.ErrUserBack
			}
			return users.UpdateUser(in)
		case "Cancel":
			return ui.ErrUserBack
		}
	}
}

// updateInputHasChanges reports whether at least one editable field has
// been touched by the admin during the edit session.
func updateInputHasChanges(in users.UpdateInput) bool {
	return in.NewPassword != nil ||
		in.NewRole != nil ||
		in.NewMaxSessions != nil ||
		in.NewMaxSpeed != nil ||
		in.NewMaxTotalMB != nil ||
		in.NewLog != nil ||
		in.NewBlockedDom != nil ||
		in.NewBlockedIPs != nil
}

// splitCSV converts a comma-separated string into a clean slice of values.
func splitCSV(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
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
				return ui.ErrUserBack
			}
			return err
		}
		if choice == "Back" {
			return ui.ErrUserBack
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
			return ui.ErrUserBack
		}
		return err
	}
	if choice == "Cancel" {
		return ui.ErrUserBack
	}

	// Show full record so the admin can verify before deletion.
	if account, lookupErr := users.GetUser(choice); lookupErr == nil {
		renderUserDetails(account)
	}

	ui.WarningBox("Confirm Deletion", fmt.Sprintf("You are about to permanently remove user '%s'.\nThe Abdal 4iProto Server service will be restarted afterwards.", choice))
	ok, err := ui.AskConfirm(fmt.Sprintf("Are you sure you want to remove '%s'?", choice), false)
	if err != nil {
		if err == terminal.InterruptErr {
			return ui.ErrUserBack
		}
		return err
	}
	if !ok {
		return ui.ErrUserBack
	}

	// users.RemoveUser updates users.json and restarts the server service.
	return users.RemoveUser(choice)
}

// handleServerConfigMenu groups the read-only view and the editor for
// server_config.json under a single submenu.
func handleServerConfigMenu() error {
	for {
		ui.ClearAndBanner()
		choice, err := ui.AskSelect("Server Configuration", []string{
			"View Configuration",
			"Edit Configuration",
			"Back",
		}, "View Configuration")
		if err != nil {
			if err == terminal.InterruptErr {
				return ui.ErrUserBack
			}
			return err
		}
		if choice == "Back" {
			return ui.ErrUserBack
		}

		ui.ClearAndBanner()
		var subErr error
		switch choice {
		case "View Configuration":
			subErr = handleViewServerConfig()
			reportIfNotBack(subErr, "View Server Config Failed")
		case "Edit Configuration":
			subErr = handleEditServerConfig()
			reportIfNotBack(subErr, "Edit Server Config Failed")
		}
		if !ui.IsBack(subErr) {
			ui.PressEnter()
		}
	}
}

// handlePanelConfigMenu groups the read-only view and the editor for
// abdal-4iproto-panel.json under a single submenu.
func handlePanelConfigMenu() error {
	for {
		ui.ClearAndBanner()
		choice, err := ui.AskSelect("Panel Configuration", []string{
			"View Configuration",
			"Edit Configuration",
			"Back",
		}, "View Configuration")
		if err != nil {
			if err == terminal.InterruptErr {
				return ui.ErrUserBack
			}
			return err
		}
		if choice == "Back" {
			return ui.ErrUserBack
		}

		ui.ClearAndBanner()
		var subErr error
		switch choice {
		case "View Configuration":
			subErr = handleViewPanelConfig()
			reportIfNotBack(subErr, "View Panel Config Failed")
		case "Edit Configuration":
			subErr = handleEditPanelConfig()
			reportIfNotBack(subErr, "Edit Panel Config Failed")
		}
		if !ui.IsBack(subErr) {
			ui.PressEnter()
		}
	}
}

// handleViewServerConfig prints server_config.json in a colored key/value
// box so the admin can inspect every field at a glance.
func handleViewServerConfig() error {
	cfg, err := configmgr.GetServerConfig()
	if err != nil {
		return err
	}
	portsStr := network.FormatPortList(cfg.Ports)
	if portsStr == "" {
		portsStr = "—"
	}
	ui.KeyValueBox("Server Configuration", [][2]string{
		{"Listener Ports", portsStr},
		{"Shell", cfg.Shell},
		{"Max Auth Attempts", fmt.Sprintf("%d", cfg.MaxAuthAttempts)},
		{"Server Banner", cfg.ServerVersion},
		{"Private Key File", cfg.PrivateKeyFile},
		{"Public Key File", cfg.PublicKeyFile},
	})
	return nil
}

// handleViewPanelConfig prints abdal-4iproto-panel.json in a colored
// key/value box with the password masked for safety.
func handleViewPanelConfig() error {
	cfg, err := configmgr.GetPanelConfig()
	if err != nil {
		return err
	}
	ui.KeyValueBox("Panel Configuration", [][2]string{
		{"Port", fmt.Sprintf("%d", cfg.Port)},
		{"Username", cfg.Username},
		{"Password", maskPassword(cfg.Password)},
		{"Logging", boolLabel(cfg.Logging)},
		{"Max Login Attempts", fmt.Sprintf("%d", cfg.MaxLoginAttempts)},
		{"Login Attempt Window (sec)", fmt.Sprintf("%d", cfg.LoginAttemptWindow)},
		{"Block Duration (sec)", fmt.Sprintf("%d", cfg.BlockDuration)},
		{"Theme", cfg.Theme},
		{"Blocked IPs", joinOrDash(cfg.BlockedIPs)},
	})

	ui.KeyValueBox("Telegram Bot", [][2]string{
		{"Status", boolLabel(cfg.TelegramBot.Enabled)},
		{"Token", maskToken(cfg.TelegramBot.Token)},
		{"Admins", joinAdmins(cfg.TelegramBot.Admins)},
	})

	showPass, err := ui.AskConfirm("Reveal panel password in clear text?", false)
	if err == nil && showPass {
		ui.Box("Panel Password (clear text)", cfg.Password)
	}
	return nil
}

// maskToken hides the body of the Telegram bot token while keeping a few
// characters at the start and end so the admin can recognise it.
func maskToken(t string) string {
	t = strings.TrimSpace(t)
	if t == "" {
		return "— (not configured)"
	}
	if len(t) <= 8 {
		return strings.Repeat("•", len(t))
	}
	return t[:4] + strings.Repeat("•", len(t)-8) + t[len(t)-4:]
}

// joinAdmins renders the list of admin Telegram IDs, returning an em dash
// when the slice is empty so the value column never looks blank.
func joinAdmins(ids []int64) string {
	if len(ids) == 0 {
		return "—"
	}
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("%d", id))
	}
	return strings.Join(parts, ", ")
}

// boolLabel returns a human-readable label for a boolean configuration.
func boolLabel(b bool) string {
	if b {
		return "enabled"
	}
	return "disabled"
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
	ui.SectionHeader(config.AppName + " v" + config.AppVersion + " – Help & Reference")
	ui.KeyValueBox("Version", [][2]string{
		{"Application", config.AppName},
		{"Version", config.AppVersion},
		{"Command", config.AppCommandName},
		{"Repository", config.CliRepoURL},
	})
	ui.Box("Programmer", config.ProgrammerName+"\n"+config.ProgrammerMail+"\n"+config.ProgrammerTG)
	ui.Box("Interactive Menu", strings.TrimSpace(`
Main Menu
  Install Abdal 4iProto
  Uninstall Abdal 4iProto
  Manage Services
  Manage Users  ──► List & View Users / Add User / Edit User / Remove User
  Server Configuration ──► View Configuration / Edit Configuration
  Panel Configuration  ──► View Configuration / Edit Configuration
  Manage CLI Command (abdal-4iproto-cli) ──► Install CLI / Update CLI / Back
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
abdal-4iproto-cli self-update
abdal-4iproto-cli help
`))
	ui.Box("Service Names", strings.TrimSpace(`
Linux systemd units : abdal-4iproto-server, abdal-4iproto-panel
Windows services    : Abdal4iProtoServer, Abdal4iProtoPanel
`))
	ui.Box("Diagnostics", "Use the Diagnostics menu entry or: abdal-4iproto-cli service diagnostics")
}
