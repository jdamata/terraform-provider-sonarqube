# terraform-provider-sonarqube
Terraform provider for managing Sonarqube configuration


## Development
```bash
docker run -d --name sonarqube -p 9000:9000 sonarqube:latest
export TF_LOG=TRACE
make
terraform init
terraform apply
```
