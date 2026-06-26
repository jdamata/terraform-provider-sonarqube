data "sonarqube_measures" "my_project" {
  component   = "my_project"
  metric_keys = ["coverage", "reliability_rating", "security_hotspots"]
}

output "coverage" {
  value = data.sonarqube_measures.my_project.measures
}
