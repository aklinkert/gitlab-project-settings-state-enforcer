package config

import (
	"errors"

	gitlab "github.com/xanzy/go-gitlab"
)

var (
	errFileDoesNotExist                      = errors.New("given config file does not exist")
	errOnlyOneOfBlacklistAndWhitelistAllowed = errors.New("only one of ProjectBlacklist and ProjectWhitelist is allowd, not both")

	errSettingsNameMustBeEmpty      = errors.New("settings.name must be empty")
	errSettingsNamespaceMustBeEmpty = errors.New("settings.namespace must be empty")
)

// Config stores the root group name and some additional configuration values
// settings documented at https://godoc.org/github.com/xanzy/go-gitlab#CreateProjectOptions

type Config struct {
	GroupName        string                     `json:"groupName"`
	Recursive        bool                       `json:"recursive"`
	ProjectBlacklist []string                   `json:"projectBlacklist"`
	ProjectWhitelist []string                   `json:"projectWhitelist"`
	Settings         *gitlab.EditProjectOptions `json:"settings"`
}
