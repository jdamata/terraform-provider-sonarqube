# sonarqube_qualitygate_project_association
Provides a Sonarqube Quality Gate Project association resource. This can be used to associate a Quality Gate to a Project

## Example: create a quality gate project association
```terraform
resource "sonarqube_qualitygate" "main" {
    name = "my_qualitygate"
}

resource "sonarqube_project" "main" {
    name       = "SonarQube"
    project    = "my_project"
    visibility = "public" 
}

resource "sonarqube_qualitygate_project_association" "main" {
    gateid     = sonarqube_qualitygate.main.id
    projectkey = sonarqube_project.main.project
}
```

## Argument Reference
The following arguments are supported:

- gateid - (Required) The id of the Quality Gate
- projectkey - (Required) Key of the project. Maximum length 400. All letters, digits, dash, underscore, period or colon.