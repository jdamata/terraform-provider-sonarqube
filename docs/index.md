# Provider configuration

The sonarqube provider is used to configure sonarqube. The provider needs to be configured with a url, user and password.

## Example Usage
```terraform
provider "sonarqube" {
    user   = "admin"
    pass   = "admin" 
    host   = "http://127.0.0.1:9000"
}
```

## Argument Reference
The following arguments are supported:

- user - (Required) Sonarqube user. This can also be set via the SONARQUBE_USER environment variable.
- pass - (Required) Sonarqube pass. This can also be set via the SONARQUBE_PASS environment variable.
- host - (Required) Sonarqube url. This can be also be set via the SONARQUBE_HOST environment variable.
- version - (Optional) The version of Sonarqube. When specified, the provider will avoid requesting this from the server during the initialization process. This can be helpful when using the same Terraform code to install Sonarqube and configure it.  
