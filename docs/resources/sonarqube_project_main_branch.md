# sonarqube_project_main_branch
Provides a Sonarqube Project main branch resource. This can be used to create and manage a Sonarqube Projects main branch.

## Example: Change the name of the main branch
```terraform
resource "sonarqube_project" "main" {
    name       = "SonarQube"
    project    = "my_project"
    visibility = "public" 
}
resource "sonarqube_project_main_branch" "mybranch" {
  name    = "release"
  project = sonarqube_project.main.project
}
```

## Argument Reference
The following arguments are supported:

- name - (Required) The name you want the main branch to have.
- project - (Required) Key of the project. Maximum length 400. All letters, digits, dash, underscore, period or colon.

## Attributes Reference
The following attributes are exported:
- name - The name of the main branch
- project - Key of the project

## Import

Main branch can be imported using their ID

```terraform
terraform import sonarqube_project_main_branch.mybranch project/name
```