export GO111MODULE=on
export TF_LOG=DEBUG
SRC=$(shell find . -name '*.go')

.PHONY: all clean release install

all: 
	go build -o terraform-provider-sonarqube

run: 
	go build -o terraform-provider-sonarqube
	terraform init
	terraform apply --auto-approve

clean:
	rm -rf terraform-provider-sonarqube .terraform terraform.tfstate crash.log terraform.tfstate.backup
