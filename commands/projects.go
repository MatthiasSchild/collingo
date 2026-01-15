package commands

import "github.com/spf13/cobra"

var ProjectsCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"project"},
	Short:   "Manage your projects",
	GroupID: "management",
}

func init() {
	RootCmd.AddCommand(ProjectsCmd)
}
