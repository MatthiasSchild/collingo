package api

import (
	"collingo/config"
	"collingo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func ListGroupSummary(userConfig *config.UserConfig, project string) (models.ManyResultModel[models.GroupSummaryModel], error) {
	var result models.ManyResultModel[models.GroupSummaryModel]

	url := fmt.Sprintf("api/v1/projects/%s/group-summary", project)
	req, err := prepareGetRequest(userConfig, url)
	if err != nil {
		return result, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return result, handleErrorResponse(resp)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
