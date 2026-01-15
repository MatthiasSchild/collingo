package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/dialogs"
	"collingo/models"
	"collingo/utils"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var ProjectsUpdateCmd = &cobra.Command{
	Use:           "update",
	Aliases:       []string{"edit"},
	Short:         "Update the data of a project",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, _ := os.UserHomeDir()
		userConfig, err := config.LoadUserConfigFromFileRequiresAuth(homeDir)
		if err != nil {
			return err
		}

		// Get the project
		var project *models.ProjectModel
		projectId, err := cmd.Flags().GetString("project")
		if err != nil {
			return err
		} else if projectId == "" {
			project2, err := dialogs.ProjectSelection(userConfig)
			if err != nil {
				return err
			}
			project = &project2
			projectId = project.ID
		}
		if !utils.IsObjectID(projectId) {
			return errors.New("The id is not a valid id")
		}
		if project == nil {
			project2, err := api.GetProject(userConfig, projectId)
			if err != nil {
				return err
			}
			project = &project2
		}

		// Update name
		updateName, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		if updateName == "" {
			shouldUpdateName, err := dialogs.ConfirmF(
				"Do you want to update the name (currently '%s')?",
				project.Name,
			)
			if err != nil {
				return err
			}

			if shouldUpdateName {
				updateName = console.String("Please enter a new name")
			}
		}

		// Update the project
		input := api.UpdateProjectInput{}
		updateAnything := false
		if updateName != "" {
			input.Name = &updateName
			updateAnything = true
		}

		if !updateAnything {
			console.Info("Nothing to update")
			return nil
		}
		projectResult, err := api.UpdateProject(userConfig, projectId, input)
		if err != nil {
			return err
		}

		console.SuccessF("Successfully updated project '%s'", projectResult.Name)
		return nil
	},
}

func init() {
	ProjectsCmd.AddCommand(ProjectsUpdateCmd)
	ProjectsUpdateCmd.Flags().String("project", "", "The ID of the project to edit")
	ProjectsUpdateCmd.Flags().String("name", "", "The new name of the project")
}
