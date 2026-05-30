/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : generator.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Generates server-side JSON configuration files with
 *                randomized passwords and OS-aware shell selection.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package filesgen

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"

	"abdal-4iproto-cli/core/config"
)

// BlockedIPsConfig mirrors blocked_ips.json.
type BlockedIPsConfig struct {
	Blocked []string `json:"blocked"`
}

// ServerConfig mirrors server_config.json.
type ServerConfig struct {
	Ports           []int  `json:"ports"`
	Shell           string `json:"shell"`
	MaxAuthAttempts int    `json:"max_auth_attempts"`
	ServerVersion   string `json:"server_version"`
	PrivateKeyFile  string `json:"private_key_file"`
	PublicKeyFile   string `json:"public_key_file"`
}

// UserAccount mirrors a single entry in users.json.
type UserAccount struct {
	Username        string   `json:"username"`
	Password        string   `json:"password"`
	Role            string   `json:"role"`
	BlockedDomains  []string `json:"blocked_domains"`
	BlockedIPs      []string `json:"blocked_ips"`
	Log             string   `json:"log"`
	MaxSessions     int      `json:"max_sessions"`
	SessionTTL      int      `json:"session_ttl_seconds"`
	MaxSpeedKbps    int      `json:"max_speed_kbps"`
	MaxTotalMB      int      `json:"max_total_mb"`
}

// TelegramBotConfig holds the Telegram bot integration block of the
// panel configuration file. The bot is disabled out of the box; the
// operator activates it after pasting a token and adding admin IDs.
type TelegramBotConfig struct {
	Enabled bool    `json:"enabled"`
	Token   string  `json:"token"`
	Admins  []int64 `json:"admins"`
}

// PanelConfig mirrors abdal-4iproto-panel.json.
type PanelConfig struct {
	Port               int               `json:"port"`
	Username           string            `json:"username"`
	Password           string            `json:"password"`
	Logging            bool              `json:"logging"`
	BlockedIPs         []string          `json:"blocked_ips"`
	MaxLoginAttempts   int               `json:"max_login_attempts"`
	LoginAttemptWindow int               `json:"login_attempt_window"`
	BlockDuration      int               `json:"block_duration"`
	Theme              string            `json:"theme"`
	TelegramBot        TelegramBotConfig `json:"telegram_bot"`
}

// KeyFileNames holds the private/public key filenames written by keygen.
type KeyFileNames struct {
	Private string
	Public  string
}

// WriteBlockedIPs creates blocked_ips.json with the default sample list.
func WriteBlockedIPs(installDir string) error {
	cfg := BlockedIPsConfig{
		Blocked: []string{"192.168.1.12", "10.0.0.7"},
	}
	return writeJSON(filepath.Join(installDir, config.BlockedIPsFileName), cfg)
}

// WriteServerConfig creates server_config.json using the provided ports and
// key filenames.
func WriteServerConfig(installDir string, ports []int, keys KeyFileNames) error {
	shell := config.LinuxShellPath
	if runtime.GOOS == "windows" {
		shell = config.WindowsShellPath
	}
	cfg := ServerConfig{
		Ports:           ports,
		Shell:           shell,
		MaxAuthAttempts: 30,
		ServerVersion:   config.DefaultServerBanner,
		PrivateKeyFile:  keys.Private,
		PublicKeyFile:   keys.Public,
	}
	return writeJSON(filepath.Join(installDir, config.ServerConfigFileName), cfg)
}

// WriteUsers creates users.json with randomized passwords for every account.
func WriteUsers(installDir string) error {
	users := []UserAccount{
		{
			Username: "ebrasha", Role: config.UserRoleAdmin,
			BlockedDomains: []string{}, BlockedIPs: []string{},
			Log: "no", MaxSessions: 1, SessionTTL: 120,
			MaxSpeedKbps: 10240, MaxTotalMB: 0,
		},
		{
			Username: "user1", Role: config.UserRoleUser,
			BlockedDomains: []string{
				"facebook.com", "*.facebook.com", "twitter.com", "*.twitter.com",
				"instagram.com", "*.instagram.com",
			},
			BlockedIPs: []string{"192.168.1.100", "10.0.0.*", "172.16.*.*"},
			Log: "yes", MaxSessions: 2, SessionTTL: 120,
			MaxSpeedKbps: 512, MaxTotalMB: 10240,
		},
		{
			Username: "user2", Role: config.UserRoleUser,
			BlockedDomains: []string{
				"youtube.com", "*.youtube.com", "netflix.com", "*.netflix.com",
			},
			BlockedIPs: []string{"192.168.10.1", "10.10.10.10"},
			Log: "yes", MaxSessions: 5, SessionTTL: 120,
			MaxSpeedKbps: 512, MaxTotalMB: 5120,
		},
	}

	for i := range users {
		pwd, err := RandomPassword(16)
		if err != nil {
			return err
		}
		users[i].Password = pwd
	}
	return writeJSON(filepath.Join(installDir, config.UsersFileName), users)
}

// WritePanelConfig creates abdal-4iproto-panel.json. The Telegram bot
// section is always written with safe defaults (disabled, no token,
// empty admins) so operators get a discoverable block to fill in later.
func WritePanelConfig(installDir string, port int, username, password string) error {
	cfg := PanelConfig{
		Port:               port,
		Username:           username,
		Password:           password,
		Logging:            config.DefaultPanelLogging,
		BlockedIPs:         []string{},
		MaxLoginAttempts:   config.DefaultMaxLoginAttempts,
		LoginAttemptWindow: config.DefaultLoginAttemptWindow,
		BlockDuration:      config.DefaultBlockDurationSeconds,
		Theme:              config.DefaultPanelTheme,
		TelegramBot: TelegramBotConfig{
			Enabled: config.DefaultTelegramBotEnabled,
			Token:   config.DefaultTelegramBotToken,
			Admins:  []int64{},
		},
	}
	return writeJSON(filepath.Join(installDir, config.PanelConfigFileName), cfg)
}

// writeJSON marshals v with indentation and writes it atomically.
func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal %s: %w", path, err)
	}
	data = append(data, '\n')
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write temp %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename %s: %w", path, err)
	}
	return nil
}

// RandomPassword generates a cryptographically secure alphanumeric password.
func RandomPassword(length int) (string, error) {
	const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if length <= 0 {
		length = 16
	}
	out := make([]byte, length)
	max := big.NewInt(int64(len(alphabet)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		out[i] = alphabet[n.Int64()]
	}
	return string(out), nil
}
