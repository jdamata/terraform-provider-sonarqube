resource "sonarqube_permissions" "john_project_read" {
  login_name  = "johndoe"
  project_key = "my-project"
  permissions = ["codeviewer", "user"]
}
