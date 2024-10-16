---
page_title: "sonarqube_qualityprofile_project_association Resource - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Provides a Sonarqube Quality Profile Project association resource. This can be used to associate a Quality Profile to a Project
---

# sonarqube_qualityprofile_project_association (Resource)

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
	project         = sonarqube_project.main.project
	language        = "js"
}
```

## Argument Reference

The following arguments are supported:

- quality_profile - (Required) Name of the Quality Profile
- project - (Required) Name of the project
- language - (Required) Quality profile language. Must be a langauge in this list https://next.sonarqube.com/sonarqube/web_api/api/languages/list

## Import

Quality Profiles Project Associations can be imported using a combination of quality profile name, project name and language

```terraform
terraform import sonarqube_qualityprofile_project_association.main my_quality_profile/SonarQube/js
```
