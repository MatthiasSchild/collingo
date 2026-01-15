package api

import (
	"collingo/config"
	"collingo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type UpdateProjectInput struct {
	Name *string `json:"name,omitempty"`
}

func UpdateProject(userConfig *config.UserConfig, projectId string, input UpdateProjectInput) (models.ProjectModel, error) {
	url := fmt.Sprintf("/api/v1/projects/%s/", projectId)
	req, err := preparePatchRequest(userConfig, url, input)
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
