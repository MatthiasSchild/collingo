package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/partials"
	"os"

	"github.com/spf13/cobra"
)

var GroupsCreateCmd = &cobra.Command{
	Use:           "create",
	Aliases:       []string{"add"},
	Short:         "Create a group",
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

		// Get the display name
		displayName, err := cmd.Flags().GetString("display-name")
		if err != nil {
			return err
		} else if displayName == "" {
			displayName = console.StringRequired("Please enter a readable name for the group")
		}

		// Get the technical name
		technicalName, err := cmd.Flags().GetString("technical-name")
		if err != nil {
			return err
		} else if technicalName == "" {
			technicalName = console.TechnicalName("Please enter a technical name for the group")
		}

		// Select a parent group
		parent, err := partials.GetParentFromCommand(userConfig, workspaceConfig, cmd)
		if err != nil {
			return err
		}
		var parentId string
		if parent.ID != "" {
			parentId = parent.ID
		}

		// Create the group
		group, err := api.CreateGroup(userConfig, workspaceConfig.ProjectId, api.CreateGroupInput{
			DisplayName:   displayName,
			TechnicalName: technicalName,
			Parent:        parentId,
		})
		if err != nil {
			return err
		}

		console.SuccessF(
			"Group '%s' successfully created with ID '%s'",
			group.DisplayName,
			group.ID,
		)
		return nil
	},
}

func init() {
	GroupsCmd.AddCommand(GroupsCreateCmd)
	GroupsCreateCmd.Flags().String("parent", "", "The path of the parent group")
	GroupsCreateCmd.Flags().String("parent-id", "", "The ID of the parent group")
	GroupsCreateCmd.Flags().String("display-name", "", "A human readable name for the group")
	GroupsCreateCmd.Flags().String("technical-name", "", "A technical identifier for this group")
	GroupsCreateCmd.Flags().String("working-dir", "", "Set the working directory")
}
