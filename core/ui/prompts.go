/*
 **********************************************************************
 * -------------------------------------------------------------------
 * Project Name : Abdal 4iProto Cli
 * File Name    : prompts.go
 * Programmer   : Ebrahim Shafiei (EbraSha)
 * Email        : Prof.Shafiei@Gmail.com
 * Created On   : 2026-05-27 20:24:00
 * Description  : Survey/Inquirer-based prompt helpers with validators
 *                for ports, integers, choices, and string fields.
 * -------------------------------------------------------------------
 *
 * "Coding is an engaging and beloved hobby for me. I passionately and insatiably pursue knowledge in cybersecurity and programming."
 * – Ebrahim Shafiei
 *
 **********************************************************************
 */

package ui

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// styledMessage wraps a prompt label with the neon prompt style.
func styledMessage(message string) string {
	return StylePrompt.Render(message)
}

// AskSelect renders a single-choice list prompt.
func AskSelect(message string, options []string, defaultOption string) (string, error) {
	answer := ""
	prompt := &survey.Select{
		Message: styledMessage(message),
		Options: options,
		Default: defaultOption,
	}
	err := survey.AskOne(prompt, &answer, survey.WithValidator(survey.Required))
	if errors.Is(err, terminal.InterruptErr) {
		return "", err
	}
	return answer, err
}

// AskMultiSelect renders a multi-choice list prompt.
func AskMultiSelect(message string, options []string, defaults []string) ([]string, error) {
	var answers []string
	prompt := &survey.MultiSelect{
		Message: styledMessage(message),
		Options: options,
		Default: defaults,
	}
	err := survey.AskOne(prompt, &answers)
	return answers, err
}

// AskConfirm asks a yes/no question.
func AskConfirm(message string, defaultYes bool) (bool, error) {
	answer := defaultYes
	prompt := &survey.Confirm{
		Message: styledMessage(message),
		Default: defaultYes,
	}
	err := survey.AskOne(prompt, &answer)
	return answer, err
}

// AskString asks for a plain free-text string.
func AskString(message string, defaultValue string, required bool) (string, error) {
	answer := ""
	prompt := &survey.Input{
		Message: styledMessage(message),
		Default: defaultValue,
	}
	opts := []survey.AskOpt{}
	if required {
		opts = append(opts, survey.WithValidator(survey.Required))
	}
	err := survey.AskOne(prompt, &answer, opts...)
	return strings.TrimSpace(answer), err
}

// AskPassword asks for a hidden password field.
func AskPassword(message string, required bool) (string, error) {
	answer := ""
	prompt := &survey.Password{
		Message: styledMessage(message),
	}
	opts := []survey.AskOpt{}
	if required {
		opts = append(opts, survey.WithValidator(survey.Required))
	}
	err := survey.AskOne(prompt, &answer, opts...)
	return answer, err
}

// AskInt asks for an integer constrained between min and max.
func AskInt(message string, defaultValue int, min int, max int) (int, error) {
	def := strconv.Itoa(defaultValue)
	for {
		answer := ""
		prompt := &survey.Input{
			Message: styledMessage(message),
			Default: def,
		}
		if err := survey.AskOne(prompt, &answer, survey.WithValidator(survey.Required)); err != nil {
			return 0, err
		}
		n, err := strconv.Atoi(strings.TrimSpace(answer))
		if err != nil {
			Warning(fmt.Sprintf("'%s' is not a valid integer. Please try again.", answer))
			continue
		}
		if n < min || n > max {
			Warning(fmt.Sprintf("Value must be between %d and %d.", min, max))
			continue
		}
		return n, nil
	}
}

// AskPort asks for a TCP port between 1 and 65535 (validation is local;
// reservation checks are performed by the caller through network package).
func AskPort(message string, defaultValue int) (int, error) {
	return AskInt(message, defaultValue, 1, 65535)
}
