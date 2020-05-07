provider "sonarqube" {
    user = "admin"
    pass = "admin" 
    url = "http://127.0.0.1:9000"
}

resource "sonarqube_qualitygate" "main" {
    name = "example"
}

resource "sonarqube_qualitygate_condition" "main" {
    gateid = sonarqube_qualitygate.main.id
    metric = "vulnerabilities"
    error = 10
}