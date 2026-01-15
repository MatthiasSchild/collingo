package api

import (
	"collingo/config"
	"fmt"
	"net/http"
)

func DeleteProject(userConfig *config.UserConfig, id string) error {
	url := fmt.Sprintf("/api/v1/projects/%s", id)
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
