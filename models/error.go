package models

type ErrorResponse struct {
	Error     string            `json:"error"`
	ErrorCode string            `json:"errorCode"`
	Fields    map[string]string `json:"fields"`
}
