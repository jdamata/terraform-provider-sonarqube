resource "sonarqube_permission_template" "template" {
  name                = "Internal-Projects"
  description         = "These are internal projects"
  project_key_pattern = "internal.*"
}
