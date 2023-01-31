terraform {
  required_providers {
    sonarqube = {
      source = "terraform-example.com/exampleprovider/example"
    }
  }
}

provider "sonarqube" {
  user = "admin"
  pass = "admin1"
  host = "http://localhost:9001/"
}

resource "sonarqube_project" "tf-postfix-test" {
  name       = "tf-postfix"
  project    = "tf-postfix"
  visibility = "public"
}
#resource "sonarqube_project_main_branch" "mybranch" {
#  name    = "main3"
#  project = sonarqube_project.tf-postfix-test.project
#}