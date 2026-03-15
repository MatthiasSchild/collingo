package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/partials"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ExportVueI18nCommand = &cobra.Command{
	Use:           "vue-i18n",
	Short:         "Export data in vue-i18n format (for src/i18n.json)",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, _ := os.UserHomeDir()
		userConfig, err := config.LoadUserConfigFromFileRequiresAuth(homeDir)
		if err != nil {
			return err
		}

		workingDir := partials.WorkingDirFromFlags(cmd, "working-dir")
		workspaceConfig, err := config.LoadWorkspaceConfigFromFile(workingDir)
		if err != nil {
			return err
		}

		baseUrl := config.EffectiveServerUrl(userConfig, workspaceConfig)
		content, err := api.ExportVueI18n(userConfig, baseUrl, workspaceConfig.ProjectId)
		if err != nil {
			return err
		}

		fmt.Print(content)
		return nil
	},
}

func init() {
	ExportCommand.AddCommand(ExportVueI18nCommand)
	ExportVueI18nCommand.Flags().String("working-dir", "", "Set the working directory")
}
