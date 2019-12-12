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
