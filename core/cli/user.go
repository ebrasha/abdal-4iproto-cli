/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : user.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Cobra user management commands (add/remove accounts).
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package cli

import (
	"github.com/spf13/cobra"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/users"
)

func newUserCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "user",
		Short: "Manage SSH server user accounts",
	}

	var (
		username     string
		password     string
		role         string
		maxSessions  int
		maxSpeedKbps int
		maxTotalMB   int
	)

	add := &cobra.Command{
		Use:   "add",
		Short: "Add a new user to users.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			return users.AddUser(users.AddInput{
				Username: username, Password: password, Role: role,
				MaxSessions: maxSessions, MaxSpeedKbps: maxSpeedKbps, MaxTotalMB: maxTotalMB,
			})
		},
	}
	add.Flags().StringVar(&username, "username", "", "Account username (required)")
	add.Flags().StringVar(&password, "password", "", "Account password (required)")
	add.Flags().StringVar(&role, "role", config.UserRoleUser, "Role: admin|user")
	add.Flags().IntVar(&maxSessions, "max-sessions", 1, "Max concurrent sessions (1-10000)")
	add.Flags().IntVar(&maxSpeedKbps, "max-speed-kbps", 512, "Bandwidth cap in Kbps")
	add.Flags().IntVar(&maxTotalMB, "max-total-mb", 0, "Total transfer cap in MB (0=unlimited)")
	_ = add.MarkFlagRequired("username")
	_ = add.MarkFlagRequired("password")

	remove := &cobra.Command{
		Use:   "remove",
		Short: "Remove a user from users.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			return users.RemoveUser(username)
		},
	}
	remove.Flags().StringVar(&username, "username", "", "Username to remove (required)")
	_ = remove.MarkFlagRequired("username")

	root.AddCommand(add, remove)
	return root
}
