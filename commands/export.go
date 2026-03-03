package commands

import "github.com/spf13/cobra"

var ExportCommand = &cobra.Command{
	Use:           "export",
	Short:         "Export data to a file",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	RootCmd.AddCommand(ExportCommand)
}
