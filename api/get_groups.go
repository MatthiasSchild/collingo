package api

import (
	"collingo/config"
	"collingo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetGroup(userConfig *config.UserConfig, baseUrl string, project string, group string) (models.GroupModel, error) {
	path := fmt.Sprintf("/api/v1/projects/%s/groups/%s", project, group)
	req, err := prepareGetRequest(userConfig, baseUrl, path)
	if err != nil {
		return models.GroupModel{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.GroupModel{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.GroupModel{}, handleErrorResponse(resp)
	}

	var result models.ResultModel[models.GroupModel]
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return models.GroupModel{}, err
	}

	return result.Result, nil
}
