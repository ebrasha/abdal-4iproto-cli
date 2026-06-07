/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : service.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Installs and manages Abdal 4iProto server/panel services
 *                on Linux (systemd) and Windows (Service Control Manager).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/paths"
	"abdal-4iproto-cli/core/ui"
)

// Component identifies which binary is managed as a service.
type Component string

const (
	ComponentServer Component = "server"
	ComponentPanel  Component = "panel"
)

// Install registers and starts the requested component as a system service.
func Install(installDir string, component Component) error {
	switch runtime.GOOS {
	case "linux":
		return installLinux(installDir, component)
	case "windows":
		return installWindows(installDir, component)
	default:
		return fmt.Errorf("unsupported OS for service install: %s", runtime.GOOS)
	}
}

// Uninstall stops and removes the service for the given component.
func Uninstall(component Component) error {
	switch runtime.GOOS {
	case "linux":
		return uninstallLinux(component)
	case "windows":
		return uninstallWindows(component)
	default:
		return fmt.Errorf("unsupported OS for service uninstall: %s", runtime.GOOS)
	}
}

// Restart restarts the service backing the component.
func Restart(component Component) error {
	switch runtime.GOOS {
	case "linux":
		name := linuxUnitName(component)
		return runCmd("systemctl", "restart", name)
	case "windows":
		name := windowsServiceName(component)
		_ = runCmd("sc", "stop", name)
		return runCmd("sc", "start", name)
	default:
		return fmt.Errorf("unsupported OS for service restart: %s", runtime.GOOS)
	}
}

