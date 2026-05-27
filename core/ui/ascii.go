/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : ascii.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : ASCII banner rendering for the Abdal 4iProto Cli with
 *                neon-green colorization and a programmer credit line.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package ui

import (
	"fmt"
	"strings"

	"abdal-4iproto-cli/core/config"
)

// asciiBanner is the fixed banner rendered at every program start.
const asciiBanner = `           _         _       _   _  _   _ _____           _           _____ _ _ 
     /\   | |       | |     | | | || | (_)  __ \         | |         / ____| (_)
    /  \  | |__   __| | __ _| | | || |_ _| |__) | __ ___ | |_ ___   | |    | |_ 
   / /\ \ | '_ \ / _` + "`" + ` |/ _` + "`" + ` | | |__   _| |  ___/ '__/ _ \| __/ _ \  | |    | | |
  / ____ \| |_) | (_| | (_| | |    | | | | |   | | | (_) | || (_) | | |____| | |
 /_/    \_\_.__/ \__,_|\__,_|_|    |_| |_|_|   |_|  \___/ \__\___/   \_____|_|_|
`

// PrintBanner prints the colored ASCII banner followed by the credit line.
func PrintBanner() {
	fmt.Println(StyleBanner.Render(asciiBanner))
	fmt.Println(StyleProgrammer.Render(strings.Repeat(" ", 4) + "Programmer: " + config.ProgrammerName + "  |  " + config.ProgrammerMail))
	fmt.Println(StyleHint.Render(strings.Repeat(" ", 4) + "GitHub: " + config.ProgrammerGH + "  |  Telegram: " + config.ProgrammerTG))
	fmt.Println()
}

// BannerString returns the banner as a string (without printing it).
func BannerString() string {
	return StyleBanner.Render(asciiBanner)
}
