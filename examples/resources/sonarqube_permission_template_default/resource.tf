resource "sonarqube_permission_template_default" "projects" {
  template_id = sonarqube_permission_template.default.id
  qualifier   = "TRK"
}
