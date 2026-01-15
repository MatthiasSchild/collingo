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
	WorkspaceConfigFileName = ".collingo.json"
)

var (
	ErrWorkspaceConfigNotFound = errors.New("workspace config not found")
)

type WorkspaceConfig struct {
	ProjectId string `json:"project,omitempty"`
	Template  string `json:"template,omitempty"`
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
