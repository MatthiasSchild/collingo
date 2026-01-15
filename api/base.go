package api

import (
	"bytes"
	"collingo/config"
	"collingo/console"
	"collingo/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	serverUrl = "https://collingo.app"
)

var Token string

func buildRequestUrl(config *config.UserConfig, path string) (string, error) {
	var err error
	var requestUrl *url.URL

	if config.ServerUrl != "" {
		requestUrl, err = url.Parse(config.ServerUrl)
		if err != nil {
			return "", err
		}
	} else {
		requestUrl, _ = url.Parse(serverUrl)
	}

	requestUrl.Path = path
	return requestUrl.String(), nil
}

func buildRequestUrlWithPagination(config *config.UserConfig, path string, limit uint32, offset uint32) (string, error) {
	var err error
	var requestUrl *url.URL

	if config.ServerUrl != "" {
		requestUrl, err = url.Parse(config.ServerUrl)
		if err != nil {
			return "", err
		}
	} else {
		requestUrl, _ = url.Parse(serverUrl)
	}

	query := requestUrl.Query()
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("offset", fmt.Sprintf("%d", offset))

	requestUrl.Path = path
	requestUrl.RawQuery = query.Encode()
	return requestUrl.String(), nil
}

func prepareGetRequest(config *config.UserConfig, path string) (*http.Request, error) {
	url, err := buildRequestUrl(config, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("token", config.ApiToken)

	return req, nil
}

func prepareGetRequestWithPagination(config *config.UserConfig, path string, limit uint32, offset uint32) (*http.Request, error) {
	url, err := buildRequestUrlWithPagination(config, path, limit, offset)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("token", config.ApiToken)

	return req, nil
}

func preparePostRequest(config *config.UserConfig, path string, body any) (*http.Request, error) {
	url, err := buildRequestUrl(config, path)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("token", config.ApiToken)

	return req, nil
}

func preparePatchRequest(config *config.UserConfig, path string, body any) (*http.Request, error) {
	url, err := buildRequestUrl(config, path)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("token", config.ApiToken)

	return req, nil
}

func prepareDeleteRequest(config *config.UserConfig, path string) (*http.Request, error) {
	url, err := buildRequestUrl(config, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("token", config.ApiToken)

	return req, nil
}

func handleErrorResponse(resp *http.Response) error {
	errorResponse := models.ErrorResponse{}
	err := json.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		console.Error(errors.New("failed to decode error response"))
	} else {
		console.InfoF("Status Code: %d", resp.StatusCode)
		console.Info("Error: " + errorResponse.Error)
		console.Info("Error Code: " + errorResponse.ErrorCode)
		if len(errorResponse.Fields) > 0 {
			console.Info("Fields:")
			for field, error := range errorResponse.Fields {
				console.InfoF("  - %s: %s", field, error)
			}
		}
	}

	return errors.New("request failed")
}
