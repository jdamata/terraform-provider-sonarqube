# The token must already have been issued for this project. Ensure its value is
# delivered to a durable secret store before retiring a previous generation.
resource "sonarqube_user_token" "analysis_metadata" {
  name          = "my-project-analysis-generation-1"
  type          = "PROJECT_ANALYSIS_TOKEN"
  project_key   = "my-project"
  metadata_only = true

  lifecycle {
    create_before_destroy = true
  }
}
