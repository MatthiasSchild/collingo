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
	Short:         "Export data for vue-i18n",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, _ := os.UserHomeDir()
		userConfig, err := config.LoadUserConfigFromFileRequiresAuth(homeDir)
		if err != nil {
			return err
		}

		workingDir := partials.WorkingDirFromFlags(cmd)
		workspaceConfig, err := config.LoadWorkspaceConfigFromFile(workingDir)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetBool("format")

		baseUrl := config.EffectiveServerUrl(userConfig, workspaceConfig)
		content, err := api.ExportVueI18n(userConfig, baseUrl, workspaceConfig.ProjectId, format)
		if err != nil {
			return err
		}

		fmt.Print(content)
		return nil
	},
}

func init() {
	ExportVueI18nCommand.Flags().Bool("format", false, "Format the output")
	ExportCommand.AddCommand(ExportVueI18nCommand)
}
