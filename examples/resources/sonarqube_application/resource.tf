resource "sonarqube_project" "api" {
  name       = "Payments API"
  project    = "payments-api"
  visibility = "private"
}

resource "sonarqube_project" "web" {
  name       = "Payments Web"
  project    = "payments-web"
  visibility = "private"
}

resource "sonarqube_application" "payments" {
  key         = "payments"
  name        = "Payments"
  description = "Everything that makes up the payments product"
  visibility  = "private"
  projects    = [sonarqube_project.api.project, sonarqube_project.web.project]
  tags        = ["payments", "team-a"]
}
