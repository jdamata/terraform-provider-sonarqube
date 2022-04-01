# Provider configuration

The sonarqube provider is used to configure sonarqube. The provider needs to be configured with a url, and either with user and password or token.

## Example: Authenticate with username and password
```terraform
terraform {
  required_providers {
    sonarqube = {
      source = "jdamata/sonarqube"
    }
  }
}

provider "sonarqube" {
    user   = "admin"
    pass   = "admin" 
    host   = "http://127.0.0.1:9000"
}
```

## Example: Authenticate with token
```terraform
terraform {
  required_providers {
    sonarqube = {
      source = "jdamata/sonarqube"
    }
  }
}

provider "sonarqube" {
    token  = "d4at55a6f7r199bd707h39625685510880gbf7ff"
    host   = "http://127.0.0.1:9000"
}
```

## Argument Reference
The following arguments are supported:

- user - (Optional) Sonarqube user. This can also be set via the SONARQUBE_USER environment variable.
- pass - (Optional) Sonarqube pass. This can also be set via the SONARQUBE_PASS environment variable.
- token - (Optional) Sonarqube token. This can also be set via the SONARQUBE_TOKEN environment variable.
- host - (Required) Sonarqube url. This can be also be set via the SONARQUBE_HOST environment variable.
- installed_version - (Optional) The version of the Sonarqube server. When specified, the provider will avoid requesting this from the server during the initialization process. This can be helpful when using the same Terraform code to install Sonarqube and configure it.
- tls_insecure_skip_verify - (Optional) Allows ignoring insecure certificates when set to true. Defaults to false. Disabling TLS verification is dangerous and should only be done for local testing.
