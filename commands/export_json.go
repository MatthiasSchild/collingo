package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/partials"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ExportJsonCommand = &cobra.Command{
	Use:           "json",
	Short:         "Export data to a json file",
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
		jsonContent, err := api.ExportJson(userConfig, baseUrl, workspaceConfig.ProjectId)
		if err != nil {
			return err
		}

		fmt.Print(jsonContent)
		return nil
	},
}

func init() {
	ExportCommand.AddCommand(ExportJsonCommand)
}
