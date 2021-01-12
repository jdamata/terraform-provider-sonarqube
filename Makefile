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
	docker run --name sonarqube1 -d -p 9001:9000 sonarqube:latest
	sleep 45
	-TF_ACC=1 SONAR_HOST=http://localhost:9001 SONAR_USER=admin SONAR_PASS=admin go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	docker stop sonarqube1
	docker rm sonarqube1