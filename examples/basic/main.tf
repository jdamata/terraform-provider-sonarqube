terraform {
  required_providers {
    sonarqube = {
      source = "jdamata/sonarqube"
    }
  }
}

provider "sonarqube" {
    user   = "admin"
    pass   = "admin" 
    host   = "http://127.0.0.1:9000"
}

resource "sonarqube_project" "tf-postfix-test" {
  name       = "tf-postfix-test"
  project    = "tf-postfix-test"
  visibility = "public"
}
