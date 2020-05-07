# terraform-provider-sonarqube
Terraform provider for managing Sonarqube configuration

## Installation
Download the binary from the [Releases](https://github.com/jdamata/terraform-provider-sonarqube/releases/latest) page and place it in: ```~/.terraform.d/plugins``` or ```%APPDATA%\terraform.d\plugins```

## Docs
[Provider configuration](docs/provider.md)

Resources:
- [sonarqube_qualitygate](docs/sonarqube_qualitygate.md)

## Development
```bash
docker run -d --name sonarqube -p 9000:9000 sonarqube:latest
export TF_LOG=TRACE
make
terraform init
terraform apply
```
