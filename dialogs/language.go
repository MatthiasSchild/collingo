package dialogs

import (
	"collingo/console"
	"strings"
)

var langRegex = `^[a-z]{2}-[A-Z]{2}$`
var langsRegex = `^([a-z]{2}-[A-Z]{2})?$`

func LanguageSelection(label string) (string, error) {
	lang := console.StringRegex(label, langRegex)
	return lang, nil
}

func MultiLanguageSelection(label string) ([]string, error) {
	var result []string

	for {
		newEntry := console.StringRegex(label+" (Finish with empty line)", langsRegex)
		if newEntry == "" {
			return result, nil
		}

		result = append(result, newEntry)
		console.InfoF("Current list: %s", strings.Join(result, ", "))
	}
}
