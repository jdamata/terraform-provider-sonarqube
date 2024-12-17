---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarqube_qualityprofile Resource - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Provides a Sonarqube Quality Profile resource. This can be used to create and manage Sonarqube Quality Profiles.
---

# sonarqube_qualityprofile (Resource)

Provides a Sonarqube Quality Profile resource. This can be used to create and manage Sonarqube Quality Profiles.

## Example Usage

```terraform
resource "sonarqube_qualityprofile" "main" {
  name       = "example"
  language   = "js"
  is_default = false
  parent     = "sonar way"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `language` (String) Quality profile language. Must be one of "cs", "css", "flex", "go", "java", "js", "jsp", "kotlin", "php", "py", "ruby", "scala", "ts", "vbnet", "web", "xml"
- `name` (String) The name of the Quality Profile to create. Maximum length 100

### Optional

- `is_default` (Boolean) When set to true this will make the added Quality Profile default
- `parent` (String) When a parent is provided the quality profile will inherit it's rules

### Read-Only

- `id` (String) The ID of this resource.
- `key` (String) ID of the Sonarqube Quality Profile