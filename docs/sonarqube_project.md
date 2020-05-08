# sonarqube_project
Provides a Sonarqube Project resource. This can be used to create and manage Sonarqube Project.

## Example: create a quality gate
```terraform
resource "sonarqube_project" "main" {
    name       = "SonarQube"
    project    = "my_project"
    visibility = "public" 
}

## Argument Reference
The following arguments are supported:

- name - (Required) The name of the Project to create
- project - (Required) Key of the project. Maximum length 400. All letters, digits, dash, underscore, period or colon.
- visibility - (Optional) Whether the created project should be visible to everyone, or only specific user/groups. If no visibility is specified, the default project visibility of the organization will be used.

## Attributes Reference
The following attributes are exported:
- project - (Required) Key of the project

## Import 
Projects can be imported using their project key

```terraform
terraform import sonarqube_qualitygate.main my_project
```

