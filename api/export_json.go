package api

import (
	"collingo/config"
	"fmt"
	"io"
	"net/http"
)

func ExportJson(userConfig *config.UserConfig, project string) (string, error) {
	url := fmt.Sprintf("/api/v1/projects/%s/export/json", project)
	req, err := prepareGetRequest(userConfig, url)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", handleErrorResponse(resp)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
