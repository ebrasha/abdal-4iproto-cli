/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : install.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Orchestrates the full or partial installation workflow
 *                including port checks, downloads, keygen, configs, and
 *                service registration.
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
	"os"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/downloader"
	"abdal-4iproto-cli/core/filesgen"
	gh "abdal-4iproto-cli/core/github"
	"abdal-4iproto-cli/core/keygen"
	"abdal-4iproto-cli/core/network"
	"abdal-4iproto-cli/core/paths"
	"abdal-4iproto-cli/core/service"
	"abdal-4iproto-cli/core/ui"
)

// Run executes the installation pipeline according to opts.
func Run(opts Options) error {
	installDir, err := paths.InstallDir()
	if err != nil {
		return fmt.Errorf("resolve install directory: %w", err)
	}

	ui.SectionHeader("Abdal 4iProto Installation")
	ui.KeyValueBox("Target", [][2]string{
		{"Install Directory", installDir},
		{"Component", targetLabel(opts.Target)},
	})

	// Suggest server ports when none were provided (interactive/CLI may set them).
	if (opts.Target == TargetAll || opts.Target == TargetServer) && len(opts.ServerPorts) == 0 {
		suggested, err := network.SuggestFreePorts(config.ServerSuggestedCount)
		if err != nil {
			return err
		}
		opts.ServerPorts = suggested
	}

	// Port validation must happen before any download (requirement).
	if err := validatePortsBeforeDownload(opts); err != nil {
		return err
	}

	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return fmt.Errorf("create install directory: %w", err)
	}

	needServer := opts.Target == TargetAll || opts.Target == TargetServer
	needPanel := opts.Target == TargetAll || opts.Target == TargetPanel

	// Always download keygen when server stack is requested (keys are required).
	if needServer {
		if err := downloadKeygen(installDir); err != nil {
			return err
		}
		if !opts.SkipKeygen {
			keys, err := keygen.Generate(installDir, opts.Keygen)
			if err != nil {
				return err
			}
			if err := writeServerFiles(installDir, opts, keys); err != nil {
				return err
			}
		}
		if err := downloadAndInstallServer(installDir, opts); err != nil {
			return err
		}
	}

	if needPanel {
		if err := downloadAndInstallPanel(installDir, opts); err != nil {
			return err
		}
	}

	ui.SuccessBox("Installation Complete", "Abdal 4iProto components are ready in:\n"+installDir)
	return nil
}

func validatePortsBeforeDownload(opts Options) error {
	if opts.Target == TargetAll || opts.Target == TargetServer {
		if err := network.ValidatePorts(opts.ServerPorts); err != nil {
			return fmt.Errorf("server ports validation failed: %w", err)
		}
		ui.Info("Server ports: " + network.FormatPortList(opts.ServerPorts))
	}
	if opts.Target == TargetAll || opts.Target == TargetPanel {
		if !network.IsPortAvailable(opts.PanelPort) {
			return &network.PortCheckError{Port: opts.PanelPort}
		}
		ui.Info(fmt.Sprintf("Panel port: %d", opts.PanelPort))
	}
	return nil
}

func downloadKeygen(installDir string) error {
	ui.Step(1, 6, "Fetching latest SSH KeyGen release metadata")
	rel, err := gh.FetchLatestRelease(config.SSHKeygenLatestReleaseAPI)
	if err != nil {
		return err
	}
	sel, err := gh.ChooseAsset(rel,
		config.WindowsKeygenAssetPattern, config.LinuxKeygenAssetPattern,
		config.WindowsKeygenBinaryName, config.LinuxKeygenBinaryName,
	)
	if err != nil {
		return err
	}
	ui.Step(2, 6, "Downloading SSH KeyGen")
	_, err = downloader.DownloadAsset(sel, installDir)
	return err
}

func downloadAndInstallServer(installDir string, opts Options) error {
	ui.Step(3, 6, "Fetching latest Server release metadata")
	rel, err := gh.FetchLatestRelease(config.ServerLatestReleaseAPI)
	if err != nil {
		return err
	}
	sel, err := gh.ChooseAsset(rel,
		config.WindowsServerAssetPattern, config.LinuxServerAssetPattern,
		config.WindowsServerBinaryName, config.LinuxServerBinaryName,
	)
	if err != nil {
		return err
	}
	ui.Step(4, 6, "Downloading Abdal 4iProto Server")
	if _, err := downloader.DownloadAsset(sel, installDir); err != nil {
		return err
	}
	if opts.InstallServices {
		ui.Step(5, 6, "Registering server service")
		if err := service.Install(installDir, service.ComponentServer); err != nil {
			return err
		}
	}
	return nil
}

func downloadAndInstallPanel(installDir string, opts Options) error {
	ui.Step(6, 6, "Fetching latest Panel release metadata")
	rel, err := gh.FetchLatestRelease(config.PanelLatestReleaseAPI)
	if err != nil {
		return err
	}
	sel, err := gh.ChooseAsset(rel,
		config.WindowsPanelAssetPattern, config.LinuxPanelAssetPattern,
		config.WindowsPanelBinaryName, config.LinuxPanelBinaryName,
	)
	if err != nil {
		return err
	}
	ui.Info("Downloading Abdal 4iProto Panel")
	if _, err := downloader.DownloadAsset(sel, installDir); err != nil {
		return err
	}
	if err := filesgen.WritePanelConfig(installDir, opts.PanelPort, opts.PanelUsername, opts.PanelPassword); err != nil {
		return err
	}
	if opts.InstallServices {
		if err := service.Install(installDir, service.ComponentPanel); err != nil {
			return err
		}
	}
	return nil
}

func writeServerFiles(installDir string, opts Options, keys filesgen.KeyFileNames) error {
	if len(opts.ServerPorts) == 0 {
		suggested, err := network.SuggestFreePorts(config.ServerSuggestedCount)
		if err != nil {
			return err
		}
		opts.ServerPorts = suggested
	}
	if err := filesgen.WriteBlockedIPs(installDir); err != nil {
		return err
	}
	if err := filesgen.WriteServerConfig(installDir, opts.ServerPorts, keys); err != nil {
		return err
	}
	return filesgen.WriteUsers(installDir)
}

func targetLabel(t Target) string {
	switch t {
	case TargetServer:
		return "Server only"
	case TargetPanel:
		return "Panel only"
	default:
		return "Full stack (Server + Panel + KeyGen)"
	}
}
