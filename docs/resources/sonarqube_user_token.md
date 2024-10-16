---
page_title: "sonarqube_user_token Resource - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Provides a Sonarqube User token resource. This can be used to manage Sonarqube User tokens.
---

# sonarqube_user_token (Resource)

Provides a Sonarqube User token resource. This can be used to manage Sonarqube User tokens.

## Example: create a user, user token and output the token value

```terraform
resource "sonarqube_user" "user" {
  login_name = "terraform-test"
  name       = "terraform-test"
  password   = "secret-sauce37!"
}

resource "sonarqube_user_token" "token" {
  login_name = sonarqube_user.user.login_name
  name       = "my-token"
}

output "user_token" {
  value = sonarqube_user_token.token.token
}
```

## Example: create an expiring global analysis token and output the token value

```terraform
resource "sonarqube_user_token" "token" {
  name = "my-token"
  type = "GLOBAL_ANALYSIS_TOKEN"
  expiration_date = "2099-01-01"
}

output "global_analysis_token" {
  value = sonarqube_user_token.token.token
}
```

## Example: create a project, project analysis token, and output the token value

```terraform
resource "sonarqube_user_token" "token" {
  name = "my-token"
  type = "PROJECT_ANALYSIS_TOKEN"
  project_key = "my-project"
}

output "project_analysis_token" {
  value = sonarqube_user_token.token.token
}
```

## Argument Reference

The following arguments are supported:

- login_name - (Optional) The login name of the User for which the token should be created. If not set, the token is created for the authenticated user. Changing this forces a new resource to be created.
- name - (Required) The name of the Token to create. Changing this forces a new resource to be created.
- expiration_date - (Optional) The expiration date of the token being generated, in ISO 8601 format (YYYY-MM-DD). If not set, default to no expiration.
- type - (Optional) The kind of Token to create. Changing this forces a new resource to be created. Possible values are USER_TOKEN, GLOBAL_ANALYSIS_TOKEN, or PROJECT_ANALYSIS_TOKEN. Defaults to USER_TOKEN. If set to PROJECT_ANALYSIS_TOKEN, then the project_key must also be specified.
- project_key - (Optional) The key of the only project that can be analyzed by the PROJECT_ANALYSIS TOKEN being created. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

- id - The ID of the Token.
- token - The Token value.

## Import

Import is not supported for this resource.
