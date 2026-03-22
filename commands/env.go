package commands

import (
	"collingo/config"
	"collingo/console"
	"collingo/partials"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var EnvCmd = &cobra.Command{
	Use:           "env",
	Short:         "Show current environment and configuration",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, _ := os.UserHomeDir()
		userConfig, err := config.LoadUserConfigFromFile(homeDir)
		if err != nil {
			return err
		}

		workingDir := partials.WorkingDirFromFlags(cmd)
		workspaceConfig, workspaceConfigPath, err := config.LoadWorkspaceConfigFromFileWithPath(workingDir)
		if err != nil {
			// Not in a workspace
			workspaceConfig = nil
		}

		effectiveServerUrl := config.EffectiveServerUrl(userConfig, workspaceConfig)

		// Workspace
		if workspaceConfig != nil {
			console.Info("Workspace: yes")
			console.InfoF("  Config file: %s", workspaceConfigPath)
			console.InfoF("  Project ID: %s", workspaceConfig.ProjectId)
			if workspaceConfig.Template != nil && workspaceConfig.Template.HasKind() {
				lines := strings.Split(workspaceConfig.Template.EnvDescription(), "\n")
				for i, line := range lines {
					if i == 0 {
						console.InfoF("  Template: %s", line)
					} else {
						console.InfoF("  %s", line)
					}
				}
			} else {
				console.Info("  Template: (none)")
			}
		} else {
			console.Info("Workspace: no")
		}

		// Logged in
		if userConfig.ApiToken != "" {
			console.Info("Logged in: yes")
		} else {
			console.Info("Logged in: no")
		}

		// Server URL
		console.InfoF("Server URL: %s", effectiveServerUrl)

		return nil
	},
}

func init() {
	RootCmd.AddCommand(EnvCmd)
}
