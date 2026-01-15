package commands

import "github.com/spf13/cobra"

var EntriesCmd = &cobra.Command{
	Use:     "entries",
	Aliases: []string{"entry"},
	Short:   "Manage entries in your project",
	GroupID: "management",
}

func init() {
	RootCmd.AddCommand(EntriesCmd)
}
