export GO111MODULE=on
export TF_LOG=DEBUG
SRC=$(shell find . -name '*.go')

.PHONY: all vet build test

all: fmt vet build

build:
	go build -a -tags netgo -o terraform-provider-sonarqube
	mkdir -p ~/.terraform.d/plugins/github.com/jdamata/sonarqube/0.1/linux_amd64/
	cp terraform-provider-sonarqube ~/.terraform.d/plugins/github.com/jdamata/sonarqube/0.1/linux_amd64/

fmt:
	go fmt ./...

vet:
	go vet ./...

testacc:
	TF_ACC=1 SONAR_HOST=localhost:9000 SONAR_USER=admin SONAR_PASS=admin go test -cover ./...