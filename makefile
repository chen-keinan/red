SHELL := /bin/bash

GOCMD=go
GOMOCKS=$(GOCMD) generate ./...
GOMOD=$(GOCMD) mod
GOTEST=$(GOCMD) test


all:
	$(info  "completed running make file for go-simple-config")
fmt:
	@go fmt ./...
build:
	go build -o dev-cli
tidy:
	$(GOMOD) tidy -v

.PHONY: fmt build