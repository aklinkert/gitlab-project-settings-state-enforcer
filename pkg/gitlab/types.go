package gitlab

import (
	"net/http"

	gitlab "github.com/xanzy/go-gitlab"
)

// Project holds all relevant infos for a GitLab repo
type Project struct {
	ID       int
	Name     string
	FullPath string
}

type groupsClient interface {
	ListGroupProjects(gid interface{}, opt *gitlab.ListGroupProjectsOptions, options ...gitlab.OptionFunc) ([]*gitlab.Project, *gitlab.Response, error)
}

type projectsClient interface {
	EditProject(pid interface{}, opt *gitlab.EditProjectOptions, options ...gitlab.OptionFunc) (*gitlab.Project, *gitlab.Response, error)
}

type protectedBranchesClient interface {
	ProtectRepositoryBranches(pid interface{}, opt *gitlab.ProtectRepositoryBranchesOptions, options ...gitlab.OptionFunc) (*gitlab.ProtectedBranch, *gitlab.Response, error)
	UnprotectRepositoryBranches(pid interface{}, branch string, options ...gitlab.OptionFunc) (*gitlab.Response, error)
	// ListProtectedBranches(pid interface{}, opt *gitlab.ListProtectedBranchesOptions, options ...gitlab.OptionFunc) ([]*gitlab.ProtectedBranch, *gitlab.Response, error)
}

type branchesClient interface {
	CreateBranch(pid interface{}, opt *gitlab.CreateBranchOptions, options ...gitlab.OptionFunc) (*gitlab.Branch, *gitlab.Response, error)
	GetBranch(pid interface{}, branch string, options ...gitlab.OptionFunc) (*gitlab.Branch, *gitlab.Response, error)
}

var (
	listGroupProjectOps = &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
	}

	addIncludeSubgroups = gitlab.OptionFunc(func(req *http.Request) error {
		v := req.URL.Query()
		v.Add("include_subgroups", "true")
		req.URL.RawQuery = v.Encode()
		return nil
	})
)
