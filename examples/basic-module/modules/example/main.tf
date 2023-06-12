terraform {
  required_providers {
    sonarqube = {
      source  = "jdamata/sonarqube"
    }
  }
}

# Sonarqube project configuration
#
# This file contains all configurations for a sonarqube project, which is used for multiple code validations

resource "sonarqube_project" "this" {
  name       = var.name
  project    = var.name
  visibility = "public"
}

# Associate custom quality gates to the project. 
# Configuration taken from `quality_gates` map
resource "sonarqube_qualitygate" "this" {
  for_each = var.quality_gates

  name = each.key

  condition {
    metric    = each.value.metric
    threshold = each.value.threshold
    op        = each.value.operator
  }
}

resource "sonarqube_qualitygate_project_association" "this" {
  for_each = var.quality_gates

  gatename   = sonarqube_qualitygate.this[each.key].id
  projectkey = sonarqube_project.this.project
}
