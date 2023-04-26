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

* [Terraform](https://www.terraform.io/downloads.html)
* [Go](http://www.golang.org)
* [Docker Engine](https://docs.docker.com/engine/install/)

You will also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `${GOPATH}/bin` to your `$PATH`.

To compile the provider, run `make`. This will install the provider into your GOPATH.

In order to run the full suite of Acceptance tests, run `make -i testacc`. These tests require Docker to be installed on the machine that runs them, and do not create any remote resources.

```sh
$ make -i testacc
```

## Debugging the Provider

You can debug the provider using a tool similar to [delve](https://github.com/go-delve/delve) as described in the Hashicorp article [Debugger-Based Debugging](https://developer.hashicorp.com/terraform/plugin/debugging#debugger-based-debugging).

### Starting a debug session

#### Visual Studio Code

Using [Visual Studio Code](https://code.visualstudio.com/) with the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.go) installed, you would add a configuration similar to this to your `launch.json` file:
```json
{
    "name": "Debug Terraform Provider",
    "type": "go",
    "request": "launch",
    "mode": "debug",
    "program": "${workspaceFolder}",
    "env": {
        <any environment variables you need>
    },
    "args": [
        "-debug"
    ],
    "showLog": true
}
```
Start the process using this configuration and then follow the instructions provided in the debug console:
```shell
Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable with the following:

	Command Prompt:	set "TF_REATTACH_PROVIDERS={"registry.terraform.io/jdamata/sonarqube":{"Protocol":"grpc","ProtocolVersion":5,"Pid":2748,"Test":true,"Addr":{"Network":"tcp","String":"127.0.0.1:56560"}}}"

	PowerShell:	$env:TF_REATTACH_PROVIDERS='{"registry.terraform.io/jdamata/sonarqube":{"Protocol":"grpc","ProtocolVersion":5,"Pid":2748,"Test":true,"Addr":{"Network":"tcp","String":"127.0.0.1:56560"}}}'
```

#### Delve CLI

With the delve CLI you would start a delve debugging session:
```shell
dlv exec --accept-multiclient --continue --headless ./terraform-provider-example -- -debug
```

and copy the line starting `TF_REATTACH_PROVIDERS` from your provider's output, again setting it according to your command prompt/shell.

### Attaching to the debug session

In order to use the debug instance of the provider, set `TF_REATTACH_PROVIDERS` as described previously, set some breakpoints in the provider, and start your terraform script as normal.

Happy debugging!
