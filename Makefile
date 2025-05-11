export GO111MODULE=on
export CGO_ENABLED=1
export TF_LOG=DEBUG
SRC=$(shell find . -name '*.go')
SONARQUBE_IMAGE?=sonarqube:latest
SONARQUBE_START_SLEEP?=60
GO_VER ?= go

.PHONY: all vet build test tools docs

all: fmt vet build

tools:
	cd tools && $(GO_VER) install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

docs:
	@tfplugindocs generate

build:
	go build -a -tags netgo -o terraform-provider-sonarqube

fmt:
	go fmt ./...

vet:
	go vet ./...

testacc:
	docker run --name sonarqube1 -d -p 9001:9000 ${SONARQUBE_IMAGE}
	timeout 300 bash -c 'while [[ "$$(curl -s -o /dev/null -w "%{http_code}" admin:admin@localhost:9001/api/system/info)" != "200" ]]; do echo "waiting for SonarQube to start"; sleep 15; done'
	-TF_ACC=1 SONAR_HOST=http://localhost:9001 SONAR_USER=admin SONAR_PASS=admin go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	docker stop sonarqube1
	docker rm sonarqube1

testacc-no-docker-windows:
#	docker run --name sonarqube1 -d -p 9001:9000 sonarqube:latest
	set TF_ACC=1\
	&& set SONAR_HOST=http://localhost:9001\
	&& set SONAR_USER=admin\
	&& set SONAR_PASS=admin\
	&& go test -race -coverprofile=coverage.txt -covermode=atomic ./... -run TestAccSonarqubeUserLocal