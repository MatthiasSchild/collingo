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

var Token string

func buildRequestUrl(baseUrl string, path string) (string, error) {
	if baseUrl == "" {
		baseUrl = config.DefaultServerUrl
	}
	requestUrl, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}
	requestUrl.Path = path
	return requestUrl.String(), nil
}

func buildRequestUrlWithPagination(baseUrl string, path string, limit uint32, offset uint32) (string, error) {
	if baseUrl == "" {
		baseUrl = config.DefaultServerUrl
	}
	requestUrl, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}
	query := requestUrl.Query()
	query.Set("limit", fmt.Sprintf("%d", limit))
	query.Set("offset", fmt.Sprintf("%d", offset))
	requestUrl.Path = path
	requestUrl.RawQuery = query.Encode()
	return requestUrl.String(), nil
}

func prepareGetRequest(config *config.UserConfig, baseUrl string, path string) (*http.Request, error) {
	urlStr, err := buildRequestUrl(baseUrl, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("token", config.ApiToken)

	return req, nil
}

func prepareGetRequestWithPagination(config *config.UserConfig, baseUrl string, path string, limit uint32, offset uint32) (*http.Request, error) {
	urlStr, err := buildRequestUrlWithPagination(baseUrl, path, limit, offset)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("token", config.ApiToken)

	return req, nil
}

func preparePostRequest(config *config.UserConfig, baseUrl string, path string, body any) (*http.Request, error) {
	urlStr, err := buildRequestUrl(baseUrl, path)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("token", config.ApiToken)

	return req, nil
}

func preparePatchRequest(config *config.UserConfig, baseUrl string, path string, body any) (*http.Request, error) {
	urlStr, err := buildRequestUrl(baseUrl, path)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPatch, urlStr, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("token", config.ApiToken)

	return req, nil
}

func prepareDeleteRequest(config *config.UserConfig, baseUrl string, path string) (*http.Request, error) {
	urlStr, err := buildRequestUrl(baseUrl, path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, urlStr, nil)
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
