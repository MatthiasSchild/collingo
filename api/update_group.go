package api

import (
	"collingo/config"
	"collingo/models"
	"collingo/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type UpdateGroupInput struct {
	Parent        *utils.NullableString `json:"parent,omitempty"`
	DisplayName   *string               `json:"displayName,omitempty"`
	TechnicalName *string               `json:"technicalName,omitempty"`
}

func UpdateGroup(userConfig *config.UserConfig, projectId string, groupId string, input UpdateGroupInput) (models.GroupModel, error) {
	url := fmt.Sprintf("/api/v1/projects/%s/groups/%s/", projectId, groupId)
	req, err := preparePatchRequest(userConfig, url, input)
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
