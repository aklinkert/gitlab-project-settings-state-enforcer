package gitlab

import (
	"fmt"

	"github.com/Scalify/gitlab-project-settings-state-enforcer/pkg/config"
	"github.com/Scalify/gitlab-project-settings-state-enforcer/pkg/internal/stringslice"

	"github.com/Sirupsen/logrus"
)

// ProjectManager fetches a list of repositories from GitLab
type ProjectManager struct {
	logger         *logrus.Entry
	groupsClient   groupsClient
	projectsClient projectsClient
	config         *config.Config
}

// NewProjectManager returns a new ProjectManager instance
func NewProjectManager(logger *logrus.Entry, groupsClient groupsClient, projectsClient projectsClient, config *config.Config) *ProjectManager {
	return &ProjectManager{
		logger:         logger,
		groupsClient:   groupsClient,
		projectsClient: projectsClient,
		config:         config,
	}
}

// GetProjects fetches a list of accessible repos within the groups set in config file
func (m *ProjectManager) GetProjects() ([]Project, error) {
	var repos []Project

	m.logger.Debugf("Fetching gitlab repos for group %s ...", m.config.GroupName)

	projects, _, err := m.groupsClient.ListGroupProjects(m.config.GroupName, listGroupProjectOps, addIncludeSubgroups)
	if err != nil {
		return []Project{}, fmt.Errorf("failed to fetch GitLab projects or group %q: %v", m.config.GroupName, err)
	}

	for _, p := range projects {
		if len(m.config.ProjectWhitelist) > 0 && !stringslice.Contains(p.PathWithNamespace, m.config.ProjectWhitelist) {
			m.logger.Debugf("Skipping repo %s as it's not whitelisted", p.PathWithNamespace)
			continue
		}
		if stringslice.Contains(p.PathWithNamespace, m.config.ProjectBlacklist) {
			m.logger.Debugf("Skipping repo %s as it's blacklisted", p.PathWithNamespace)
			continue
		}

		repos = append(repos, Project{
			ID:       p.ID,
			Name:     p.Name,
			FullPath: p.PathWithNamespace,
		})
	}

	m.logger.Debugf("Fetching gitlab repos done. Got %d repos.", len(repos))

	return repos, nil
}

// UpdateSettings updates the project settings on gitlab
func (m *ProjectManager) UpdateSettings(project Project) error {
	m.logger.Debugf("Updating settings of project %s ...", project.FullPath)

	if _, _, err := m.projectsClient.EditProject(project.ID, m.config.Settings); err != nil {
		return fmt.Errorf("failed to update settings or project %s: %v", project.FullPath, err)
	}

	m.logger.Debugf("Updating settings of project %s done.", project.FullPath)

	return nil
}
