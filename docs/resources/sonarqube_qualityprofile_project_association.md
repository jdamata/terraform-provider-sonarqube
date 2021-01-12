# sonarqube_qualityprofile_project_association
Provides a Sonarqube Quality Profile Project association resource. This can be used to associate a Quality Profile to a Project

## Example: create a quality profile project association
```terraform
resource "sonarqube_qualityprofile" "main" {
	name     = "my_quality_profile"
	language = "js"
}

resource "sonarqube_project" "main" {
    name       = "SonarQube"
    project    = "my_project"
    visibility = "public" 
}

resource "sonarqube_qualityprofile_project_association" "main" {
	quality_profile = sonarqube_qualityprofile.main.name
	project         = sonarqube_project.main.name
	language        = "js"
}
```

## Argument Reference
The following arguments are supported:

- quality_profile - (Required) Name of the Quality Profile
- project         - (Required) Name of the project
- language        - (Required) Quality profile language. Must be one of "cs", "css", "flex", "go", "java", "js", "jsp", "kotlin", "php", "py", "ruby", "scala", "ts", "vbnet", "web", "xml"

## Import 
Quality Profiles Project Associations can be imported using a combination of quality profile name and project name

```terraform
terraform import sonarqube_qualityprofile_project_association.main my_quality_profile/SonarQube
```