package api

import (
	"collingo/config"
	"collingo/models"
	"encoding/json"
	"net/http"
)

type CreateProjectInput struct {
	Name                 string   `json:"name"`
	BaseLanguage         string   `json:"baseLanguage"`
	TranslationLanguages []string `json:"translationLanguages"`
}

func CreateProject(userConfig *config.UserConfig, input CreateProjectInput) (models.ProjectModel, error) {
	req, err := preparePostRequest(userConfig, "/api/v1/projects", input)
	if err != nil {
		return models.ProjectModel{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.ProjectModel{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ProjectModel{}, handleErrorResponse(resp)
	}

	var result models.ResultModel[models.ProjectModel]
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return models.ProjectModel{}, err
	}

	return result.Result, nil
}
