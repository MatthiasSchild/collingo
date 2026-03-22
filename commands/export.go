package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/partials"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var ExportCommand = &cobra.Command{
	Use:   "export",
	Short: "Export data to a file",
	Long: "Export using the workspace template in .collingo.json. " +
		"Use subcommands (e.g. export flutter) for an explicit format and paths.",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runWorkspaceExport,
}

func runWorkspaceExport(cmd *cobra.Command, args []string) error {
	homeDir, _ := os.UserHomeDir()
	userConfig, err := config.LoadUserConfigFromFileRequiresAuth(homeDir)
	if err != nil {
		return err
	}

	workingDir := partials.WorkingDirFromFlags(cmd)
	workspaceConfig, cfgPath, err := config.LoadWorkspaceConfigFromFileWithPath(workingDir)
	if err != nil {
		return err
	}

	tmpl := workspaceConfig.Template
	if err := tmpl.ValidateForDefaultExport(); err != nil {
		return err
	}

	workspaceRoot := filepath.Dir(cfgPath)
	baseURL := config.EffectiveServerUrl(userConfig, workspaceConfig)
	projectID := workspaceConfig.ProjectId

	switch tmpl.Kind {
	case config.TemplateKindFlutter:
		dir, err := config.ResolvePathUnderWorkspace(workspaceRoot, tmpl.Directory)
		if err != nil {
			return err
		}
		if err := api.ExportFlutter(userConfig, baseURL, projectID, dir, tmpl.Formatted); err != nil {
			return err
		}
		console.Success("Successfully exported Flutter ARB files to " + dir)
	case config.TemplateKindI18next:
		dir, err := config.ResolvePathUnderWorkspace(workspaceRoot, tmpl.Directory)
		if err != nil {
			return err
		}
		if err := api.ExportI18next(userConfig, baseURL, projectID, dir, tmpl.Formatted); err != nil {
			return err
		}
		console.Success("Successfully exported i18next translations to " + dir)
	case config.TemplateKindVueI18n:
		content, err := api.ExportVueI18n(userConfig, baseURL, projectID, tmpl.Formatted)
		if err != nil {
			return err
		}
		if tmpl.OutputFile == "" {
			fmt.Print(content)
			return nil
		}
		outPath, err := config.ResolvePathUnderWorkspace(workspaceRoot, tmpl.OutputFile)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
			return err
		}
		console.Success("Successfully exported vue-i18n data to " + outPath)
	case config.TemplateKindJSON:
		content, err := api.ExportJson(userConfig, baseURL, projectID, tmpl.Formatted)
		if err != nil {
			return err
		}
		if tmpl.OutputFile == "" {
			fmt.Print(content)
			return nil
		}
		outPath, err := config.ResolvePathUnderWorkspace(workspaceRoot, tmpl.OutputFile)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
			return err
		}
		console.Success("Successfully exported JSON to " + outPath)
	default:
		return fmt.Errorf("%w: %q", config.ErrTemplateExportUnknownKind, tmpl.Kind)
	}
	return nil
}

func init() {
	RootCmd.AddCommand(ExportCommand)
}
