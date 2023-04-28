# Debugging the Provider

You can debug the provider using a tool similar to [delve](https://github.com/go-delve/delve) as described in the Hashicorp article [Debugger-Based Debugging](https://developer.hashicorp.com/terraform/plugin/debugging#debugger-based-debugging).

## Starting a debug session

### Visual Studio Code

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

### Delve CLI

With the delve CLI you would start a delve debugging session:
```shell
dlv exec --accept-multiclient --continue --headless ./terraform-provider-example -- -debug
```

and copy the line starting `TF_REATTACH_PROVIDERS` from your provider's output, again setting it according to your command prompt/shell.

## Attaching to the debug session

In order to use the debug instance of the provider, set `TF_REATTACH_PROVIDERS` as described previously, set some breakpoints in the provider, and start your terraform script as normal.

Happy debugging!
