---
page_title: "sonarqube_gitlab_binding Resource - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Provides a Sonarqube GitLab binding resource. This can be used to create and manage the binding between a
  GitLab repository and a SonarQube project
---

# sonarqube_gitlab_binding (Resource)

Provides a Sonarqube GitLab binding resource. This can be used to create and manage the binding between a
GitLab repository and a SonarQube project

## Example: Create a GitLab binding

```terraform
resource "sonarqube_alm_gitlab" "gitlab-alm" {
  key                    = "myalm"
  personal_access_token  = "my_personal_access_token"
  url                    = "https://gitlab.com/api/v4"
}

resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "my_project"
  visibility = "public"
}

resource "sonarqube_gitlab_binding" "gitlab-binding" {
  alm_setting = sonarqube_alm_gitlab.gitlab-alm.key
  project    = "my_project"
  repository = "123"
}
```

## Argument Reference

The following arguments are supported:

- alm_setting - (Required) - GitLab ALM setting key
- monorepo - (Optional) - Is this project part of a monorepo. Default value: false
- project - (Required) - SonarQube project key. Changing this will force a new resource to be created.
- repository - (Required) - The GitLab project ID.

## Attributes Reference

The following attributes are exported:

- project - SonarQube unique project key.
- repository - GitLab project ID.
- alm_setting - The unique key of the GitLab instance setting.

## Import

Bindings can be imported using their ID

```terraform
terraform import sonarqube_gitlab_binding.gitlab-binding project/repository
```
