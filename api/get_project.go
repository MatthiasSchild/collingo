package api

import (
	"collingo/config"
	"collingo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetProject(userConfig *config.UserConfig, id string) (models.ProjectModel, error) {
	url := fmt.Sprintf("/api/v1/projects/%s", id)
	req, err := prepareGetRequest(userConfig, url)
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
