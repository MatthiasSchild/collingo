package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/dialogs"
	"collingo/partials"
	"collingo/utils"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var EntriesDeleteCmd = &cobra.Command{
	Use:           "delete",
	Aliases:       []string{"rm"},
	Short:         "Delete an entry from a group",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, _ := os.UserHomeDir()
		userConfig, err := config.LoadUserConfigFromFileRequiresAuth(homeDir)
		if err != nil {
			return err
		}

		// Get current project
		workingDir := partials.WorkingDirFromFlags(cmd)
		workspaceConfig, err := config.LoadWorkspaceConfigFromFile(workingDir)
		if err != nil {
			return err
		}
		projectId := workspaceConfig.ProjectId
		baseUrl := config.EffectiveServerUrl(userConfig, workspaceConfig)

		// Get the ID of the group
		groupId, err := cmd.Flags().GetString("group")
		if err != nil {
			return err
		} else if groupId == "" {
			group, err := dialogs.GroupSelection(userConfig, baseUrl, projectId)
			if err != nil {
				return err
			}
			groupId = group.ID
		}
		if !utils.IsObjectID(groupId) {
			return errors.New("The id is not a valid id")
		}

		// Get the technical name of the entry
		technicalName, err := cmd.Flags().GetString("technical-name")
		if err != nil {
			return err
		} else if technicalName == "" {
			technicalName = console.TechnicalName("Please enter the technical name")
		}

		// Get and display the entry, ask to confirm
		noConfirm, _ := cmd.Flags().GetBool("yes")
		if !noConfirm {
			group, err := api.GetGroup(userConfig, baseUrl, projectId, groupId)
			if err != nil {
				return err
			}

			found := false
			for _, entry := range group.Entries {
				if entry.TechnicalName == technicalName {
					console.InfoF(
						"Entry to delete: '%s' (%s) in group '%s' (%s)",
						entry.BaseTerm,
						entry.TechnicalName,
						group.DisplayName,
						group.ID,
					)

					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("the entry has not been found in the group")
			}

			ok, err := dialogs.Confirm("Do you really want to delete the entry?")
			if err != nil {
				return err
			}
			if !ok {
				console.Info("Cancelled")
				return nil
			}
		}

		// Delete the entry
		err = api.DeleteEntry(userConfig, baseUrl, projectId, groupId, technicalName)
		if err != nil {
			return err
		}

		console.SuccessF(
			"Entry '%s' successfully deleted from group '%s'",
			technicalName,
			groupId,
		)
		return nil
	},
}

func init() {
	EntriesCmd.AddCommand(EntriesDeleteCmd)
	EntriesDeleteCmd.Flags().String("group", "", "The id of the group")
	EntriesDeleteCmd.Flags().String("technical-name", "", "The technical name of the entry")
	EntriesDeleteCmd.Flags().Bool("yes", false, "Confirm and don't ask")
}
