# sonarqube_azure_binding

Provides a Sonarqube Azure Devops binding resource. This can be used to create and manage the binding between an
Azure Devops repository and a SonarQube project

## Example: Create an Azure Devops binding

```terraform
resource "sonarqube_alm_azure" "az1" {
  key           = "az1"
  personal_access_token    = "my_pat"
  url           = "https://dev.azure.com/my-org"
}

resource "sonarqube_project" "main" {
  name       = "SonarQube"
  project    = "my_project"
  visibility = "public"
}

resource "sonarqube_azure_binding" "az1-my_project-my_repo" {
  alm_setting = sonarqube_alm_azure.az1.key
  project    = sonarqube_project.main.project
  repository = "my_repo"
}
```

## Argument Reference

The following arguments are supported:

- alm_setting - (Required) - azure ALM setting key
- monorepo - (Optional) - Is this project part of a monorepo. Default value: false
- project - (Required) - Project key
- repository - (Required) - The name of your Azure Devops repository.

## Attributes Reference

The following attributes are exported:

- project - Project key.
- repository - Azure Devops Repository.
- alm_setting - The unique key of the azure alm instance setting.

## Import

Bindings can be imported using their ID

```terraform
terraform import sonarqube_azure_binding.az1-my_project-my_repo project/repository
```