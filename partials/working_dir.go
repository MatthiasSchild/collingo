package partials

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func WorkingDirFromFlags(cmd *cobra.Command, flagName string) string {
	workingDir, err := cmd.Flags().GetString("working-dir")
	if err != nil {
		slog.Warn("Could not parse working dir flag", "err", err.Error())
	} else if workingDir != "" {
		return workingDir
	}

	workingDir, _ = os.Getwd()
	return workingDir
}
