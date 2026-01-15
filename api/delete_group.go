package api

import (
	"collingo/config"
	"fmt"
	"net/http"
)

func DeleteGroup(userConfig *config.UserConfig, project string, group string) error {
	url := fmt.Sprintf("/api/v1/projects/%s/groups/%s", project, group)
	req, err := prepareDeleteRequest(userConfig, url)
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
