---
page_title: "sonarqube_group_member Resource - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Provides a Sonarqube Group Member resource. This can be used to add or remove user to or from Sonarqube Groups.
---

# sonarqube_group_member (Resource)

Provides a Sonarqube Group Member resource. This can be used to add or remove user to or from Sonarqube Groups.

## Example: add a user to a group

```terraform
resource "sonarqube_user" "user" {
    login_name = "terraform-test"
    name       = "terraform-test"
    password   = "secret-sauce37!"
}

resource "sonarqube_group" "project_users" {
    name        = "Project-Users"
    description = "This is a group"
}

resource "sonarqube_group_member" "project_users_member" {
    name       = sonarqube_group.project_users.name
    login_name = sonarqube_user.user.login_name
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the Group to add a member to. Changing this forces a new resource to be created.
- `login_name` - (Required) The `login_name` of the User to add as a member. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

- id - The ID of the Group Membership.

## Import

Group Members can be imported using their ID (`<name>[<login_name>]`):

```terraform
terraform import sonarqube_group_member.member group[user]
```
