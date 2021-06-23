# Basic sonarqube configuration example

Start sonarqube:

```sh
docker run --name sonarqube -d -p 9000:9000 sonarqube:latest
```

Run terraform commands to create the a sonarqube project, remove it from state and then re-add it to state.

```sh
terraform init
terraform plan
terraform apply
terraform state rm 'sonarqube_project.tf-postfix-test'
terraform import 'sonarqube_project.tf-postfix-test' tf-postfix-test
```