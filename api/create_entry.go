package api

import (
	"collingo/config"
	"collingo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateEntryInput struct {
	TechnicalName string `json:"technicalName"`
	BaseTerm      string `json:"baseTerm"`
	ContextInfo   string `json:"contextInfo,omitempty"`
}

func CreateEntry(userConfig *config.UserConfig, project string, group string, input CreateEntryInput) (models.GroupModel, error) {
	url := fmt.Sprintf("/api/v1/projects/%s/groups/%s/entries", project, group)
	req, err := preparePostRequest(userConfig, url, input)
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
