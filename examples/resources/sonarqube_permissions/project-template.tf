resource "sonarqube_permissions" "internal_admins" {
  group_name  = "my-internal-admins"
  template_id = sonarqube_permission_template.template.id
  permissions = ["admin"]
}
