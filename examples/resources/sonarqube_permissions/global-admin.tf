resource "sonarqube_permissions" "my_global_admins" {
  group_name  = "my-admins"
  permissions = ["admin"]
}
