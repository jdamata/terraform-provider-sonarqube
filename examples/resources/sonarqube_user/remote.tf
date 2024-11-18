resource "sonarqube_user" "remote_user" {
  login_name = "terraform-test"
  name       = "terraform-test"
  email      = "terraform-test@sonarqube.com"
  is_local   = false
}
