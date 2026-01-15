package partials

import (
	"collingo/api"
	"collingo/config"
	"collingo/dialogs"
	"collingo/models"
	"collingo/utils"
	"errors"

	"github.com/spf13/cobra"
)

func GetGroupFromCommand(
	userConfig *config.UserConfig,
	workspaceConfig *config.WorkspaceConfig,
	cmd *cobra.Command,
) (group models.GroupModel, err error) {
	var groupId string
	var groupPath string

	// Get the flags
	groupId, err = cmd.Flags().GetString("group-id")
	if err != nil {
		return
	}
	groupPath, err = cmd.Flags().GetString("group")
	if err != nil {
		return
	}

	if groupId != "" && groupPath != "" {
		err = errors.New("group and group-id can not be used at the same time")
		return
	}

	// Check if group-id is set in the flags
	if groupId != "" {
		if !utils.IsObjectID(groupId) {
			err = errors.New("The id is not a valid object id")
			return
		}

		// Group ID set, get the group by group id
		group, err = api.GetGroup(userConfig, workspaceConfig.ProjectId, groupId)
		return
	}

	// group-id is not set, try to get the group by path
	if groupPath != "" {
		// Try to resolve the group path and return the result
		group, err = ResolveGroupPath(userConfig, workspaceConfig.ProjectId, groupPath)
		return
	}

	// Neither group-id nor group is set, open the selection dialog
	var groupSummary models.GroupSummaryModel
	groupSummary, err = dialogs.GroupSelection(userConfig, workspaceConfig.ProjectId)
	if err != nil {
		return
	}
	group, err = api.GetGroup(userConfig, workspaceConfig.ProjectId, groupSummary.ID)
	return
}

func GetParentFromCommand(
	userConfig *config.UserConfig,
	workspaceConfig *config.WorkspaceConfig,
	cmd *cobra.Command,
) (group models.GroupModel, err error) {
	var groupId string
	var groupPath string

	// Get the flags
	groupId, err = cmd.Flags().GetString("parent-id")
	if err != nil {
		return
	}
	groupPath, err = cmd.Flags().GetString("parent")
	if err != nil {
		return
	}

	if groupId != "" && groupPath != "" {
		err = errors.New("parent and parent-id can not be used at the same time")
		return
	}

	// Check if group-id is set in the flags
	if groupId != "" {
		if !utils.IsObjectID(groupId) {
			err = errors.New("The id is not a valid object id")
			return
		}

		// Group ID set, get the group by group id
		group, err = api.GetGroup(userConfig, workspaceConfig.ProjectId, groupId)
		return
	}

	// group-id is not set, try to get the group by path
	if groupPath != "" {
		// Try to resolve the group path and return the result
		group, err = ResolveGroupPath(userConfig, workspaceConfig.ProjectId, groupPath)
		return
	}

	// Neither group-id nor group is set, open the selection dialog
	selectParent, err := dialogs.Confirm("Do you want to select a parent group?")
	if err != nil {
		return group, err
	}
	if selectParent {
		var groupSummary models.GroupSummaryModel
		groupSummary, err = dialogs.GroupSelection(userConfig, workspaceConfig.ProjectId)
		if err != nil {
			return
		}
		group, err = api.GetGroup(userConfig, workspaceConfig.ProjectId, groupSummary.ID)
	}
	return
}
