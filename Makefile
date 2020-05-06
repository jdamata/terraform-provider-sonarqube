export GO111MODULE=on
SRC=$(shell find . -name '*.go')

.PHONY: all clean release install

all: 
	go build -o terraform-provider-sonarqube

clean:
	rm -f terraform-provider-sonarqube
