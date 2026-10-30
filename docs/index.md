---
page_title: "Provider: sonarqube"
description: |-
    The sonarqube provider is used to configure sonarqube. The provider needs to be configured with a url, and either with user and password or token.
---

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

- `user` - (Optional) Sonarqube user. This can also be set via the `SONARQUBE_USER` environment variable.
- `pass` - (Optional) Sonarqube pass. This can also be set via the `SONARQUBE_PASS` environment variable.
- `token` - (Optional) Sonarqube token. This can also be set via the `SONARQUBE_TOKEN` environment variable.
- `host` - (Required) Sonarqube url. This can be also be set via the `SONARQUBE_HOST` environment variable.
- `installed_version` - (Optional) The version of the Sonarqube server. When specified, the provider will avoid requesting this from the
  server during the initialization process. This can be helpful when using the same Terraform code to install Sonarqube and configure it.
- `tls_insecure_skip_verify` - (Optional) Allows ignoring insecure certificates when set to true. Defaults to false. Disabling TLS verification
  is dangerous and should only be done for local testing.
- `anonymize_user_on_delete` - (Optional) Allows anonymizing users on destroy. Requires Sonarqube version >= `9.7`. This can be helpful
  to comply with regulations like [GDPR](https://en.wikipedia.org/wiki/General_Data_Protection_Regulation).

## Badge-token sensitivity feature flag

Badge tokens remain non-sensitive by default for compatibility with public badge-URL outputs. Set `SONARQUBE_BADGE_TOKEN_SENSITIVE=true` in the Terraform process environment to opt in to sensitive badge attributes on project resources and data sources. Keep this schema flag consistent between plan and apply and across provider aliases. Root outputs using these attributes must then declare `sensitive = true`. Sensitivity does not remove values from state. See the project resource documentation for migration examples.
