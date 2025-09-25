resource "sonarqube_permissions" "my_project_admins" {
  group_name  = "my-project-admins"
  project_key = "my-project"
  permissions = ["admin"]
}
