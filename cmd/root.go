package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"

	"github.com/Scalify/gitlab-project-settings-state-enforcer/pkg/config"
)

type envCfg struct {
	GitlabEndpoint string `split_words:"true"`
	GitlabToken    string `split_words:"true" required:"true"`
	ConfigFile     string `split_words:"true" default:"./config.json"`
	Verbose        bool
}

var (
	env    = &envCfg{}
	logger = logrus.New()
	cfg    *config.Config
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitlab-project-state-enforcer",
	Short: "Enforces the settings of a bunch of GitLab repos",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := envconfig.Process("", env)
		if err != nil {
			logger.Fatal(err)
		}

		logger.Infof("Loading config file from %v", env.ConfigFile)

		cfg, err = config.Parse(env.ConfigFile)
		if err != nil {
			logger.Fatal(err)
		}

		if env.Verbose {
			logger.SetLevel(logrus.DebugLevel)
		} else {
			logger.SetLevel(logrus.InfoLevel)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
