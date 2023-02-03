# sonarqube_user_external_entity

Updates the *external identity* of a *non local* Sonarqube User. This can be used to set the *Identity Provider* which should be used to
authenticate a specific user.

The Sonarqube API currently does not provide an endpoint to read the *external identity* setting of an user. 

## Example: change the external identity to SAML

```terraform
resource "sonarqube_user" "remote_user" {
  login_name = "terraform-test"
  name       = "terraform-test"
  email      = "terraform-test@sonarqube.com"
  is_local   = false
}

resource "sonarqube_user_external_identity" "remote_user" {
  login_name        = sonarqube_user.remote_user.login_name
  external_identity = "terraform-test@sonarqube.com"
  external_provider = "saml"
}
```

## Argument Reference

The following arguments are supported:

- `login_name` - (Required) The login name of the User to update. Changing this forces a new resource to be created.
- `external_identity` - (Required) The identifier of the User used by the Authentication Provider. Changing this forces a new resource to be 
  created.
- `external_provider` - (Required) The key of the Authentication Provider. The Authentication Provider must be activated on Sonarqube. Changing 
  this forces a new resource to be created.
