package console

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"slices"

	"github.com/manifoldco/promptui"
)

func String(label string) string {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()
	if errors.Is(err, promptui.ErrInterrupt) {
		os.Exit(0)
	} else if err != nil {
		panic(err)
	}

	return result
}

func StringRequired(label string) string {
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(s string) error {
			if len(s) == 0 {
				return fmt.Errorf("Please enter a valid value")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if errors.Is(err, promptui.ErrInterrupt) {
		os.Exit(0)
	} else if err != nil {
		panic(err)
	}

	return result
}

func StringRegex(label string, pattern string) string {
	re := regexp.MustCompile(pattern)

	prompt := promptui.Prompt{
		Label: label,
		Validate: func(s string) error {
			if !re.MatchString(s) {
				return fmt.Errorf("Please enter a valid value")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if errors.Is(err, promptui.ErrInterrupt) {
		os.Exit(0)
	} else if err != nil {
		panic(err)
	}

	return result
}

func TechnicalName(label string) string {
	return StringRegex(label, `^[a-z][a-zA-Z0-9]{0,31}$`)
}

func TechnicalNameExcept(label string, except []string) string {
	re := regexp.MustCompile(`^[a-z][a-zA-Z0-9]{0,31}$`)

	prompt := promptui.Prompt{
		Label: label,
		Validate: func(s string) error {
			if !re.MatchString(s) {
				return fmt.Errorf("Please enter a valid value")
			}
			if slices.Contains(except, s) {
				return fmt.Errorf("An entry with this technical name already exists")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if errors.Is(err, promptui.ErrInterrupt) {
		os.Exit(0)
	} else if err != nil {
		panic(err)
	}

	return result

}
