package dialogs

import (
	"collingo/api"
	"collingo/config"
	"collingo/models"
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"
)

func EntrySelection(userConfig *config.UserConfig, baseUrl string, projectId string, groupId string) (models.EntryModel, error) {
	group, err := api.GetGroup(userConfig, baseUrl, projectId, groupId)
	if err != nil {
		return models.EntryModel{}, err
	}

	return EntrySelectionFromGroup(group)
}

func EntrySelectionFromGroup(group models.GroupModel) (models.EntryModel, error) {
	if len(group.Entries) == 0 {
		return models.EntryModel{}, errors.New("No entries to select from")
	}

	entries := group.Entries
	entryNames := make([]string, 0)
	for _, entry := range entries {
		name := fmt.Sprintf("%s (%s)", entry.TechnicalName, entry.BaseTerm)
		entryNames = append(entryNames, name)
	}

	prompt := promptui.Select{
		Label: "Select an entry",
		Items: entryNames,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return models.EntryModel{}, err
	}

	return entries[index], nil
}
