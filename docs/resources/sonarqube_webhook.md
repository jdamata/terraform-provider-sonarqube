# sonarqube_user

Provides a Sonarqube Webhook resource. This can be used to manage Sonarqube webhooks.

## Example: create a webhook

```terraform
resource "sonarqube_webhook" "webhook" {
  name = "terraform-webhook"
  url  = "https://my-webhook-destination.example.com"
}
```

## Argument Reference

The following arguments are supported:

- name - (Required) The name of the webhook to create. This will be displayed in the Sonarqube administration console.
- url - (Required) The URL to send event payloads to. This must begin with either `https://` or `http://`.
- secret 0 (Optional) The secret to send with the event payload.


## Attributes Reference

The following attributes are exported:

- id - The ID (key) of the webhook.

## Import

Webhooks can be imported using their ID (key):

```bash
terraform import sonarqube_webhook.webhook AXnN9NuxdWLvsEEPOr2g
```
