resource "sonarqube_qualitygate" "main" {
  name       = "example"
  is_default = true

  condition {
    metric    = "new_coverage"
    op        = "LT"
    threshold = "50"
  }

  condition {
    metric    = "vulnerabilities"
    threshold = "10"
    op        = "GT"
  }
}
