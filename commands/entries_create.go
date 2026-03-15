package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/partials"
	"os"

	"github.com/spf13/cobra"
)

var EntriesCreateCmd = &cobra.Command{
	Use:           "create",
	Aliases:       []string{"add"},
	Short:         "Create an entry",
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

		// Get the technical name
		technicalName, err := cmd.Flags().GetString("technical-name")
		if err != nil {
			return err
		} else if technicalName == "" {
			technicalName = console.TechnicalName("Please enter a technical name for the entry")
		}

		// Get the base term
		baseTerm, err := cmd.Flags().GetString("base-term")
		if err != nil {
			return err
		} else if baseTerm == "" {
			baseTerm = console.StringRequired("Please enter a base term")
		}

		// Get the context info
		contextInfo, err := cmd.Flags().GetString("context-info")
		if err != nil {
			return err
		} else if contextInfo == "" {
			contextInfo = console.String("Enter optionally a context info")
		}

		group, err = api.CreateEntry(userConfig, baseUrl, workspaceConfig.ProjectId, group.ID, api.CreateEntryInput{
			TechnicalName: technicalName,
			BaseTerm:      baseTerm,
			ContextInfo:   contextInfo,
		})
		if err != nil {
			return err
		}

		console.SuccessF(
			"Entry '%s' successfully created in group '%s'",
			technicalName,
			group.ID,
		)
		return nil
	},
}

func init() {
	EntriesCmd.AddCommand(EntriesCreateCmd)
	EntriesCreateCmd.Flags().String("working-dir", "", "Set the working directory")
	EntriesCreateCmd.Flags().String("group", "", "group path (e.g. \"base.footer\")")
	EntriesCreateCmd.Flags().String("group-id", "", "The group where the entry should be added")
	EntriesCreateCmd.Flags().String("technical-name", "", "The technical name for this entry")
	EntriesCreateCmd.Flags().String("base-term", "", "The term in the base language")
	EntriesCreateCmd.Flags().String("context-info", "", "Optional additional context information about the usage of the term")
}
