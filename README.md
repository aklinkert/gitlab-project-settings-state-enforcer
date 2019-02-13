# gitlab-project-settings-state-enforcer
> Enforces the settings of a bunch of gitlab repos by reading a config file and talking to the GitLab API.

# Usage

This project get's built automatically and is usable as a docker image at hub.docker.com: [`scalify/gitlab-project-settings-state-enforcer`](https://hub.docker.com/r/scalify/gitlab-project-settings-state-enforcer)

# Configuration

Configuration of project interaction is currently possible via JSON files providing a Config object. The config object has the following fields:


| field                   | type              | required | content                                                                                                          | default |
|-------------------------|-------------------|----------|------------------------------------------------------------------------------------------------------------------|---------|
| `group_name`            | string            | yes      | The path of the root group (e.g. `example` or `some/nested/example`)                                             |         |
| `project_blacklist`     | []string          | no       | A list of projects to blacklist (cannot be set when project_whitelist is used)                                   | []      |
| `project_whitelist`     | []string          | no       | A list of projects to whitelist (cannot be set when project_blacklist is used)                                   | []      |
| `create_default_branch` | bool              | no       | Whether the default branch configured in `settings.default_branch` should be created if it doesn't exist already |         |
| `protected_branches`    | []ProtectedBranch | no       | A list of branches to protect, together with the infos which roles are allowed to merge or push.                 |         |
| settings                | Object            | yes      | the gitlab settings to change. Possible keys: https://docs.gitlab.com/ce/api/projects.html#edit-project          |         |


`ProtectedBranch` 

| field                | type   | required | content                                                                              |
|----------------------|--------|----------|--------------------------------------------------------------------------------------|
| `name`               | string | yes      | he name of the branch to protect                                                     |
| `push_access_level`  | string | yes      | Which role is allowed to push (possible values: `maintainer`, `developer`, `noone`)  |
| `merge_access_level` | string | yes      | Which role is allowed to merge (possible values: `maintainer`, `developer`, `noone`) |


## Env vars

To control the GitLab API endpoint and the authentication as well as further internal flags please use the following env vars:

| Name              | Required | Description                                                                       | Default      |
|-------------------|----------|-----------------------------------------------------------------------------------|--------------|
| `VERBOSE`         | no       | Enables debug logging when enabled                                                | `false`      |
| `GITLAB_TOKEN`    | yes      | The GitLab API token used for authentication                                      |              |
| `GITLAB_ENDPOINT` | no       | Only override when using GitLab on premise, set this to your GitLab Server Domain | (gitlab.com) |


## Config Example
An example config might look like the following:

```json
{
  "group_name": "example",
  "project_blacklist": [
    "example/path-to/ignored-project"
  ],
  "project_whitelist": [],
  "create_default_branch": true,
  "protected_branches": [
    {"name":  "develop", "push_access_level":  "maintainer", "merge_access_level":  "developer"},
    {"name":  "master", "push_access_level":  "maintainer", "merge_access_level":  "developer"}
  ],
  "settings": {
    "default_branch": "develop",
    "issues_enabled": true,
    "merge_requests_enabled": true,
    "jobs_enabled": true,
    "wiki_enabled": false,
    "snippets_enabled": false,
    "resolve_outdated_diff_discussions": true,
    "container_registry_enabled": true,
    "shared_runners_enabled": false,
    "only_allow_merge_if_pipeline_succeeds": false,
    "only_allow_merge_if_all_discussions_are_resolved": true,
    "merge_method": "merge",
    "public_builds": false,
    "lfs_enabled": true,
    "request_access_enabled": false,
    "tag_list": [],
    "printing_merge_request_link_enabled": true,
    "ci_config_path": null,
    "approvals_before_merge": 1
  }
}
```


# License

    MIT License
    
    Copyright (c) 2019 Scalify GmbH
