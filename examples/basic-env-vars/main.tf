terraform {
  required_providers {
    sonarqube = {
      source = "jdamata/sonarqube"
    }
  }
}

provider "sonarqube" {}

resource "sonarqube_project" "tf-postfix-test" {
  name       = "tf-postfix-test"
  project    = "tf-postfix-test"
  visibility = "public"
}
