provider "sonarqube" {
    user   = "admin"
    pass   = "admin" 
    host   = "127.0.0.1:9000"
    scheme = "http"
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
    gateid     = sonarqube_qualitygate.main.id
    projectkey = sonarqube_project.main.project
}

resource "sonarqube_qualitygate_condition" "main" {
    gateid = sonarqube_qualitygate.main.id
    metric = "vulnerabilities"
    error  = 12
    op     = "GT"
}

resource "sonarqube_plugin" "main" {
    key = "cloudformation"
}
