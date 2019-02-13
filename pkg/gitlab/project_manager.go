package gitlab

import (
	"fmt"
	"net/http"

	"github.com/xanzy/go-gitlab"

	"github.com/Scalify/gitlab-project-settings-state-enforcer/pkg/config"
	"github.com/Scalify/gitlab-project-settings-state-enforcer/pkg/internal/stringslice"

	"github.com/Sirupsen/logrus"
)

// ProjectManager fetches a list of repositories from GitLab
type ProjectManager struct {
	logger                  *logrus.Entry
	groupsClient            groupsClient
	projectsClient          projectsClient
	protectedBranchesClient protectedBranchesClient
	branchesClient          branchesClient
	config                  *config.Config
}

// NewProjectManager returns a new ProjectManager instance
func NewProjectManager(
	logger *logrus.Entry,
	groupsClient groupsClient,
	projectsClient projectsClient,
	protectedBranchesClient protectedBranchesClient,
	branchesClient branchesClient,
	config *config.Config,
) *ProjectManager {

	return &ProjectManager{
		logger:                  logger,
		groupsClient:            groupsClient,
		projectsClient:          projectsClient,
		protectedBranchesClient: protectedBranchesClient,
		branchesClient:          branchesClient,
		config:                  config,
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

// EnsureBranchesAndProtection ensures that 1) the default branch exists and 2) all of the protected branches are configured correctly
func (m *ProjectManager) EnsureBranchesAndProtection(project Project) error {
	if err := m.ensureDefaultBranch(project); err != nil {
		return err
	}

	for _, b := range m.config.ProtectedBranches {
		if resp, err := m.protectedBranchesClient.UnprotectRepositoryBranches(project.ID, b.Name); err != nil && resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("failed to unprotect branch %v befor protection: %v", b.Name, err)
		}

		opt := &gitlab.ProtectRepositoryBranchesOptions{
			Name:             gitlab.String(b.Name),
			PushAccessLevel:  b.PushAccessLevel.Value(),
			MergeAccessLevel: b.MergeAccessLevel.Value(),
		}

		if _, _, err := m.protectedBranchesClient.ProtectRepositoryBranches(project.ID, opt); err != nil {
			return fmt.Errorf("failed to protect branch %s: %v", b.Name, err)
		}
	}

	return nil
}

func (m *ProjectManager) ensureDefaultBranch(project Project) error {
	if !m.config.CreateDefaultBranch ||
		m.config.Settings.DefaultBranch == nil ||
		*m.config.Settings.DefaultBranch == "master" {
		return nil
	}

	opt := &gitlab.CreateBranchOptions{
		Branch: m.config.Settings.DefaultBranch,
		Ref:    gitlab.String("master"),
	}

	m.logger.Debugf("Ensuring default branch %s existence ... ", *opt.Branch)

	_, resp, err := m.branchesClient.GetBranch(project.ID, *opt.Branch)
	if err == nil {
		m.logger.Debugf("Ensuring default branch %s existence ... already exists!", *opt.Branch)
		return nil
	}

	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to check for default branch existence, got unexpected response status code %d", resp.StatusCode)
	}

	if _, _, err := m.branchesClient.CreateBranch(project.ID, opt); err != nil {
		return fmt.Errorf("failed to create default branch %s: %v", *opt.Branch, err)
	}

	return nil
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
