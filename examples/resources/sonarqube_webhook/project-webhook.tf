resource "sonarqube_project" "project" {
  name       = "project"
  project    = "project"
  visibility = "public"
}

resource "sonarqube_webhook" "webhook" {
  name    = "terraform-webhook"
  url     = "https://my-webhook-destination.example.com"
  project = sonarqube_project.project.name
}
