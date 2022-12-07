# sonarqube_permissions

Provides a Sonarqube Permissions resource. This can be used to manage global and project permissions.

## Example: Set global admin permissions for a group called "my-admins"

```terraform
resource "sonarqube_permissions" "my_global_admins" {
    group_name  = "my-admins"
    permissions = ["admin"]
}
```

## Example: Set project admin permissions for a group called "my-project-admins"

```terraform
resource "sonarqube_permissions" "my_project_admins" {
    group_name  = "my-project-admins"
    project_key = "my-project"
    permissions = ["admin"]
}
```

## Example: Set project admin permissions for a group called "my-project-admins on a permission template"

```terraform
resource "sonarqube_permissions" "internal_admins" {
    group_name  = "my-internal-admins"
    template_id = sonarqube_permission_template.template.id
    permissions = ["admin"]
}
```

## Example: Set codeviewer & user permissions on project level for a user called "johndoe"

```terraform
resource "sonarqube_permissions" "john_project_read" {
    login_name  = "johndoe"
    project_key = "my-project"
    permissions = ["codeviewer", "user"]
}
```

## Argument Reference

The following arguments are supported:

- login_name - (Optional) The name of the user that should get the specified permissions. Changing this forces a new resource to be created. Cannot be used with `group_name`
- group_name - (Optional) The name of the Group that should get the specified permissions. Changing this forces a new resource to be created. Cannot be used with `login_name`
- project_key - (Optional) Specify if you want to apply project level permissions. Changing this forces a new resource to be created. Cannot be used with `template_id & template_name`
- template_id - (Optional) Specify if you want to apply the permissions to a permission template. Changing this forces a new resource to be created. Cannot be used with `project_key & template_name`
- template_name - (Optional) Specify if you want to apply the permissions to a permission template. Changing this forces a new resource to be created. Cannot be used with `project_key & template_id`
- permissions - (Required) A list of permissions that should be applied. Changing this forces a new resource to be created. Possible values are: `admin`, `codeviewer`, `issueadmin`, `securityhotspotadmin`, `scan`, `user`.

**Note:** To prevent unwanted diffs, you should sort the permissions alphabetically.

## Attributes Reference

The following attributes are exported:

- id - A randomly generated UUID for the permission entry.

## Import

Importing is not supported for the `sonarqube_permissions` resource.
