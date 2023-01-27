# Basic sonarqube configuration using env vars 
Start sonarqube:

```sh
docker run --name sonarqube -d -p 9000:9000 sonarqube:latest
```

Run terraform commands to create the a sonarqube project, remove it from state and then re-add it to state.

```sh
export SONAR_HOST=http://127.0.0.1:9000
export SONAR_USER=admin
export SONAR_PASS=admin
terraform init
terraform plan
terraform apply --auto-approve
terraform state rm 'sonarqube_project.tf-postfix-test'
terraform import 'sonarqube_project.tf-postfix-test' tf-postfix-test
terraform destroy --auto-approve
```