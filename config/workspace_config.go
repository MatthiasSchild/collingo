package config

import (
	"collingo/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

const (
	WorkspaceConfigFileName = ".collingo.json"
)

var (
	ErrWorkspaceConfigNotFound = errors.New("workspace config not found")
)

type WorkspaceConfig struct {
	ProjectId string          `json:"project,omitempty"`
	Template  *TemplateConfig `json:"template,omitempty"`
	ServerUrl string          `json:"serverUrl,omitempty"`
}

func LoadWorkspaceConfigFromFile(currentDir string) (*WorkspaceConfig, error) {
	for !utils.FileExists(path.Join(currentDir, WorkspaceConfigFileName)) {
		parent := path.Dir(currentDir)
		if currentDir == parent {
			return nil, ErrWorkspaceConfigNotFound
		}
		currentDir = parent
	}

	filename := path.Join(currentDir, WorkspaceConfigFileName)
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config WorkspaceConfig
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}

// LoadWorkspaceConfigFromFileWithPath loads the workspace config and returns
// the absolute path to the config file. Returns (nil, "", ErrWorkspaceConfigNotFound) when not in a workspace.
func LoadWorkspaceConfigFromFileWithPath(currentDir string) (*WorkspaceConfig, string, error) {
	for !utils.FileExists(path.Join(currentDir, WorkspaceConfigFileName)) {
		parent := path.Dir(currentDir)
		if currentDir == parent {
			return nil, "", ErrWorkspaceConfigNotFound
		}
		currentDir = parent
	}

	filename := path.Join(currentDir, WorkspaceConfigFileName)
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, "", err
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, "", err
	}

	var cfg WorkspaceConfig
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		return nil, "", err
	}

	return &cfg, absPath, nil
}

func (p *WorkspaceConfig) WriteToFile(workingDir string) error {
	filename := path.Join(workingDir, WorkspaceConfigFileName)

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal workspace config: %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write workspace config file: %v", err)
	}

	return err
}
