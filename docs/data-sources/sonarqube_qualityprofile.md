---
page_title: "sonarqube_qualityprofile Data Source - terraform-provider-sonarqube"
subcategory: ""
description: |-
  Use this data source to get a Sonarqube qualityprofile resource
---

# sonarqube_qualityprofile (Data Source)

Use this data source to get a Sonarqube qualityprofile resource

## Example usage

```terraform
data "sonarqube_qualityprofile" "main" {
  name     = "example"
}
```

## Argument Reference

The following arguments are supported:

- name - (Required) The name of the quality profile

## Attributes Reference

The following attributes are exported:

- name - The name of the Quality Profile
- language - Quality profile language.
- is_default - Quality Profile default
