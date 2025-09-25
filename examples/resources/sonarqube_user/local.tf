resource "sonarqube_user" "user" {
  login_name = "terraform-test"
  name       = "terraform-test"
  password   = "secret-sauce37!"
}
