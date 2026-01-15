package commands

import (
	"collingo/config"
	"collingo/console"
	"collingo/dialogs"
	"collingo/utils"
	"errors"
	"os"
	"path"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:           "init",
	Short:         "Initialize project directory",
	Long:          "Initialize the current directory as project directory",
	GroupID:       "setup",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		update, err := cmd.Flags().GetBool("update")
		if err != nil {
			return err
		}

		if update {
			return initCmdUpdate()
		}
		return initCmdNew()
	},
}

func init() {
	RootCmd.AddCommand(InitCmd)
	InitCmd.Flags().Bool("update", false, "Update the current project")
}

func initCmdNew() error {
	homeDir, _ := os.UserHomeDir()
	currentDir, _ := os.Getwd()

	// Check if current dir is a project
	if utils.FileExists(path.Join(currentDir, config.WorkspaceConfigFileName)) {
		return errors.New("current directory is already initialized")
	}

	// Get user config
	userConfig, err := config.LoadUserConfigFromFileRequiresAuth(homeDir)
	if err != nil {
		return err
	}

	// Select project
	project, err := dialogs.ProjectSelection(userConfig)
	if err != nil {
		return err
	}

	// Select template
	template, err := dialogs.TemplateSelection()
	if err != nil {
		return err
	}

	workspaceConfig := config.WorkspaceConfig{
		ProjectId: project.ID,
		Template:  template,
	}
	err = workspaceConfig.WriteToFile(currentDir)
	if err != nil {
		return err
	}

	console.Success("Workspace configuration successfully initialized")
	return nil
}

func initCmdUpdate() error {
	homeDir, _ := os.UserHomeDir()
	currentDir, _ := os.Getwd()

	// Get user config
	userConfig, err := config.LoadUserConfigFromFileRequiresAuth(homeDir)
	if err != nil {
		return err
	}

	// Get the current configuration
	workspaceConfig, err := config.LoadWorkspaceConfigFromFile(currentDir)
	if err != nil {
		return err
	}

	// Select project
	updateProject, err := dialogs.Confirm("Do you want to update the project?")
	if err != nil {
		return err
	}
	if updateProject {
		project, err := dialogs.ProjectSelection(userConfig)
		if err != nil {
			return err
		}
		workspaceConfig.ProjectId = project.ID
	}

	// Select template
	updateTemplate, err := dialogs.Confirm("Do you want to update the template?")
	if err != nil {
		return err
	}
	if updateTemplate {
		template, err := dialogs.TemplateSelection()
		if err != nil {
			return err
		}
		workspaceConfig.Template = template
	}

	err = workspaceConfig.WriteToFile(currentDir)
	if err != nil {
		return err
	}

	console.Success("Workspace configuration successfully updated")
	return nil
}
