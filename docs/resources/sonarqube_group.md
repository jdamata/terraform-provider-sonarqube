---
page_title: "sonarqube_group Resource - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Provides a Sonarqube Group resource. This can be used to create and manage Sonarqube Groups.
---

# sonarqube_group (Resource)

Provides a Sonarqube Group resource. This can be used to create and manage Sonarqube Groups.

## Example: create a group

```terraform
resource "sonarqube_group" "project_users" {
    name        = "Project-Users"
    description = "This is a group"
}
```

## Argument Reference

The following arguments are supported:

- name - (Required) The name of the Group to create. Changing this forces a new resource to be created.
- description - (Optional) Description of the Group.

## Attributes Reference

The following attributes are exported:

- id - The ID of the Group.

## Import

Groups can be imported using their ID, and only support in the Sonarqube version <= 9.9

```terraform
terraform import sonarqube_group.group 101
```
