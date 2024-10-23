resource "sonarqube_user_token" "token" {
  name        = "my-token"
  type        = "PROJECT_ANALYSIS_TOKEN"
  project_key = "my-project"
}

output "project_analysis_token" {
  value = sonarqube_user_token.token.token
}
