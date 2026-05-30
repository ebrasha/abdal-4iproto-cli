/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : config-constants.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Central configuration constants for repositories,
 *                installation paths, file names, default ports,
 *                service identifiers, and embedded JSON templates.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package config

// Application metadata.
const (
	AppName        = "Abdal 4iProto Cli"
	AppCommandName = "abdal-4iproto-cli"
	AppVersion     = "3.0"
	ProgrammerName = "Ebrahim Shafiei (EbraSha)"
	ProgrammerMail = "Prof.Shafiei@Gmail.com"
	ProgrammerTG   = "https://t.me/ProfShafiei"
	ProgrammerX    = "https://x.com/ProfShafiei"
	ProgrammerGH   = "https://github.com/ebrasha"
	ProgrammerLI   = "https://www.linkedin.com/in/profshafiei/"
)

// GitHub repositories for the components managed by this installer.
const (
	SSHKeygenRepoURL = "https://github.com/ebrasha/abdal-4iproto-server-ssh-keygen"
	PanelRepoURL     = "https://github.com/ebrasha/abdal-4iproto-panel"
	ServerRepoURL    = "https://github.com/ebrasha/abdal-4iproto-server"
)

// GitHub REST endpoints exposing the latest release as JSON.
const (
	SSHKeygenLatestReleaseAPI = "https://api.github.com/repos/ebrasha/abdal-4iproto-server-ssh-keygen/releases/latest"
	PanelLatestReleaseAPI     = "https://api.github.com/repos/ebrasha/abdal-4iproto-panel/releases/latest"
	ServerLatestReleaseAPI    = "https://api.github.com/repos/ebrasha/abdal-4iproto-server/releases/latest"
	CliLatestReleaseAPI       = "https://api.github.com/repos/ebrasha/abdal-4iproto-cli/releases/latest"
	CliRepoURL                = "https://github.com/ebrasha/abdal-4iproto-cli"
)

// Installation directories for each supported operating system.
// The Windows path is expanded from %LOCALAPPDATA% at runtime.
const (
	WindowsInstallDirRelative = "abdal-4iproto-server"
	LinuxInstallDir           = "/usr/local/abdal-4iproto-server"
)

// Final canonical binary file names after download (Linux).
const (
	LinuxServerBinaryName   = "abdal_4iproto_server_linux"
	LinuxPanelBinaryName    = "abdal_4iproto_panel_linux"
	LinuxKeygenBinaryName   = "abdal_4iproto_server_ssh_keygen_linux"
)

// Final canonical binary file names after download (Windows).
// Per requirement, the canonical file name keeps the "-linux" tail string
// while the executable adds the .exe extension at install time.
const (
	WindowsServerBinaryName = "abdal-4iproto-server-linux"
	WindowsPanelBinaryName  = "abdal-4iproto-panel-linux"
	WindowsKeygenBinaryName = "abdal-4iproto-server-ssh-keygen-linux"
)

// Asset naming patterns used to detect the correct file in a GitHub release.
// {arch} is replaced by runtime.GOARCH (amd64, 386, arm64, ...).
const (
	WindowsServerAssetPattern = "abdal-4iproto-server-windows-{arch}.exe"
	WindowsPanelAssetPattern  = "abdal-4iproto-panel-windows-{arch}.exe"
	WindowsKeygenAssetPattern = "abdal-4iproto-server-ssh-keygen-windows-{arch}.exe"

	LinuxServerAssetPattern = "abdal_4iproto_server_linux_{arch}"
	LinuxPanelAssetPattern  = "abdal_4iproto_panel_linux_{arch}"
	LinuxKeygenAssetPattern = "abdal_4iproto_server_ssh_keygen_linux_{arch}"

	// Self-update asset patterns for the CLI itself. Windows uses
	// hyphenated names with the .exe suffix, while Linux uses
	// underscore-separated names with no extension.
	WindowsCliAssetPattern = "abdal-4iproto-cli-windows-{arch}.exe"
	LinuxCliAssetPattern   = "abdal_4iproto_cli_linux_{arch}"
)

