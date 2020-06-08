# sonarqube_user_token

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

## Argument Reference

The following arguments are supported:

- login_name - (Required) The login name of the User for which the token should be created. Changing this forces a new resource to be created.
- name - (Required) The name of the Token to create. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

- id - The ID of the Token.
- token - The Token value.

## Import

Import is not supported for this resource.
