package api

import (
	"collingo/config"
	"collingo/models"
	"encoding/json"
	"net/http"
)

func ListProjects(userConfig *config.UserConfig, baseUrl string, limit uint32, offset uint32) (models.ManyResultModel[models.ProjectModel], error) {
	var result models.ManyResultModel[models.ProjectModel]

	req, err := prepareGetRequestWithPagination(userConfig, baseUrl, "/api/v1/projects", limit, offset)
	if err != nil {
		return result, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
