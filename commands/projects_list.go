package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"os"

	"github.com/spf13/cobra"
)

var ProjectsListCmd = &cobra.Command{
	Use:           "list",
	Aliases:       []string{"ls"},
	Short:         "List all your projects",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, _ := os.UserHomeDir()
		userConfig, err := config.LoadUserConfigFromFileRequiresAuth(homeDir)
		if err != nil {
			return err
		}

		// Get current user
		info, err := api.Info(userConfig)
		if err != nil {
			return err
		}

		// Fetch pages of projects
		page := uint32(0)
		for {
			result, err := api.ListProjects(userConfig, 10, page*10)
			if err != nil {
				return err
			}

			for _, project := range result.Result {
				suffix := ""
				if project.Owner == info.ID {
					suffix = " (owner)"
				}

				console.InfoF(
					"[%s] %s%s",
					project.ID,
					project.Name,
					suffix,
				)
			}

			page += 1
			if !result.HasMore {
				break
			}
		}

		return nil
	},
}

func init() {
	ProjectsCmd.AddCommand(ProjectsListCmd)
}