// Canonical CLI binary names (without the OS extension). Used by the
// self-update flow to derive the on-disk file name after download.
const (
	WindowsCliBinaryName = "abdal-4iproto-cli-windows"
	LinuxCliBinaryName   = "abdal_4iproto_cli_linux"
)

// Service identifiers (systemd unit name on Linux, Windows service name).
const (
	LinuxServerServiceName   = "abdal-4iproto-server"
	LinuxPanelServiceName    = "abdal-4iproto-panel"
	WindowsServerServiceName = "Abdal4iProtoServer"
	WindowsPanelServiceName  = "Abdal4iProtoPanel"
)

// Configuration file names placed next to the server binary.
const (
	ServerConfigFileName  = "server_config.json"
	BlockedIPsFileName    = "blocked_ips.json"
	UsersFileName         = "users.json"
	PanelConfigFileName   = "abdal-4iproto-panel.json"
	DefaultKeyBaseName    = "id_ed25519"
	DefaultKeyPublicName  = "id_ed25519.pub"
	EnvFileLinuxPath      = "/etc/default/abdal-4iproto-server"
	SystemdUnitDir        = "/etc/systemd/system"
	ServerSystemdUnitFile = "abdal-4iproto-server.service"
	PanelSystemdUnitFile  = "abdal-4iproto-panel.service"
)

// Default ports and ranges.
const (
	DefaultPanelPort     = 52202
	ServerSuggestedCount = 4
	PortSuggestionMin    = 60000
	PortSuggestionMax    = 65000
)

// Default credentials (panel) and keygen flags.
const (
	DefaultPanelUsername = "ebrasha"
	DefaultPanelPassword = "ebrasha1366"

	DefaultKeygenType  = "ed25519"
	DefaultKeygenBits  = 4096
	DefaultKeygenForce = true
)

// Shell command bound to user accounts in server_config.json.
const (
	LinuxShellPath   = "/bin/bash"
	WindowsShellPath = "cmd.exe"
)

// Default server protocol identifier embedded into server_config.json.
const (
	DefaultServerBanner = "SSH-2.0-Abdal-4iProto-Server"
)

// User role identifiers used by the user management commands.
const (
	UserRoleAdmin = "admin"
	UserRoleUser  = "user"
)

// Limits used to validate user inputs when adding/modifying accounts.
const (
	MinSessions      = 1
	MaxSessions      = 10000
	MinSpeedKbps     = 0
	MaxSpeedKbps     = 1024 * 1024
	MinPort          = 1
	MaxPort          = 65535
	CountdownSeconds = 5
)

// HTTP timeouts used by GitHub API and downloader.
const (
	HTTPTimeoutSeconds       = 60
	DownloadTimeoutSeconds   = 1800
	GithubAcceptHeaderJSON   = "application/vnd.github+json"
	GithubAcceptHeaderBinary = "application/octet-stream"
	GithubAPIVersion         = "2022-11-28"
)

// UserAgentHeader returns a freshly composed User-Agent so it always stays
// in sync with AppVersion without requiring a constant update.
var UserAgentHeader = "Abdal-4iProto-Cli/" + AppVersion + " (+" + ProgrammerGH + ")"

// Default panel configuration values.
const (
	DefaultPanelLogging         = true
	DefaultMaxLoginAttempts     = 5
	DefaultLoginAttemptWindow   = 300
	DefaultBlockDurationSeconds = 36000
	DefaultPanelTheme           = "ebrasha-dark"
)

// Default Telegram bot integration values for the panel. The bot is
// disabled out of the box; the operator turns it on after pasting a
// token and listing the admin Telegram user IDs.
const (
	DefaultTelegramBotEnabled = false
	DefaultTelegramBotToken   = ""
)
