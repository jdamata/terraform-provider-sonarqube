terraform {
  required_providers {
    sonarqube = {
      source  = "jdamata/sonarqube"
      version = "~> 0.15.8"
    }
  }
}

provider "sonarqube" {
  host  = "http://localhost:9001"
  user = "admin"
  pass = "admin"
}

resource "sonarqube_qualitygate" "base" {
    name = "base_gate"
}

resource "sonarqube_qualitygate_condition" "qualitygate_condition" {
    gatename  = sonarqube_qualitygate.base.id
    metric    = "comment_lines_density"
    threshold = "68"
    op        = "LT"
}

resource "sonarqube_qualitygate" "copy" {
    depends_on = [sonarqube_qualitygate.base, sonarqube_qualitygate_condition.qualitygate_condition]
    name = "copy_gate"
    copy_from = sonarqube_qualitygate.base.name
}