package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	gl "github.com/Scalify/gitlab-project-settings-state-enforcer/pkg/gitlab"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync gitlab's project settings with the config",
	Run: func(cmd *cobra.Command, args []string) {
		client := gitlab.NewClient(nil, env.GitlabToken)
		if env.GitlabEndpoint != "" {
			if err := client.SetBaseURL(env.GitlabEndpoint); err != nil {
				logger.Fatal(err)
			}
		}

		manager := gl.NewProjectManager(
			logger.WithField("module", "project_manager"),
			client,
			cfg,
		)

		projects, err := manager.GetProjects()
		if err != nil {
			logger.Fatal(err)
		}

		for _, project := range projects {
			logger.Infof("Updating project %s", project.FullPath)

			if err := manager.EnsureBranchesAndProtection(project); err != nil {
				logger.Errorf("failed to ensure branches of repo %v: %v", project.FullPath, err)
			}

			if err := manager.UpdateSettings(project); err != nil {
				logger.Errorf("failed to update settings of repo %v: %v", project.FullPath, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
