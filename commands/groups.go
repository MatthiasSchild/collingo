package commands

import "github.com/spf13/cobra"

var GroupsCmd = &cobra.Command{
	Use:     "groups",
	Aliases: []string{"group"},
	Short:   "Manage groups in your project",
	GroupID: "management",
}

func init() {
	RootCmd.AddCommand(GroupsCmd)
}
