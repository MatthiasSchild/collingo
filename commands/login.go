package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/partials"
	"os"

	"github.com/spf13/cobra"
)

var LoginCmd = &cobra.Command{
	Use:           "login",
	Short:         "Login to collingo",
	Long:          "Use an API key to login this device to collingo",
	GroupID:       "setup",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load user config
		homeDir, _ := os.UserHomeDir()
		userConfig, err := config.LoadUserConfigFromFile(homeDir)
		if err != nil {
			return err
		}

		// Ask for the new api token
		userConfig.ApiToken = console.StringRegex(
			"Please enter your API token",
			`^[a-z0-9_]{1,32}$`,
		)

		workingDir := partials.WorkingDirFromFlags(cmd, "working-dir")
		workspaceConfig, _ := config.LoadWorkspaceConfigFromFile(workingDir)
		baseUrl := config.EffectiveServerUrl(userConfig, workspaceConfig)

		// Validate api token
		info, err := api.Info(userConfig, baseUrl)
		if err != nil {
			return err
		}
		console.Info("API token check was successful")

		// Update user config
		err = userConfig.WriteToFile(homeDir)
		if err != nil {
			return err
		}

		console.InfoF("Hello, %s!", info.FirstName)
		console.Success("API token has been successfully updated in the user settings")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(LoginCmd)
}
