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

# Run sonarqube locally via docker -> docker run -d -p 9000:9000 sonarqube:latest
# Export these for local testing   -> export SONAR_HOST=localhost:9000 SONAR_USER=admin SONAR_PASS=admin
testacc:
	TF_ACC=1 go test -race -coverprofile=coverage.txt -covermode=atomic -cover ./...