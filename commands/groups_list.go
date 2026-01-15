package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/partials"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var GroupsListCmd = &cobra.Command{
	Use:           "list",
	Aliases:       []string{"ls"},
	Short:         "List the groups of a project",
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

		// List groups
		summary, err := api.ListGroupSummary(userConfig, workspaceConfig.ProjectId)
		if err != nil {
			return err
		}

		for _, entry := range summary.Result {
			name := strings.Join(append(entry.BreadcrumbNames, entry.DisplayName), " > ")
			console.InfoF(
				"[%s] %s (%s)",
				entry.ID,
				name,
				entry.TechnicalName,
			)
		}

		return nil
	},
}

func init() {
	GroupsCmd.AddCommand(GroupsListCmd)
	GroupsListCmd.Flags().String("working-dir", "", "Set the working directory")
}
