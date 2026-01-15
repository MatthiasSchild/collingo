package config

import (
	"collingo/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

const (
	UserConfigFileName = ".collingo"
)

var (
	ErrUserConfigNotAuthenticated = errors.New("login is required")
)

type UserConfig struct {
	ApiToken  string `json:"apiToken,omitempty"`
	ServerUrl string `json:"serverUrl,omitempty"`
}

func LoadUserConfigFromFile(homeDir string) (*UserConfig, error) {
	filename := path.Join(homeDir, UserConfigFileName)

	// If file does not exist, create an empty config
	if !utils.FileExists(filename) {
		return &UserConfig{}, nil
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config UserConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadUserConfigFromFileRequiresAuth(homeDir string) (*UserConfig, error) {
	userConfig, err := LoadUserConfigFromFile(homeDir)
	if err != nil {
		return nil, err
	}
	if err = userConfig.RequiresAuth(); err != nil {
		return nil, err
	}
	return userConfig, nil
}

func (u *UserConfig) RequiresAuth() error {
	if u.ApiToken == "" {
		return ErrUserConfigNotAuthenticated
	}
	return nil
}

func (u *UserConfig) WriteToFile(homeDir string) error {
	filename := path.Join(homeDir, UserConfigFileName)

	data, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal user config: %v", err)
	}

	err = os.WriteFile(filename, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write user config file: %v", err)
	}

	return err
}
