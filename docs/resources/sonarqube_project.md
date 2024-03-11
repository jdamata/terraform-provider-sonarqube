# sonarqube_project
Provides a Sonarqube Project resource. This can be used to create and manage Sonarqube Project.

## Example: create a project
```terraform
resource "sonarqube_project" "main" {
    name       = "SonarQube"
    project    = "my_project"
    visibility = "public" 
}
```

## Example: a project with associated settings
```terraform
resource "sonarqube_project" "main" {
    name       = "SonarQube"
    project    = "my_project"
    visibility = "public" 

    setting {
        key   = "sonar.demo"
        value = "sonarqube@example.org"
    }
}
```

## Argument Reference
The following arguments are supported:

- name - (Required) The name of the Project to create
- project - (Required) Key of the project. Maximum length 400. All letters, digits, dash, underscore, period or colon.
- visibility - (Optional) Whether the created project should be visible to everyone, or only specific user/groups. If no visibility is specified, the default project visibility of the organization will be used. The only possible values are `public`
and `private`. 
- tags - (Optional) A list of tags to put on the project.
- setting - (Optional) The definition of a Setting to be used by this Portfolio as documented in the `setting` block below.

A `setting` block supports:

- key - (Required) Setting key
- value - (Optional) Single valued setting value
- values - (Optional) Multi-valued setting values
- field_values - (Optional) Multi-field setting values

One of value, values, field_values _must_ be supplied

## Attributes Reference
The following attributes are exported:
- project - (Required) Key of the project

## Import 
Projects can be imported using their project key

```terraform
terraform import sonarqube_project.main my_project
```

