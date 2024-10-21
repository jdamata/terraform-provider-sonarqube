## Example: change the external identity to SAML

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
