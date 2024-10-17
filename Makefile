export GO111MODULE=on
export TF_LOG=DEBUG
SRC=$(shell find . -name '*.go')
SONARQUBE_IMAGE?=sonarqube:latest
SONARQUBE_START_SLEEP?=60
GO_VER ?= go

.PHONY: all vet build test

all: fmt vet build

tools:
	@if ! [ -x "$$(command -v tfplugindocs)" ]; then \
		echo "Installing tfplugindocs..."; \
		$(GO_VER) get github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs && $(GO_VER) install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs; \
	else \
		echo "tfplugindocs is already installed."; \
	fi

docs:
	rm -f docs/data-sources/*.md
	rm -f docs/resources/*.md
	rm -f docs/*.md
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
