# sonarqube_plugin
Provides a Sonarqube Plugin resource. This can be used to create and manage Sonarqube Plugins.

## Example: create a project
```terraform
resource "sonarqube_plugin" "main" {
    key = "cloudformation"
}
```

## Argument Reference
The following arguments are supported:

- key - (Required) The key identifying the plugin to uninstall

## Attributes Reference
The following attributes are exported:
- key - (Required) The key identifying the plugin to uninstall

## Import 
Projects can be imported using their plugin key

```terraform
terraform import sonarqube_plugin.main key
```

## Notes
After installing or uninstalling a plugin, the sonarqube server needs to be restarted.