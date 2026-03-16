package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/partials"
	"os"

	"github.com/spf13/cobra"
)

var ExportFlutterCommand = &cobra.Command{
	Use:           "flutter",
	Short:         "Export data for Flutter as ARB format",
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
		err = api.ExportFlutter(userConfig, baseUrl, workspaceConfig.ProjectId, directory, format)
		if err != nil {
			return err
		}

		console.Success("Successfully exported Flutter ARB files to " + directory)
		return nil
	},
}

func init() {
	ExportFlutterCommand.Flags().StringP("directory", "d", "", "Directory where ARB files should be extracted (required)")
	ExportFlutterCommand.Flags().Bool("format", false, "Format the output")
	ExportFlutterCommand.MarkFlagRequired("directory")
	ExportCommand.AddCommand(ExportFlutterCommand)
}
