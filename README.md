# terraform-provider-sonarqube
Terraform provider for managing Sonarqube configuration


## Development
```bash
docker run -d --name sonarqube -p 9000:9000 sonarqube:latest
make
terraform init
terraform plan
```
