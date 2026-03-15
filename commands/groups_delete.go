package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/dialogs"
	"collingo/partials"
	"os"

	"github.com/spf13/cobra"
)

var GroupsDeleteCmd = &cobra.Command{
	Use:           "delete",
	Aliases:       []string{"rm"},
	Short:         "Delete a group of a project",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, _ := os.UserHomeDir()
		userConfig, err := config.LoadUserConfigFromFileRequiresAuth(homeDir)
		if err != nil {
			return err
		}

		// Get current project
		workingDir := partials.WorkingDirFromFlags(cmd, "working-dir")
		workspaceConfig, err := config.LoadWorkspaceConfigFromFile(workingDir)
		if err != nil {
			return err
		}
		projectId := workspaceConfig.ProjectId
		baseUrl := config.EffectiveServerUrl(userConfig, workspaceConfig)

		// Get the ID of the group
		group, err := partials.GetGroupFromCommand(
			userConfig,
			workspaceConfig,
			baseUrl,
			cmd,
		)
		if err != nil {
			return err
		}

		// Get and display the group, ask to confirm
		noConfirm, _ := cmd.Flags().GetBool("yes")
		if !noConfirm {
			console.InfoF(
				"Group to delete: '%s' (%s)",
				group.DisplayName,
				group.ID,
			)

			ok, err := dialogs.Confirm("Do you really want to delete this group?")
			if err != nil {
				return err
			}
			if !ok {
				console.Info("Cancelled")
				return nil
			}
		}

		// Delete the group
		err = api.DeleteGroup(userConfig, baseUrl, projectId, group.ID)
		if err != nil {
			return err
		}

		console.SuccessF(
			"Group '%s' successfully delete from project '%s'",
			group.DisplayName,
			projectId,
		)
		return nil
	},
}

func init() {
	GroupsCmd.AddCommand(GroupsDeleteCmd)
	GroupsDeleteCmd.Flags().String("working-dir", "", "Set the working directory")
	GroupsDeleteCmd.Flags().String("group", "", "The path of the group (e.g. main.footer)")
	GroupsDeleteCmd.Flags().String("group-id", "", "The id of the group")
	GroupsDeleteCmd.Flags().Bool("yes", false, "Confirm and don't ask")
}