// Start launches an already-installed service for the component.
func Start(component Component) error {
	switch runtime.GOOS {
	case "linux":
		if err := runCmd("systemctl", "start", linuxUnitName(component)); err != nil {
			return err
		}
	case "windows":
		if err := runCmd("sc", "start", windowsServiceName(component)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported OS for service start: %s", runtime.GOOS)
	}
	ui.Success(fmt.Sprintf("Service '%s' started", serviceLabel(component)))
	return nil
}

// Stop halts the running service for the component without removing it,
// so it can be started again later or have its binary replaced safely.
func Stop(component Component) error {
	switch runtime.GOOS {
	case "linux":
		if err := runCmd("systemctl", "stop", linuxUnitName(component)); err != nil {
			return err
		}
	case "windows":
		if err := runCmd("sc", "stop", windowsServiceName(component)); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported OS for service stop: %s", runtime.GOOS)
	}
	ui.Success(fmt.Sprintf("Service '%s' stopped", serviceLabel(component)))
	return nil
}

// Enable configures the service to start automatically at system boot.
func Enable(component Component) error {
	switch runtime.GOOS {
	case "linux":
		if err := runCmd("systemctl", "enable", linuxUnitName(component)); err != nil {
			return err
		}
	case "windows":
		// "start= auto" makes the Windows service start automatically.
		if err := runCmd("sc", "config", windowsServiceName(component), "start=", "auto"); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported OS for service enable: %s", runtime.GOOS)
	}
	ui.Success(fmt.Sprintf("Service '%s' enabled (auto-start at boot)", serviceLabel(component)))
	return nil
}

// Disable prevents the service from starting automatically at system boot.
func Disable(component Component) error {
	switch runtime.GOOS {
	case "linux":
		if err := runCmd("systemctl", "disable", linuxUnitName(component)); err != nil {
			return err
		}
	case "windows":
		// "start= disabled" stops Windows from launching the service.
		if err := runCmd("sc", "config", windowsServiceName(component), "start=", "disabled"); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported OS for service disable: %s", runtime.GOOS)
	}
	ui.Success(fmt.Sprintf("Service '%s' disabled (no auto-start at boot)", serviceLabel(component)))
	return nil
}

// serviceLabel returns the OS-specific service identifier for messages.
func serviceLabel(component Component) string {
	if runtime.GOOS == "windows" {
		return windowsServiceName(component)
	}
	return linuxUnitName(component)
}

// Status prints a short status summary for the component.
func Status(component Component) error {
	switch runtime.GOOS {
	case "linux":
		name := linuxUnitName(component)
		ui.SectionHeader("Service Status: " + name)
		return runCmd("systemctl", "status", name, "--no-pager")
	case "windows":
		name := windowsServiceName(component)
		ui.SectionHeader("Service Status: " + name)
		return runCmd("sc", "query", name)
	default:
		return fmt.Errorf("unsupported OS for service status: %s", runtime.GOOS)
	}
}

// Diagnostics runs a lightweight troubleshooting bundle.
func Diagnostics(installDir string) error {
	ui.SectionHeader("Diagnostics")
	ui.KeyValueBox("Environment", [][2]string{
		{"OS", runtime.GOOS},
		{"Arch", runtime.GOARCH},
		{"Install Dir", installDir},
	})
	_ = Status(ComponentServer)
	_ = Status(ComponentPanel)
	return nil
}

// --- Linux (systemd) ---

func installLinux(installDir string, component Component) error {
	unitContent, unitName, err := linuxUnitContent(installDir, component)
	if err != nil {
		return err
	}
	unitPath := filepath.Join(config.SystemdUnitDir, unitName)
	if err := os.WriteFile(unitPath, []byte(unitContent), 0o644); err != nil {
		return fmt.Errorf("write unit file %s: %w", unitPath, err)
	}
	if err := runCmd("systemctl", "daemon-reload"); err != nil {
		return err
	}
	if err := runCmd("systemctl", "enable", strings.TrimSuffix(unitName, ".service")); err != nil {
		return err
	}
	if err := runCmd("systemctl", "start", strings.TrimSuffix(unitName, ".service")); err != nil {
		return err
	}
	ui.Success(fmt.Sprintf("systemd service '%s' installed and started", strings.TrimSuffix(unitName, ".service")))
	return nil
}

func uninstallLinux(component Component) error {
	name := linuxUnitName(component)
	_ = runCmd("systemctl", "stop", name)
	_ = runCmd("systemctl", "disable", name)
	unitPath := filepath.Join(config.SystemdUnitDir, name+".service")
	_ = os.Remove(unitPath)
	_ = runCmd("systemctl", "daemon-reload")
	ui.Success(fmt.Sprintf("systemd service '%s' removed", name))
	return nil
}

func linuxUnitName(component Component) string {
	switch component {
	case ComponentPanel:
		return config.LinuxPanelServiceName
	default:
		return config.LinuxServerServiceName
	}
}

func linuxUnitContent(installDir string, component Component) (string, string, error) {
	var execPath, description, syslogID, unitFile string
	switch component {
	case ComponentServer:
		execPath = paths.ServerBinaryPath(installDir)
		description = "Abdal 4iProto Server"
		syslogID = config.LinuxServerServiceName
		unitFile = config.ServerSystemdUnitFile
	case ComponentPanel:
		execPath = paths.PanelBinaryPath(installDir)
		description = "Abdal 4iProto Panel"
		syslogID = config.LinuxPanelServiceName
		unitFile = config.PanelSystemdUnitFile
	default:
		return "", "", fmt.Errorf("unknown component: %s", component)
	}

	content := fmt.Sprintf(`# -------------------------------------------------------------------
# Programmer       : Ebrahim Shafiei (EbraSha)
# Email            : Prof.Shafiei@Gmail.com
# -------------------------------------------------------------------
[Unit]
Description=%s
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=%s
EnvironmentFile=-%s
ExecStart=%s
Restart=always
RestartSec=3
LimitNOFILE=65536
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full
ProtectHome=true
ReadWritePaths=%s
SyslogIdentifier=%s

[Install]
WantedBy=multi-user.target
`, description, installDir, config.EnvFileLinuxPath, execPath, installDir, syslogID)

	return content, unitFile, nil
}

// --- Windows (sc) ---

func installWindows(installDir string, component Component) error {
	svcName, binPath, displayName, err := windowsServiceParams(installDir, component)
	if err != nil {
		return err
	}

	// Remove stale service if present.
	_ = runCmd("sc", "stop", svcName)
	_ = runCmd("sc", "delete", svcName)

	createArgs := []string{
		"create", svcName,
		"binPath=", binPath,
		"start=", "auto",
		"DisplayName=", displayName,
	}
	if err := runCmd("sc", createArgs...); err != nil {
		return fmt.Errorf("sc create failed: %w", err)
	}
	if err := runCmd("sc", "description", svcName, displayName); err != nil {
		ui.Warning(fmt.Sprintf("sc description failed: %v", err))
	}
	if err := runCmd("sc", "start", svcName); err != nil {
		return fmt.Errorf("sc start failed: %w", err)
	}
	ui.Success(fmt.Sprintf("Windows service '%s' installed and started", svcName))
	return nil
}

func uninstallWindows(component Component) error {
	name := windowsServiceName(component)
	_ = runCmd("sc", "stop", name)
	if err := runCmd("sc", "delete", name); err != nil {
		return err
	}
	ui.Success(fmt.Sprintf("Windows service '%s' removed", name))
	return nil
}

func windowsServiceName(component Component) string {
	switch component {
	case ComponentPanel:
		return config.WindowsPanelServiceName
	default:
		return config.WindowsServerServiceName
	}
}

func windowsServiceParams(installDir string, component Component) (name, binPath, display string, err error) {
	switch component {
	case ComponentServer:
		name = config.WindowsServerServiceName
		binPath = paths.ServerBinaryPath(installDir)
		display = "Abdal 4iProto Server"
	case ComponentPanel:
		name = config.WindowsPanelServiceName
		binPath = paths.PanelBinaryPath(installDir)
		display = "Abdal 4iProto Panel"
	default:
		err = fmt.Errorf("unknown component: %s", component)
	}
	// sc requires quoted binPath when spaces are present.
	if strings.Contains(binPath, " ") {
		binPath = `"` + binPath + `"`
	}
	return
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s %v: %w", name, args, err)
	}
	return nil
}
