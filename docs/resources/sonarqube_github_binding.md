# sonarqube_github_binding

Provides a Sonarqube GitHub binding resource. This can be used to create and manage the binding between a
GitHub repository and a SonarQube project

## Example: Create a GitHub binding

```terraform
resource "sonarqube_alm_github" "github-alm" {
  app_id         = "12345"
  client_id      = "56789"
  client_secret  = "secret"
  key           = "myalm"
  private_key    = "myprivate_key"
  url           = "https://api.github.com"
  webhook_secret = "mysecret"
}

resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "my_project"
  visibility = "public"
}
resource "sonarqube_github_binding" "github-binding" {
  alm_setting = sonarqube_alm_github.github-alm.key
  project    = "my_project"
  repository = "myorg/myrepo"
}
```

## Argument Reference

The following arguments are supported:

- alm_setting - (Required) - GitHub ALM setting key
- monorepo - (Optional) - Is this project part of a monorepo. Default value: false
- project - (Required) - Project key
- repository - (Required) - The full name of your GitHub repository, including the organization, case-sensitive. Maximum length: 256
- summary_comment_enabled - (optional) - Enable/disable summary in PR discussion tab. Default value: true

## Attributes Reference

The following attributes are exported:

- project - Project key.
- repository - GitHub Repository.
- alm_setting - The unique key of the GitHub instance setting.

## Import

Bindings can be imported using their ID

```terraform
terraform import sonarqube_github_binding.github-binding project/repository
```