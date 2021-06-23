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
}

resource "sonarqube_qualitygate_project_association" "this" {
  for_each = var.quality_gates

  gatename   = sonarqube_qualitygate.this[each.key].id
  projectkey = sonarqube_project.this.project
}

resource "sonarqube_qualitygate_condition" "this" {
  for_each = var.quality_gates

  gatename  = sonarqube_qualitygate.this[each.key].id
  metric    = each.value.metric
  threshold = each.value.threshold
  op        = each.value.operator
}