/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : manager.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Reads and updates server_config.json and the panel JSON
 *                configuration, restarting the appropriate services.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package configmgr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/filesgen"
	"abdal-4iproto-cli/core/network"
	"abdal-4iproto-cli/core/paths"
	"abdal-4iproto-cli/core/service"
	"abdal-4iproto-cli/core/ui"
)

// GetServerConfig returns the parsed server_config.json without changing it.
func GetServerConfig() (*filesgen.ServerConfig, error) {
	installDir, err := paths.InstallDir()
	if err != nil {
		return nil, err
	}
	return loadServerConfig(filepath.Join(installDir, config.ServerConfigFileName))
}

// GetPanelConfig returns the parsed abdal-4iproto-panel.json without changing it.
func GetPanelConfig() (*filesgen.PanelConfig, error) {
	installDir, err := paths.InstallDir()
	if err != nil {
		return nil, err
	}
	return loadPanelConfig(filepath.Join(installDir, config.PanelConfigFileName))
}

// UpdateServerPorts replaces the listener ports in server_config.json.
func UpdateServerPorts(ports []int) error {
	if err := network.ValidatePorts(ports); err != nil {
		return err
	}
	installDir, err := paths.InstallDir()
	if err != nil {
		return err
	}
	path := filepath.Join(installDir, config.ServerConfigFileName)

	cfg, err := loadServerConfig(path)
	if err != nil {
		return err
	}
	cfg.Ports = ports
	if err := saveServerConfig(path, cfg); err != nil {
		return err
	}
	ui.Success("Server ports updated: " + network.FormatPortList(ports))
	return service.Restart(service.ComponentServer)
}

// UpdateServerField updates a single scalar field on server_config.json.
func UpdateServerField(field string, value any) error {
	installDir, err := paths.InstallDir()
	if err != nil {
		return err
	}
	path := filepath.Join(installDir, config.ServerConfigFileName)
	cfg, err := loadServerConfig(path)
	if err != nil {
		return err
	}

	switch field {
	case "shell":
		cfg.Shell = fmt.Sprint(value)
	case "max_auth_attempts":
		if v, ok := value.(int); ok {
			cfg.MaxAuthAttempts = v
		}
	case "server_version":
		cfg.ServerVersion = fmt.Sprint(value)
	default:
		return fmt.Errorf("unsupported server field: %s", field)
	}

	if err := saveServerConfig(path, cfg); err != nil {
		return err
	}
	ui.Success(fmt.Sprintf("Server field '%s' updated", field))
	return service.Restart(service.ComponentServer)
}

// UpdatePanelConfig merges non-zero fields into abdal-4iproto-panel.json.
func UpdatePanelConfig(port int, username, password string, logging *bool) error {
	if port > 0 {
		if !network.IsPortAvailable(port) {
			return &network.PortCheckError{Port: port}
		}
	}
	installDir, err := paths.InstallDir()
	if err != nil {
		return err
	}
	path := filepath.Join(installDir, config.PanelConfigFileName)

	cfg, err := loadPanelConfig(path)
	if err != nil {
		return err
	}
	if port > 0 {
		cfg.Port = port
	}
	if username != "" {
		cfg.Username = username
	}
	if password != "" {
		cfg.Password = password
	}
	if logging != nil {
		cfg.Logging = *logging
	}

	if err := savePanelConfig(path, cfg); err != nil {
		return err
	}
	ui.Success("Panel configuration updated")
	// Requirement: restart the server service after panel JSON changes.
	return service.Restart(service.ComponentServer)
}

func loadServerConfig(path string) (*filesgen.ServerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg filesgen.ServerConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveServerConfig(path string, cfg *filesgen.ServerConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func loadPanelConfig(path string) (*filesgen.PanelConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg filesgen.PanelConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func savePanelConfig(path string, cfg *filesgen.PanelConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
