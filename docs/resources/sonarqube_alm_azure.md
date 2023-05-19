# sonarqube_alm_azure

Provides a Sonarqube Azure Devops Alm/Devops Platform Integration resource. This can be used to create and manage a Alm/Devops
Platform Integration for Azure Devops.

## Example: Create an Azure Devops Alm Integration

```terraform
resource "sonarqube_alm_azure" "az1" {
  key                   = "az1"
  personal_access_token = "my_pat"
  url                   = "https://dev.azure.com/my-org"
}
```

## Argument Reference

The following arguments are supported:

- key - (Required) - Unique key of the azure alm instance setting. Maximum length: 200
- personal_access_token - (Required) - Azure Devops Personal Access Token. Maximum length: 2000
- url - (Required) - Azure Devops Organization URL. Maximum length: 2000

## Attributes Reference

The following attributes are exported:

- key - The unique key of the azure alm instance setting.
- url - Azure Devops Organization URL.

## Import

Resource can be imported providing their Alm Instance Key and Azure DevOps Personal Access Token

```terraform
terraform import sonarqube_alm_azure.az1 key/personal_access_token
```
