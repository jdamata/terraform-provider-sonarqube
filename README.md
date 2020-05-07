# terraform-provider-sonarqube
Terraform provider for managing Sonarqube configuration

## Usage
You will first need to download the binary from the [Releases](https://github.com/jdamata/terraform-provider-sonarqube/releases/latest) page and place it in: ~/.terraform.d/plugins or %APPDATA%\terraform.d\plugins

```terraform
provider "sonarqube" {
    user = "admin"
    pass = "admin" 
    url = "http://127.0.0.1:9000"
}

resource "sonarqube_qualitygate" "test" {
    name = "test"
}
```

## Development
```bash
docker run -d --name sonarqube -p 9000:9000 sonarqube:latest
export TF_LOG=TRACE
make
terraform init
terraform apply
```
