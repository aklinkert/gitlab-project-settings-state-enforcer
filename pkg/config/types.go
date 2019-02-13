package config

import (
	"errors"

	"github.com/xanzy/go-gitlab"
)

const (
	AccessLevelDeveloper  = "developer"
	AccessLevelMaintainer = "maintainer"
)

var (
	errFileDoesNotExist                      = errors.New("given config file does not exist")
	errOnlyOneOfBlacklistAndWhitelistAllowed = errors.New("only one of ProjectBlacklist and ProjectWhitelist is allowed, not both")

	errSettingsNameMustBeEmpty      = errors.New("settings.name must be empty")
	errSettingsNamespaceMustBeEmpty = errors.New("settings.namespace must be empty")
)

// Config stores the root group name and some additional configuration values
// settings documented at https://godoc.org/github.com/xanzy/go-gitlab#CreateProjectOptions
type Config struct {
	GroupName           string                     `json:"group_name"`
	Recursive           bool                       `json:"recursive"`
	ProjectBlacklist    []string                   `json:"project_blacklist"`
	ProjectWhitelist    []string                   `json:"project_whitelist"`
	CreateDefaultBranch bool                       `json:"create_default_branch"`
	ProtectedBranches   []ProtectedBranch          `json:"protected_branches"`
	Settings            *gitlab.EditProjectOptions `json:"settings"`
}

// ProtectedBranch defines who can act on a protected branch
type ProtectedBranch struct {
	Name             string      `json:"name"`
	PushAccessLevel  AccessLevel `json:"pushAccessLevel"`
	MergeAccessLevel AccessLevel `json:"mergeAccessLevel"`
}

type AccessLevel string

func (a AccessLevel) Value() *gitlab.AccessLevelValue {
	switch a {
	case AccessLevelDeveloper:
		return gitlab.AccessLevel(gitlab.DeveloperPermissions)
	case AccessLevelMaintainer:
		return gitlab.AccessLevel(gitlab.MaintainerPermissions)
	default:
		return gitlab.AccessLevel(gitlab.NoPermissions)
	}
}
