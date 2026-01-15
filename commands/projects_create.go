package commands

import (
	"collingo/api"
	"collingo/config"
	"collingo/console"
	"collingo/dialogs"
	"os"

	"github.com/spf13/cobra"
)

var ProjectsCreateCmd = &cobra.Command{
	Use:           "create",
	Aliases:       []string{"add"},
	Short:         "Create a project",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, _ := os.UserHomeDir()
		userConfig, err := config.LoadUserConfigFromFileRequiresAuth(homeDir)
		if err != nil {
			return err
		}

		// Get name for project
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		} else if name == "" {
			name = console.StringRequired("Please enter a name for the new project")
		}

		// Select the base language
		baseLanguage, err := dialogs.LanguageSelection("Please select a base language (e.g. en-US)")
		if err != nil {
			return err
		}

		// Select the translation languages
		translationLanguages, err := dialogs.MultiLanguageSelection(
			"Add translation languages to the project",
		)
		if err != nil {
			return err
		}

		// Create the project
		project, err := api.CreateProject(userConfig, api.CreateProjectInput{
			Name:                 name,
			BaseLanguage:         baseLanguage,
			TranslationLanguages: translationLanguages,
		})
		if err != nil {
			return err
		}

		console.SuccessF(
			"Project '%s' successfully created with ID '%s'",
			project.Name,
			project.ID,
		)
		return nil
	},
}

func init() {
	ProjectsCmd.AddCommand(ProjectsCreateCmd)
	ProjectsCreateCmd.Flags().String("name", "", "The name of the new project")
}
