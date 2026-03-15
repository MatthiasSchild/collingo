package commands

import (
	"collingo/config"
	"collingo/console"
	"collingo/partials"
	"os"

	"github.com/spf13/cobra"
)

var EntriesListCmd = &cobra.Command{
	Use:           "list",
	Aliases:       []string{"ls"},
	Short:         "List all entries of a group",
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

		if len(group.Entries) == 0 {
			console.Info("Group has no entries")
			return nil
		}

		// List entries
		for _, entry := range group.Entries {
			baseTerm := entry.BaseTerm
			if len(baseTerm) > 30 {
				baseTerm = baseTerm[0:27] + "..."
			}
			console.InfoF(
				"[%s] | %s",
				entry.TechnicalName,
				baseTerm,
			)
			if len(entry.ContextInfo) > 0 {
				console.InfoF(" -> Context: %s", entry.ContextInfo)
			}
		}

		return nil
	},
}

func init() {
	EntriesCmd.AddCommand(EntriesListCmd)
	EntriesListCmd.Flags().String("group", "", "The path of the group")
	EntriesListCmd.Flags().String("group-id", "", "The id of the group")
}
