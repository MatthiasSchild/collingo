package dialogs

import (
	"github.com/manifoldco/promptui"
)

type templateDefinition struct {
	Key  string
	Name string
}

var templateDefinitions = []templateDefinition{
	{Key: "", Name: "No template"},
	{Key: "vue-i18n", Name: "Vue i18n"},
	{Key: "flutter", Name: "Flutter"},
	{Key: "json", Name: "JSON"},
	{Key: "i18next", Name: "i18next"},
}

func TemplateSelection() (string, error) {
	templateNames := make([]string, 0)
	for _, templateDefinition := range templateDefinitions {
		templateNames = append(templateNames, templateDefinition.Name)
	}

	prompt := promptui.Select{
		Label: "Select a template",
		Items: templateNames,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return templateDefinitions[index].Key, nil
}
