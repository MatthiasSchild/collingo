package dialogs

import (
	"collingo/config"
	"strings"

	"github.com/manifoldco/promptui"
)

// PromptTemplateExportOptions asks for optional default export path and formatting after a template kind was chosen.
// kind may be empty ("no template"); returns (nil, nil).
func PromptTemplateExportOptions(kind string) (*config.TemplateConfig, error) {
	if kind == "" {
		return nil, nil
	}
	t := &config.TemplateConfig{Kind: kind}

	switch kind {
	case config.TemplateKindFlutter, config.TemplateKindI18next:
		dir, err := optionalStringPrompt("Export directory relative to workspace root (empty to skip)")
		if err != nil {
			return nil, err
		}
		t.Directory = strings.TrimSpace(dir)
	case config.TemplateKindVueI18n, config.TemplateKindJSON:
		file, err := optionalStringPrompt("Output file relative to workspace root (empty = stdout for collingo export)")
		if err != nil {
			return nil, err
		}
		t.OutputFile = strings.TrimSpace(file)
	}

	useFormat, err := Confirm("Use formatted export by default?")
	if err != nil {
		return nil, err
	}
	t.Formatted = useFormat

	return t, nil
}

func optionalStringPrompt(label string) (string, error) {
	prompt := promptui.Prompt{Label: label}
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}
