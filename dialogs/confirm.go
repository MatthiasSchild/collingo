package dialogs

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func Confirm(message string) (bool, error) {
	prompt := promptui.Select{
		Label: message,
		Items: []string{"yes", "no"},
	}
	index, _, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return index == 0, nil
}

func ConfirmF(format string, a ...any) (bool, error) {
	message := fmt.Sprintf(format, a...)
	return Confirm(message)
}
