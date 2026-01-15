package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/dialogs"
	"collingo/models"
	"collingo/partials"
	"collingo/utils"
	"errors"
	"os"
	"slices"

	"github.com/spf13/cobra"
)

var EntriesUpdateCmd = &cobra.Command{
	Use:           "update",
	Aliases:       []string{"edit"},
	Short:         "Update the data of an entry",
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

		// Get the group
		group, err := partials.GetGroupFromCommand(
			userConfig,
			workspaceConfig,
			cmd,
		)
		if err != nil {
			return err
		}

		// Get the entry
		var entry models.EntryModel
		var index int
		technicalName, err := cmd.Flags().GetString("entry")
		if err != nil {
			return err
		}
		if len(technicalName) > 0 {
			index = slices.IndexFunc(group.Entries, func(e models.EntryModel) bool {
				return e.TechnicalName == technicalName
			})
			if index == -1 {
				return errors.New("Entry not found in group")
			}
			entry = group.Entries[index]
		} else {
			entry, err = dialogs.EntrySelectionFromGroup(group)
			if err != nil {
				return err
			}
		}
		originalEntry := entry

		// Check if should be moved to new group
		var newGroup *models.GroupModel
		updateGroup, err := cmd.Flags().GetString("new-group")
		if err != nil {
			return err
		}
		updateGroupId, err := cmd.Flags().GetString("new-group-id")
		if err != nil {
			return err
		}
		if updateGroup != "" && updateGroupId != "" {
			return errors.New("new-group and new-group-id can't be used at the same time")
		}
		if updateGroup != "" {
			newGroupResolved, err := partials.ResolveGroupPath(userConfig, workspaceConfig.ProjectId, updateGroup)
			if err != nil {
				return err
			}
			newGroup = &newGroupResolved
		}
		if updateGroupId != "" {
			newGroupResolved, err := api.GetGroup(userConfig, workspaceConfig.ProjectId, updateGroupId)
			if err != nil {
				return err
			}
			newGroup = &newGroupResolved
		}
		if updateGroup == "" && updateGroupId == "" {
			shouldMove, err := dialogs.Confirm("Do you want to move the entry to a different group?")
			if err != nil {
				return err
			}
			if shouldMove {
				groupSummary, err := dialogs.GroupSelection(userConfig, workspaceConfig.ProjectId, group.ID)
				if err != nil {
					return err
				}
				newGroupSelection, err := api.GetGroup(userConfig, workspaceConfig.ProjectId, groupSummary.ID)
				if err != nil {
					return err
				}
				newGroup = &newGroupSelection
			}
		}

		// Update technical name
		updateTechnicalName, err := cmd.Flags().GetString("technical-name")
		if err != nil {
			return err
		}
		if updateTechnicalName != "" {
			// Check if it is a valid technical name
			if !utils.IsTechnicalName(updateTechnicalName) {
				return errors.New("The provided technical name is not valid")
			}
			entry.TechnicalName = updateTechnicalName
		} else {
			shouldUpdateTechnicalName, err := dialogs.Confirm("Do you want to update the technical name?")
			if err != nil {
				return err
			}
			if shouldUpdateTechnicalName {
				// Prevent name collision with other entries
				var except []string
				checkGroup := &group
				if newGroup != nil {
					checkGroup = newGroup
				}
				for i, otherEntry := range checkGroup.Entries {
					if i == index {
						continue
					}
					except = append(except, otherEntry.TechnicalName)
				}

				updateTechnicalName = console.TechnicalNameExcept("Please enter the new technical name", except)
				entry.TechnicalName = updateTechnicalName
			}
		}

		// Update base term
		updateBaseTerm, err := cmd.Flags().GetString("base-term")
		if err != nil {
			return err
		}
		if updateBaseTerm != "" {
			entry.BaseTerm = updateBaseTerm
		} else {
			shouldUpdateBaseTerm, err := dialogs.Confirm("Do you want to update the base term?")
			if err != nil {
				return err
			}
			if shouldUpdateBaseTerm {
				updateBaseTerm = console.StringRequired("Please enter the new base term")
				entry.BaseTerm = updateBaseTerm
			}
		}

		// Update the context
		updateContext, err := cmd.Flags().GetString("context")
		if err != nil {
			return err
		}
		clearContext, err := cmd.Flags().GetBool("no-context")
		if err != nil {
			return err
		}
		if updateContext != "" && clearContext {
			return errors.New("context and no-context can not be used at the same time")
		}
		if updateContext != "" {
			entry.ContextInfo = updateContext
		} else if clearContext {
			entry.ContextInfo = ""
		} else {
			shouldUpdateContext, err := dialogs.Confirm("Do you want to update the context information?")
			if err != nil {
				return err
			}
			if shouldUpdateContext {
				updateContext = console.String("Please enter the new context information")
				entry.ContextInfo = updateContext
			} else if entry.ContextInfo != "" {
				shouldClearContext, err := dialogs.Confirm("Do you want to remove the context information?")
				if err != nil {
					return err
				}
				if shouldClearContext {
					entry.ContextInfo = ""
				}
			}
		}

		// Update the entry in the group or move it
		if newGroup != nil || updateTechnicalName != "" {
			// Here a complex update is necessary,
			// the old entry should be deleted and a new one should be created
			targetGroup := group
			if newGroup != nil {
				targetGroup = *newGroup
			}

			console.Info("The entry needs to be moved")
			console.InfoF(
				"From group '%s' to group '%s'",
				group.ID,
				targetGroup.ID,
			)
			console.InfoF(
				"From technical name '%s' to technical name '%s'",
				originalEntry.TechnicalName,
				entry.TechnicalName,
			)

			// Add the entry in the new group
			createEntryInput := api.CreateEntryInput{
				TechnicalName: entry.TechnicalName,
				BaseTerm:      entry.BaseTerm,
				ContextInfo:   entry.ContextInfo,
			}
			_, err = api.CreateEntry(userConfig, workspaceConfig.ProjectId, targetGroup.ID, createEntryInput)
			if err != nil {
				return err
			}
			console.Success("Creating the new entry was successful, deleting the old one...")

			// Remove the entry from the old group
			err = api.DeleteEntry(userConfig, workspaceConfig.ProjectId, group.ID, originalEntry.TechnicalName)
			if err != nil {
				return err
			}
			console.Success("Deleting the old entry was successful. Done.")
		} else {
			// A simple update
			var input api.UpdateEntryInput
			if originalEntry.BaseTerm != entry.BaseTerm {
				input.BaseTerm = &entry.BaseTerm
			}
			if originalEntry.ContextInfo != entry.ContextInfo {
				input.ContextInfo = &entry.ContextInfo
			}

			resp, err := api.UpdateEntry(
				userConfig,
				workspaceConfig.ProjectId,
				group.ID,
				technicalName,
				input,
			)
			if err != nil {
				return err
			}
			console.SuccessF(
				"Entry '%s' successfully updated in group '%s'",
				entry.TechnicalName,
				resp.DisplayName,
			)
		}

		return nil
	},
}

func init() {
	EntriesCmd.AddCommand(EntriesUpdateCmd)
	EntriesUpdateCmd.Flags().String("group", "", "The path of the group")
	EntriesUpdateCmd.Flags().String("group-id", "", "The id of the group")
	EntriesUpdateCmd.Flags().String("entry", "", "The technical name of the entry")
	EntriesUpdateCmd.Flags().String("new-group", "", "Move the entry to a new group, the group path")
	EntriesUpdateCmd.Flags().String("new-group-id", "", "Move the entry to a new group, the group id")
	EntriesUpdateCmd.Flags().String("technical-name", "", "Set the new technical name")
	EntriesUpdateCmd.Flags().String("base-term", "", "Set the new base term")
	EntriesUpdateCmd.Flags().String("context", "", "Set the new context info")
	EntriesUpdateCmd.Flags().Bool("no-context", false, "Remove the context info")
	EntriesUpdateCmd.Flags().String("working-dir", "", "Set the working directory")
}
