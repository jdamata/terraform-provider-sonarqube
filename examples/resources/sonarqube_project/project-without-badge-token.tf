resource "sonarqube_project" "without_badge_token" {
  name                 = "My project"
  project              = "my-project"
  visibility           = "private"
  retrieve_badge_token = false
}
