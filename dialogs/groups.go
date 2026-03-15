package dialogs

import (
	"collingo/api"
	"collingo/config"
	"collingo/models"
	"errors"
	"slices"
	"strings"

	"github.com/manifoldco/promptui"
)

func GroupSelection(userConfig *config.UserConfig, baseUrl string, project string, exceptIds ...string) (models.GroupSummaryModel, error) {
	groups := make([]models.GroupSummaryModel, 0)
	more := true
	currentPage := uint32(0)

	for more {
		page, err := api.ListGroupSummary(userConfig, baseUrl, project)
		if err != nil {
			return models.GroupSummaryModel{}, err
		}

		for _, result := range page.Result {
			if slices.Contains(exceptIds, result.ID) {
				continue
			}
			groups = append(groups, result)
		}
		currentPage += 1
		more = page.HasMore
	}

	if len(groups) == 0 {
		return models.GroupSummaryModel{}, errors.New("No groups to select from")
	}

	groupNames := make([]string, 0)
	for _, group := range groups {
		name := strings.Join(append(group.BreadcrumbNames, group.DisplayName), " > ")
		groupNames = append(groupNames, name)
	}

	prompt := promptui.Select{
		Label: "Select a group",
		Items: groupNames,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return models.GroupSummaryModel{}, err
	}

	return groups[index], nil
}
