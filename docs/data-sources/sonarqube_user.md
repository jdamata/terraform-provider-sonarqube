---
page_title: "sonarqube_user Data Source - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Use this data source to get a Sonarqube User resource
---

# sonarqube_user (Data Source)

Use this data source to get a Sonarqube User resource

## Example usage

```terraform
data "sonarqube_user" "user" {
  login_name = "terraform-test"
}
```

## Argument Reference

The following arguments are supported:

- login_name - (Required) The login name of the User

## Attributes Reference

The following attributes are exported:

- id - The ID of the User.
- login_name - The login name of the user
- email - The users email.
