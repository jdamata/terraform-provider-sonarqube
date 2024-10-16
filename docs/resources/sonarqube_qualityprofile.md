---
page_title: "sonarqube_qualityprofile Resource - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Provides a Sonarqube Quality Profile resource. This can be used to create and manage Sonarqube Quality Profiles.
---

# sonarqube_qualityprofile (Resource)

Provides a Sonarqube Quality Profile resource. This can be used to create and manage Sonarqube Quality Profiles.

## Example: create a quality profile

```terraform
resource "sonarqube_qualityprofile" "main" {
    name     = "example"
    language = "js"
    is_default = false
    parent = "sonar way"
}
```

## Argument Reference

The following arguments are supported:

- name - (Required) The name of the Quality Profile to create. Maximum length 100
- language - (Required) Quality profile language. Must be one of "cs", "css", "flex", "go", "java", "js", "jsp", "kotlin", "php", "py", "ruby", "scala", "ts", "vbnet", "web", "xml"
- is_default - (Optional) When set to true this will make the added Quality Profile default
- parent - (Optional) When a parent is provided the quality profile will inherit it's rules

## Attributes Reference

The following attributes are exported:

- name - Name of the Sonarqube Quality Profile
- id - ID of the Sonarqube Quality Profile
- key - ID of the Sonarqube Quality Profile

## Import

Quality Profiles can be imported using their ID

```terraform
terraform import sonarqube_qualityprofile.main AU-Tpxb--iU5OvuD2FLy
```
