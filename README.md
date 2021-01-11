# terraform-provider-sonarqube

[![Build Status](https://cloud.drone.io/api/badges/jdamata/terraform-provider-sonarqube/status.svg)](https://cloud.drone.io/jdamata/terraform-provider-sonarqube)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=jdamata_terraform-provider-sonarqube&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=jdamata_terraform-provider-sonarqube)
[![Go Report Card](https://goreportcard.com/badge/github.com/jdamata/terraform-provider-sonarqube)](https://goreportcard.com/report/github.com/jdamata/terraform-provider-sonarqube)
[![codecov](https://codecov.io/gh/jdamata/terraform-provider-sonarqube/branch/master/graph/badge.svg)](https://codecov.io/gh/jdamata/terraform-provider-sonarqube)
[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](https://opensource.org/licenses/)

Terraform provider for managing Sonarqube configuration

This is a community provider and is not supported by Hashicorp.

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

## Developing the Provider

Working on this provider requires the following:

* [Terraform](https://www.terraform.io/downloads.html) 0.14+
* [Go](http://www.golang.org) (version requirements documented in the `go.mod` file)
* [Docker Engine](https://docs.docker.com/engine/install/) 20.10+ (for running acceptance tests)

You will also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `${GOPATH}/bin` to your `$PATH`.

To compile the provider, run `make`. This will install the provider into your GOPATH.

In order to run the full suite of Acceptance tests, run `make -i testacc`. These tests require Docker to be installed on the machine that runs them, and do not create any remote resources.

```sh
$ make -i testacc
```
