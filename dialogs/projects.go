package dialogs

import (
	"collingo/api"
	"collingo/config"
	"collingo/models"

	"github.com/manifoldco/promptui"
)

func ProjectSelection(userConfig *config.UserConfig, baseUrl string) (models.ProjectModel, error) {
	projects := make([]models.ProjectModel, 0)
	more := true
	currentPage := uint32(0)

	for more {
		page, err := api.ListProjects(userConfig, baseUrl, 10, currentPage*10)
		if err != nil {
			return models.ProjectModel{}, err
		}

		projects = append(projects, page.Result...)
		currentPage += 1
		more = page.HasMore
	}

	projectNames := make([]string, 0)
	for _, project := range projects {
		projectNames = append(projectNames, project.Name)
	}

	prompt := promptui.Select{
		Label: "Select a project",
		Items: projectNames,
	}
	index, _, err := prompt.Run()
	if err != nil {
		return models.ProjectModel{}, err
	}

	return projects[index], nil
}
