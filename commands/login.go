package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
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
		config, err := config.LoadUserConfigFromFile(homeDir)
		if err != nil {
			return err
		}

		// Ask for the new api token
		config.ApiToken = console.StringRegex(
			"Please enter your API token",
			`^[a-z0-9_]{1,32}$`,
		)

		// Validate api token
		info, err := api.Info(config)
		if err != nil {
			return err
		}
		console.Info("API token check was successful")

		// Update user config
		err = config.WriteToFile(homeDir)
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
