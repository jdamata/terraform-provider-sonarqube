[![Build Status](https://cloud.drone.io/api/badges/jdamata/terraform-provider-sonarqube/status.svg)](https://cloud.drone.io/jdamata/terraform-provider-sonarqube)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=jdamata_terraform-provider-sonarqube&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=jdamata_terraform-provider-sonarqube)
[![Go Report Card](https://goreportcard.com/badge/github.com/jdamata/terraform-provider-sonarqube)](https://goreportcard.com/report/github.com/jdamata/terraform-provider-sonarqube)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=jdamata_terraform-provider-sonarqube&metric=coverage)](https://sonarcloud.io/dashboard?id=jdamata_terraform-provider-sonarqube)
[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](https://opensource.org/licenses/)

# terraform-provider-sonarqube
Terraform provider for managing Sonarqube configuration

## Installation
Download the binary from the [Releases](https://github.com/jdamata/terraform-provider-sonarqube/releases/latest) page and place it in: ```~/.terraform.d/plugins``` or ```%APPDATA%\terraform.d\plugins```

## Docs
[Provider configuration](docs/provider.md)

Resources:
- [sonarqube_group](docs/sonarqube_group.md)
- [sonarqube_permissions](docs/sonarqube_permissions.md)
- [sonarqube_permission_template](docs/sonarqube_permission_template.md)
- [sonarqube_plugin](docs/sonarqube_plugin.md)
- [sonarqube_project](docs/sonarqube_project.md)
- [sonarqube_qualitygate](docs/sonarqube_qualitygate.md)
- [sonarqube_qualitygate_condition](docs/sonarqube_qualitygate_condition.md)
- [sonarqube_qualitygate_project_association](docs/sonarqube_qualitygate_project_association.md)
- [sonarqube_user](docs/sonarqube_user.md)
- [sonarqube_user_token](docs/sonarqube_user_token.md)