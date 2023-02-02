# sonarqube_alm_github

Provides a Sonarqube GitHub Alm/Devops Platform Integration resource. This can be used to create and manage a Alm/Devops
Platform Integration for GitHub.

## Example: Create a GitHub Alm Integration

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
```

## Argument Reference

The following arguments are supported:

- appid - (Required) - GitHub App ID. Maximum length: 80
- clientid - (Required) - GitHub App Client ID. Maximum length: 80
- clientsecret - (Required) - GitHub App Client Secret. Maximum length: 160
- key - (Required) - Unique key of the GitHub instance setting. Maximum length: 200
- privatekey - (Required) - GitHub App private key. Maximum length: 2500
- url - (Required) - GitHub API URL. Maximum length: 2000
- webhooksecret - (Optional) - GitHub App Webhook Secret. Maximum length: 160

## Attributes Reference

The following attributes are exported:

- appid - GitHub App ID.
- clientid - GitHub App Client ID.
- key - The unique key of the GitHub instance setting.
- url - GitHub API URL.

