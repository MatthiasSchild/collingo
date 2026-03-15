package api

import (
	"collingo/config"
	"fmt"
	"net/http"
)

func DeleteGroup(userConfig *config.UserConfig, baseUrl string, project string, group string) error {
	path := fmt.Sprintf("/api/v1/projects/%s/groups/%s", project, group)
	req, err := prepareDeleteRequest(userConfig, baseUrl, path)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return handleErrorResponse(resp)
	}

	return nil
}
