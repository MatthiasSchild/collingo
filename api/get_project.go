package api

import (
	"collingo/config"
	"collingo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetProject(userConfig *config.UserConfig, baseUrl string, id string) (models.ProjectModel, error) {
	path := fmt.Sprintf("/api/v1/projects/%s", id)
	req, err := prepareGetRequest(userConfig, baseUrl, path)
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
