# sonarqube_github_binding

Provides a Sonarqube GitHub binding resource. This can be used to create and manage the binding between a
GitHub repository and a SonarQube project

## Example: Create a GitHub binding

```terraform
resource "sonarqube_alm_github" "github-alm" {
  appid         = "12345"
  clientid      = "56789"
  clientsecret  = "secret"
  key           = "myalm"
  privatekey    = "myprivatekey"
  url           = "https://api.github.com"
  webhooksecret = "mysecret"
}

resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "my_project"
  visibility = "public"
}
resource "sonarqube_github_binding" "github-binding" {
  almsetting = sonarqube_alm_github.github-alm.key
  project    = sonarqube_project.main.project
  repository = "myrepo"
}
```

## Argument Reference

The following arguments are supported:

- almsetting - (Required) - GitHub ALM setting key
- monorepo - (Optional) - Is this project part of a monorepo. Default value: no
- project - (Required) - Project key
- repository - (Required) - GitHub Repository. Maximum length: 256
- summarycommentenabled - (optional) - Enable/disable summary in PR discussion tab. Default value: true

## Attributes Reference

The following attributes are exported:

- project - Project key.
- repository - GitHub Repository.
- almsetting - The unique key of the GitHub instance setting.

