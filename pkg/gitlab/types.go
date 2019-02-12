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

var (
	listGroupProjectOps = &gitlab.ListGroupProjectsOptions{
		MinAccessLevel: gitlab.AccessLevel(gitlab.DeveloperPermissions),
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
