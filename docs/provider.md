# Provider configuration

The sonarqube provider is used to configure sonarqube. The provider needs to be configured with a url, user and password.

## Example Usage
```terraform
provider "sonarqube" {
    user = "admin"
    pass = "admin" 
    url = "http://127.0.0.1:9000"
}
```

## Argument Reference
The following arguments are supported:

- user - (Required) Sonarqube user. This can also be set via the SONARQUBE_USER environment variable.
- pass - (Required) Sonarqube pass. This can also be set via the SONARQUBE_PASS environment variable.
- url - (Required) Sonarqube url. This can be also be set via the SONARQUBE_URL environment variable.