provider "sonarqube" {
    user = "admin"
    pass = "admin" 
    url = "http://127.0.0.1:9000"
}

resource "sonarqube_qualitygate" "test" {
    name = "test"
}

output "guality-gate-id" {
    value = sonarqube_qualitygate.test.id
}