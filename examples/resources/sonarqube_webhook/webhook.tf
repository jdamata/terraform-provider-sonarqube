resource "sonarqube_webhook" "webhook" {
  name = "terraform-webhook"
  url  = "https://my-webhook-destination.example.com"
}
