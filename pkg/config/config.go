package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Parse takes the given configFilePath and reads the containing config file into a config struct
func Parse(configFilePath string) (*Config, error) {
	var err error
	configFilePath, err = filepath.Abs(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error resolving file path  %v", err)
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, errFileDoesNotExist
	} else if err != nil {
		return nil, fmt.Errorf("error checking config file: %v", err)
	}

	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %v", configFilePath, err)
	}

	cfg := &Config{
		ProjectBlacklist: make([]string, 0),
		ProjectWhitelist: make([]string, 0),
	}
	if err := json.Unmarshal(b, cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %q: %v", configFilePath, err)
	}

	return checkConfig(cfg)
}

func checkConfig(cfg *Config) (*Config, error) {

	if len(cfg.ProjectBlacklist) > 0 && len(cfg.ProjectWhitelist) > 0 {
		return nil, errOnlyOneOfBlacklistAndWhitelistAllowed
	}

	if cfg.Settings.Name != nil {
		return nil, errSettingsNameMustBeEmpty
	}

	if cfg.Settings.NamespaceID != nil {
		return nil, errSettingsNamespaceMustBeEmpty
	}

	return cfg, nil
}
