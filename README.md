# terraform-provider-sonarqube

[![release](https://github.com/jdamata/terraform-provider-sonarqube/actions/workflows/release.yaml/badge.svg)](https://github.com/jdamata/terraform-provider-sonarqube/actions/workflows/release.yaml)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=jdamata_terraform-provider-sonarqube&metric=sqale_rating)](https://sonarcloud.io/dashboard?id=jdamata_terraform-provider-sonarqube)
[![Go Report Card](https://goreportcard.com/badge/github.com/jdamata/terraform-provider-sonarqube)](https://goreportcard.com/report/github.com/jdamata/terraform-provider-sonarqube)
[![codecov](https://codecov.io/gh/jdamata/terraform-provider-sonarqube/branch/master/graph/badge.svg)](https://codecov.io/gh/jdamata/terraform-provider-sonarqube)
[![GPLv3 License](https://img.shields.io/badge/License-GPL%20v3-yellow.svg)](https://opensource.org/licenses/)

Terraform provider for managing Sonarqube configuration

This is a community provider and is not supported by Hashicorp.

## Installation

This provider has been published to the Terraform Registry at https://registry.terraform.io/providers/jdamata/sonarqube. Please visit the registry for documentation and installation instructions.

## Developing the Provider

Working on this provider requires the following:

- [Terraform](https://www.terraform.io/downloads.html)
- [Go](http://www.golang.org)
- [Docker Engine](https://docs.docker.com/engine/install/)

You will also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `${GOPATH}/bin` to your `$PATH`.

To compile the provider, run `make`. This will install the provider into your GOPATH.

In order to run the full suite of Acceptance tests, run `make -i testacc`. These tests require Docker to be installed on the machine that runs them. The tests do not create any remote resources.

```sh
$ make -i testacc
```

## Generate documentation

Documentation is generated using `tfplugindocs`. These are auto-generated when creating a PR to the project. 

If you wish to generate documentation locally to verify it's accuracy, run the following commands:

```sh
$ make tools
$ make docs
```

## Debugging the Provider

See [debugging.md](debugging.md)
