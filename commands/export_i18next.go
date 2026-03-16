package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/partials"
	"os"

	"github.com/spf13/cobra"
)

var ExportI18nextCommand = &cobra.Command{
	Use:           "i18next",
	Short:         "Export data for i18next",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		directory, _ := cmd.Flags().GetString("directory")
		if directory == "" {
			return cmd.Help()
		}

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
		err = api.ExportI18next(userConfig, baseUrl, workspaceConfig.ProjectId, directory, format)
		if err != nil {
			return err
		}

		console.Success("Successfully exported i18next translations to " + directory)
		return nil
	},
}

func init() {
	ExportI18nextCommand.Flags().StringP("directory", "d", "", "Directory where translation files should be extracted (required)")
	ExportI18nextCommand.Flags().Bool("format", false, "Format the output")
	ExportI18nextCommand.MarkFlagRequired("directory")
	ExportCommand.AddCommand(ExportI18nextCommand)
}
