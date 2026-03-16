package api

import (
	"collingo/config"
	"fmt"
	"io"
	"net/http"
)

func ExportVueI18n(userConfig *config.UserConfig, baseUrl string, project string, format bool) (string, error) {
	path := fmt.Sprintf("/api/v1/projects/%s/export/vue-i18n", project)
	req, err := prepareGetRequestWithFormat(userConfig, baseUrl, path, format)
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
