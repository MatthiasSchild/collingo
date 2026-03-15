package partials

import (
	"collingo/api"
	"collingo/config"
	"collingo/models"
	"fmt"
	"strings"
)

func ResolveGroupPath(userConfig *config.UserConfig, baseUrl string, project string, path string) (models.GroupModel, error) {
	groupSummaries, err := api.ListGroupSummary(userConfig, baseUrl, project)
	if err != nil {
		return models.GroupModel{}, err
	}

	groupIdIndex := make(map[string]models.GroupSummaryModel)
	for _, groupSummary := range groupSummaries.Result {
		groupIdIndex[groupSummary.ID] = groupSummary
	}

	for _, groupSummary := range groupSummaries.Result {
		parts := make([]string, 0)
		for _, groupId := range groupSummary.Breadcrumbs {
			groupElement, ok := groupIdIndex[groupId]
			if !ok {
				return models.GroupModel{}, fmt.Errorf("invalid group parent id '%s'", groupId)
			}
			parts = append(parts, groupElement.TechnicalName)
		}
		parts = append(parts, groupSummary.TechnicalName)

		currentPath := strings.Join(parts, ".")
		if currentPath == path {
			group, err := api.GetGroup(userConfig, baseUrl, project, groupSummary.ID)
			if err != nil {
				return models.GroupModel{}, err
			}
			return group, nil
		}
	}

	return models.GroupModel{}, fmt.Errorf("No group with technical name path found: %s", path)
}
