package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "collingo",
}

func init() {
	RootCmd.AddGroup(
		&cobra.Group{ID: "setup", Title: "Setup commands"},
		&cobra.Group{ID: "management", Title: "Management commands"},
	)
}
