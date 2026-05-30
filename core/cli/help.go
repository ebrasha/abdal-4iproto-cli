/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : help.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Extended help output including programmer credits,
 *                command catalog, flags, and service troubleshooting.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"abdal-4iproto-cli/core/config"
	"abdal-4iproto-cli/core/ui"
)

func newHelpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "Show detailed help, options, and programmer information",
		Run: func(cmd *cobra.Command, args []string) {
			printExtendedHelp(cmd.Root())
		},
	}
}

func printExtendedHelp(root *cobra.Command) {
	ui.SectionHeader(config.AppName + " v" + config.AppVersion + " – Help")

	ui.KeyValueBox("Version", [][2]string{
		{"Application", config.AppName},
		{"Version", config.AppVersion},
		{"Command", config.AppCommandName},
		{"Repository", config.CliRepoURL},
	})

	ui.Box("Programmer", strings.TrimSpace(fmt.Sprintf(`%s
Email   : %s
GitHub  : %s
Telegram: %s
LinkedIn: %s
X/Twitter: %s`, config.ProgrammerName, config.ProgrammerMail, config.ProgrammerGH, config.ProgrammerTG, config.ProgrammerLI, config.ProgrammerX)))

	ui.Box("Interactive Mode", "Run without arguments to open the Survey-driven main menu with colored boxes.")

	ui.Box("Commands", strings.TrimSpace(`
install       Full/partial installation from GitHub releases
uninstall     Remove services and optionally delete install directory
user add      Add SSH user (requires --username --password --role ...)
user remove   Remove SSH user by --username
service       status|restart|diagnostics --component server|panel
config server Update server ports via --ports
config panel  Update panel JSON via --port --username --password
self-install  Copy CLI binary to system path as abdal-4iproto-cli
self-update   Download the latest release from GitHub and replace the running binary
help          This screen
`))

	ui.Box("Install Flags", strings.TrimSpace(`
--server-only --panel-only
--server-ports 64235,64236,64237,64238
--panel-port 52202 --panel-username ebrasha --panel-password <secret>
--key-type ed25519 --key-bits 4096 --key-force --key-file id_ed25519
--no-services
`))

	ui.Box("KeyGen (bundled binary flags)", strings.TrimSpace(`
-f string   private key filename (public becomes <f>.pub)
-force      overwrite existing key files
-t string   rsa | ed25519 | ecdsa (default rsa in upstream; CLI defaults ed25519)
-b int      key size bits (default 4096)
`))

	ui.Box("Service Names", strings.TrimSpace(`
Linux  : abdal-4iproto-server, abdal-4iproto-panel (systemd)
Windows: Abdal4iProtoServer, Abdal4iProtoPanel (sc)
`))

	ui.Box("Troubleshooting", strings.TrimSpace(`
Linux  : systemctl status abdal-4iproto-server
         journalctl -u abdal-4iproto-server -e
Windows: sc query Abdal4iProtoServer
CLI    : abdal-4iproto-cli service diagnostics
`))

	if root != nil {
		fmt.Println(ui.StyleTitle.Render("Cobra Command Tree"))
		fmt.Println(root.UsageString())
	}
}
