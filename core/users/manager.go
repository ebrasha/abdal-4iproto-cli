/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : manager.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Adds and removes SSH server user accounts stored in
 *                users.json and restarts the server service afterwards.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package users

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/filesgen"
	"abdal-4iproto-cli/core/paths"
	"abdal-4iproto-cli/core/service"
	"abdal-4iproto-cli/core/ui"
)

// AddInput contains mandatory fields for creating a new account.
type AddInput struct {
	Username     string
	Password     string
	Role         string
	MaxSessions  int
	MaxSpeedKbps int
	MaxTotalMB   int
}

// AddUser appends a new account to users.json and restarts the server service.
func AddUser(in AddInput) error {
	if err := validateAddInput(in); err != nil {
		return err
	}
	installDir, err := paths.InstallDir()
	if err != nil {
		return err
	}
	path := filepath.Join(installDir, config.UsersFileName)

	accounts, err := loadUsers(path)
	if err != nil {
		return err
	}
	for _, u := range accounts {
		if strings.EqualFold(u.Username, in.Username) {
			return fmt.Errorf("user '%s' already exists", in.Username)
		}
	}

	accounts = append(accounts, filesgen.UserAccount{
		Username:       in.Username,
		Password:       in.Password,
		Role:           in.Role,
		BlockedDomains: []string{},
		BlockedIPs:     []string{},
		Log:            "yes",
		MaxSessions:    in.MaxSessions,
		SessionTTL:     120,
		MaxSpeedKbps:   in.MaxSpeedKbps,
		MaxTotalMB:     in.MaxTotalMB,
	})

	if err := saveUsers(path, accounts); err != nil {
		return err
	}
	ui.Success(fmt.Sprintf("User '%s' added", in.Username))
	return service.Restart(service.ComponentServer)
}

// RemoveUser deletes an account by username and restarts the server service.
func RemoveUser(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("username is required")
	}
	installDir, err := paths.InstallDir()
	if err != nil {
		return err
	}
	path := filepath.Join(installDir, config.UsersFileName)

	accounts, err := loadUsers(path)
	if err != nil {
		return err
	}

	var next []filesgen.UserAccount
	var removed bool
	for _, u := range accounts {
		if strings.EqualFold(u.Username, username) {
			removed = true
			continue
		}
		next = append(next, u)
	}
	if !removed {
		return fmt.Errorf("user '%s' not found", username)
	}
	if err := saveUsers(path, next); err != nil {
		return err
	}
	ui.Success(fmt.Sprintf("User '%s' removed", username))
	return service.Restart(service.ComponentServer)
}

// ListUsernames returns the list of usernames stored in users.json.
// It performs a read-only scan and never restarts any service.
func ListUsernames() ([]string, error) {
	installDir, err := paths.InstallDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(installDir, config.UsersFileName)
	accounts, err := loadUsers(path)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(accounts))
	for _, u := range accounts {
		names = append(names, u.Username)
	}
	return names, nil
}

// GetUser returns the full account record for a given username
// (case-insensitive). It performs a read-only lookup.
func GetUser(username string) (*filesgen.UserAccount, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	installDir, err := paths.InstallDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(installDir, config.UsersFileName)
	accounts, err := loadUsers(path)
	if err != nil {
		return nil, err
	}
	for i := range accounts {
		if strings.EqualFold(accounts[i].Username, username) {
			return &accounts[i], nil
		}
	}
	return nil, fmt.Errorf("user '%s' not found", username)
}

func validateAddInput(in AddInput) error {
	if strings.TrimSpace(in.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(in.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if in.Role != config.UserRoleAdmin && in.Role != config.UserRoleUser {
		return fmt.Errorf("role must be '%s' or '%s'", config.UserRoleAdmin, config.UserRoleUser)
	}
	if in.MaxSessions < config.MinSessions || in.MaxSessions > config.MaxSessions {
		return fmt.Errorf("max_sessions must be between %d and %d", config.MinSessions, config.MaxSessions)
	}
	return nil
}

func loadUsers(path string) ([]filesgen.UserAccount, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read users file: %w", err)
	}
	var accounts []filesgen.UserAccount
	if err := json.Unmarshal(data, &accounts); err != nil {
		return nil, fmt.Errorf("parse users file: %w", err)
	}
	return accounts, nil
}

func saveUsers(path string, accounts []filesgen.UserAccount) error {
	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
