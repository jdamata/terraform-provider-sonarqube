---
page_title: "sonarqube_alm_github Resource - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Provides a Sonarqube GitHub Alm/Devops Platform Integration resource. This can be used to create and manage a Alm/Devops
  Platform Integration for GitHub.
---

# sonarqube_alm_github (Resource)

Provides a Sonarqube GitHub Alm/Devops Platform Integration resource. This can be used to create and manage a Alm/Devops
Platform Integration for GitHub.

## Example: Create a GitHub Alm Integration

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
```

## Argument Reference

The following arguments are supported:

- app_id - (Required) - GitHub App ID. Maximum length: 80
- client_id - (Required) - GitHub App Client ID. Maximum length: 80
- client_secret - (Required) - GitHub App Client Secret. Maximum length: 160
- key - (Required) - Unique key of the GitHub instance setting. Maximum length: 200
- private_key - (Required) - GitHub App private key. Maximum length: 2500
- url - (Required) - GitHub API URL. Maximum length: 2000
- webhook_secret - (Optional) - GitHub App Webhook Secret. Maximum length: 160

## Attributes Reference

The following attributes are exported:

- app_id - GitHub App ID.
- client_id - GitHub App Client ID.
- key - The unique key of the GitHub instance setting.
- url - GitHub API URL.
