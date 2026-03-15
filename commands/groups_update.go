package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/dialogs"
	"collingo/partials"
	"collingo/utils"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var GroupsUpdateCmd = &cobra.Command{
	Use:           "update",
	Aliases:       []string{"edit"},
	Short:         "Update the data of a group",
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

		// Get the group
		group, err := partials.GetGroupFromCommand(
			userConfig,
			workspaceConfig,
			baseUrl,
			cmd,
		)
		if err != nil {
			return err
		}

		// Update Parent
		updateParent, err := cmd.Flags().GetString("parent")
		if err != nil {
			return err
		} else if updateParent != "" && !utils.IsObjectID(updateParent) {
			return errors.New("The id is not a valid id")
		}
		clearParent, err := cmd.Flags().GetBool("root-group")
		if err != nil {
			return err
		}
		if len(updateParent) > 0 && clearParent {
			return errors.New("parent and root-group can not be used at the same time")
		}
		if len(updateParent) == 0 && !clearParent {
			selectParent, err := dialogs.Confirm("Do you want to set a parent of the group?")
			if err != nil {
				return err
			}
			if selectParent {
				newParent, err := dialogs.GroupSelection(userConfig, baseUrl, projectId, group.ID)
				if err != nil {
					return err
				}
				updateParent = newParent.ID
			} else if group.Parent != nil {
				removeParent, err := dialogs.Confirm("Do you want to remove the parent then instead?")
				if err != nil {
					return err
				}

				if removeParent {
					clearParent = true
				}
			}
		}

		// Update display name
		updateDisplayName, err := cmd.Flags().GetString("display-name")
		if err != nil {
			return err
		}
		if updateDisplayName == "" {
			shouldUpdateDisplayName, err := dialogs.ConfirmF(
				"Do you want to update the display name (currently '%s')?",
				group.DisplayName,
			)
			if err != nil {
				return err
			}

			if shouldUpdateDisplayName {
				updateDisplayName = console.StringRequired("Please enter a new display name")
			}
		}

		// Update technical name
		updateTechnicalName, err := cmd.Flags().GetString("technical-name")
		if err != nil {
			return err
		}
		if updateTechnicalName == "" {
			shouldUpdateTechnicalName, err := dialogs.ConfirmF(
				"Do you want to update the technical name (currently '%s')?",
				group.TechnicalName,
			)
			if err != nil {
				return err
			}

			if shouldUpdateTechnicalName {
				updateTechnicalName = console.TechnicalName("Please enter a new technical name")
			}
		}

		// Update the group
		input := api.UpdateGroupInput{}
		updateAnything := false
		if updateParent != "" {
			input.Parent = utils.NewNullableString(updateParent)
			updateAnything = true
		} else if clearParent {
			input.Parent = utils.NullableStringNull()
			updateAnything = true
		}
		if updateDisplayName != "" {
			input.DisplayName = &updateDisplayName
			updateAnything = true
		}
		if updateTechnicalName != "" {
			input.TechnicalName = &updateTechnicalName
			updateAnything = true
		}

		if !updateAnything {
			console.Info("Nothing to update")
			return nil
		}
		groupResult, err := api.UpdateGroup(userConfig, baseUrl, projectId, group.ID, input)
		if err != nil {
			return err
		}

		console.SuccessF("Successfully updated group '%s'", groupResult.DisplayName)
		return nil
	},
}

func init() {
	GroupsCmd.AddCommand(GroupsUpdateCmd)
	GroupsUpdateCmd.Flags().String("group", "", "The path of the group to edit")
	GroupsUpdateCmd.Flags().String("group-id", "", "The ID of the group to edit")
	GroupsUpdateCmd.Flags().String("parent", "", "The ID of the new parent")
	GroupsUpdateCmd.Flags().Bool("root-group", false, "Put the group as root without a parent")
	GroupsUpdateCmd.Flags().String("display-name", "", "The new display name")
	GroupsUpdateCmd.Flags().String("technical-name", "", "The new technical name")
	GroupsUpdateCmd.Flags().String("working-dir", "", "Set the working directory")
}
