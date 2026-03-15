package partials

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

const WorkingDirFlagName = "working-dir"

func WorkingDirFromFlags(cmd *cobra.Command) string {
	workingDir, err := cmd.Flags().GetString(WorkingDirFlagName)
	if err != nil {
		slog.Warn("Could not parse working dir flag", "err", err.Error())
	} else if workingDir != "" {
		return workingDir
	}

	workingDir, _ = os.Getwd()
	return workingDir
}
