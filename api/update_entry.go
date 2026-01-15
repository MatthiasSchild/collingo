package api

import (
	"collingo/config"
	"collingo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type UpdateEntryInput struct {
	BaseTerm    *string `json:"baseTerm,omitempty"`
	ContextInfo *string `json:"contextInfo,omitempty"`
}

func UpdateEntry(
	userConfig *config.UserConfig,
	projectId string,
	groupId string,
	technicalName string,
	input UpdateEntryInput,
) (models.GroupModel, error) {
	url := fmt.Sprintf(
		"/api/v1/projects/%s/groups/%s/entries/%s",
		projectId,
		groupId,
		technicalName,
	)
	req, err := preparePatchRequest(userConfig, url, input)
	if err != nil {
		return models.GroupModel{}, nil
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
