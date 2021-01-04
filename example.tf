terraform {
  required_providers {
    sonarqube = {
      version = "~> 0.1"
      source  = "github.com/jdamata/sonarqube"
    }
  }
}

provider "sonarqube" {
    user   = "admin"
    pass   = "admin1" 
    host   = "http://127.0.0.1:9000"
}

resource "sonarqube_project" "main" {
    name       = "Project 1"
    project    = "my_project"
    visibility = "public" 
}

resource "sonarqube_qualitygate" "main" {
    name = "my_qualitygate"
}

resource "sonarqube_qualitygate_project_association" "main" {
    gatename     = sonarqube_qualitygate.main.name
    projectkey = sonarqube_project.main.project
}

resource "sonarqube_qualitygate_condition" "main" {
    gatename = sonarqube_qualitygate.main.name
    metric = "vulnerabilities"
    error  = 12
    op     = "GT"
}

resource "sonarqube_plugin" "main" {
    key = "cloudformation"
}

resource "sonarqube_group" "project_users" {
    name        = "Project-Users"
    description = "This is a group"
}

resource "sonarqube_permissions" "my_global_admins" {
    group_name  = sonarqube_group.project_users.name
    permissions = ["admin"]
}

resource "sonarqube_permission_template" "template" {
    name                = "Internal-Projects"
    description         = "These are internal projects"
    project_key_pattern = "internal.*"
}

resource "sonarqube_permissions" "internal_admins" {
    group_name  = "sonar-administrators"
    template_id = sonarqube_permission_template.template.id
    permissions = ["admin"]
}
