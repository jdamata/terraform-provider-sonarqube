# sonarqube_alm_gitlab

Provides a Sonarqube GitLab Alm/Devops Platform Integration resource. This can be used to create and manage a Alm/Devops
Platform Integration for GitLab.

## Example: Create a GitHub Alm Integration

```terraform
resource "sonarqube_alm_gitlab" "gitlab-alm" {
  key                    = "myalm"
  personal_access_token  = "my_personal_access_token"
  url                    = "https://gitlab.com/api/v4"
}
```

## Argument Reference

The following arguments are supported:

- key - (Required) - Unique key of the GitLab instance setting. Maximum length: 200
- personal_access_token - (Required) - GitLab App personal access token with the `read_api` scope. See [this doc](https://docs.sonarqube.org/latest/devops-platform-integration/gitlab-integration/#importing-your-gitlab-projects-into-sonarqube) for more information. Maximum length: 2000
- url - (Required) - GitLab API URL. Maximum length: 2000

## Attributes Reference

The following attributes are exported:

- key - The unique key of the GitLab instance setting.
- url - GitLab API URL.

