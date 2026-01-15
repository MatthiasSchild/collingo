package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/dialogs"
	"collingo/utils"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var ProjectsDeleteCmd = &cobra.Command{
	Use:           "delete",
	Aliases:       []string{"rm"},
	Short:         "Delete one of your projects",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, _ := os.UserHomeDir()
		userConfig, err := config.LoadUserConfigFromFileRequiresAuth(homeDir)
		if err != nil {
			return err
		}

		// Get the ID of the project
		projectId, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		} else if projectId == "" {
			project, err := dialogs.ProjectSelection(userConfig)
			if err != nil {
				return err
			}
			projectId = project.ID
		}
		if !utils.IsObjectID(projectId) {
			return errors.New("The id is not a valid id")
		}

		// Get and display the project, ask to confirm
		noConfirm, _ := cmd.Flags().GetBool("yes")
		if !noConfirm {
			project, err := api.GetProject(userConfig, projectId)
			if err != nil {
				return err
			}
			console.InfoF(
				"Project to delete: %s with ID %s",
				project.Name,
				project.ID,
			)
			ok, err := dialogs.Confirm("Do you really want to delete this project?")
			if err != nil {
				return err
			}
			if !ok {
				console.Info("Cancelled")
				return nil
			}
		}

		// Delete the project
		err = api.DeleteProject(userConfig, projectId)
		if err != nil {
			return err
		}

		console.Success("Project has been successfully deleted")
		return nil
	},
}

func init() {
	ProjectsCmd.AddCommand(ProjectsDeleteCmd)
	ProjectsDeleteCmd.Flags().String("id", "", "The ID of the project to delete")
	ProjectsDeleteCmd.Flags().Bool("yes", false, "Confirm and don't ask")
}
