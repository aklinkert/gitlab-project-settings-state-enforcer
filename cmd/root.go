package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitlab-repo-state-enforcer",
	Short: "Enforces the settings of a bunch of GitLab repos",
}

func init() {
	rootCmd.PersistentFlags().StringP("gitlab-endpoint", "e", "", "The API endpoint of a GitLab server (env: GITLAB_ENDPOINT)")
	rootCmd.PersistentFlags().StringP("gitlab-token", "t", "", "The API endpoint of a GitLab server (env: GITLAB_TOKEN)")
	rootCmd.PersistentFlags().StringP("config-path", "c", "./config.yaml", "The config to read the actionable settings from (env: CONFIG_PATH)")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
