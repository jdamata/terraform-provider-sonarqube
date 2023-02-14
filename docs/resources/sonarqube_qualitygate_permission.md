# sonarqube_qualitygate_permission

Provides a Sonarqube Quality Gate Permission resource. This can be used to assign `edit` permissions on a quality gate to users or groups.
This feature is available on SonarQube `9.2` or newer.

## Example: grant permission to a user

```terraform
resource "sonarqube_qualitygate" "main" {
    name = "my_qualitygate"
}

resource "sonarqube_user" "user" {
	login_name = "testuser"
	name       = "Test User"
	email      = "terraform-test@sonarqube.com"
	password   = "secret-sauce!"
}

resource "sonarqube_qualitygate_permission" "permission" {
	gatename = sonarqube_qualitygate.main.name
	type      = "user"
	subject   = sonarqube_user.user.login_name
}
```

## Argument Reference

The following arguments are supported:

- `gatename` - (Required) The name of the Quality Gate. Changing this forces a new resource to be created.
- `type` - (Required) The type of the subject to give permission. Only `user` and `group` are valid values. Changing this forces a new resource to be created.
- `subject` - (Required) A user's `login_name` or a group's `name` depending on the *type*. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the permission assignment.

## Import

Assignments can be imported using their ID (`<gatename>[<type>/<subject>]`):

```terraform
terraform import sonarqube_qualitygate_permission.permission my-qualitygate[user/my-user]
```
