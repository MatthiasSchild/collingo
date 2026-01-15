package api

import (
	"collingo/config"
	"collingo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func Info(config *config.UserConfig) (models.InfoModel, error) {
	var result models.InfoModel

	path, err := buildRequestUrl(config, "/api/v1/info")
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return result, err
	}

	req.SetBasicAuth("token", config.ApiToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return result, fmt.Errorf("API returned with status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
