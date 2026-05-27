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

// UpdateInput describes the editable fields of an account. Pointer types
// are used so callers can choose to update a subset; nil pointers mean
// "leave this field untouched".
type UpdateInput struct {
	Username       string
	NewPassword    *string
	NewRole        *string
	NewMaxSessions *int
	NewMaxSpeed    *int
	NewMaxTotalMB  *int
	NewLog         *string
	NewBlockedDom  *[]string
	NewBlockedIPs  *[]string
}

// UpdateUser modifies an existing account in users.json and restarts the
// Abdal 4iProto Server service when at least one field was changed.
func UpdateUser(in UpdateInput) error {
	in.Username = strings.TrimSpace(in.Username)
	if in.Username == "" {
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

	idx := -1
	for i := range accounts {
		if strings.EqualFold(accounts[i].Username, in.Username) {
			idx = i
			break
		}
	}
	if idx < 0 {
		return fmt.Errorf("user '%s' not found", in.Username)
	}

	target := &accounts[idx]

	if in.NewPassword != nil {
		if strings.TrimSpace(*in.NewPassword) == "" {
			return fmt.Errorf("password cannot be empty")
		}
		target.Password = *in.NewPassword
	}
	if in.NewRole != nil {
		role := strings.TrimSpace(*in.NewRole)
		if role != config.UserRoleAdmin && role != config.UserRoleUser {
			return fmt.Errorf("role must be '%s' or '%s'", config.UserRoleAdmin, config.UserRoleUser)
		}
		target.Role = role
	}
	if in.NewMaxSessions != nil {
		v := *in.NewMaxSessions
		if v < config.MinSessions || v > config.MaxSessions {
			return fmt.Errorf("max_sessions must be between %d and %d", config.MinSessions, config.MaxSessions)
		}
		target.MaxSessions = v
	}
	if in.NewMaxSpeed != nil {
		v := *in.NewMaxSpeed
		if v < config.MinSpeedKbps || v > config.MaxSpeedKbps {
			return fmt.Errorf("max_speed_kbps must be between %d and %d", config.MinSpeedKbps, config.MaxSpeedKbps)
		}
		target.MaxSpeedKbps = v
	}
	if in.NewMaxTotalMB != nil {
		if *in.NewMaxTotalMB < 0 {
			return fmt.Errorf("max_total_mb cannot be negative")
		}
		target.MaxTotalMB = *in.NewMaxTotalMB
	}
	if in.NewLog != nil {
		log := strings.ToLower(strings.TrimSpace(*in.NewLog))
		if log != "yes" && log != "no" {
			return fmt.Errorf("log must be 'yes' or 'no'")
		}
		target.Log = log
	}
	if in.NewBlockedDom != nil {
		target.BlockedDomains = normalizeList(*in.NewBlockedDom)
	}
	if in.NewBlockedIPs != nil {
		target.BlockedIPs = normalizeList(*in.NewBlockedIPs)
	}

	if err := saveUsers(path, accounts); err != nil {
		return err
	}
	ui.Success(fmt.Sprintf("User '%s' updated", target.Username))
	return service.Restart(service.ComponentServer)
}

// normalizeList trims spaces and drops empty entries while keeping order.
func normalizeList(items []string) []string {
	out := make([]string, 0, len(items))
	for _, raw := range items {
		v := strings.TrimSpace(raw)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
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
