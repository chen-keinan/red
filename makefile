GOCMD=go

all:
	$(info  "completed running make file for dev-cli")
fmt:
	@go fmt ./...
build:
	go build -o dev-cli
tidy:
	$(GOMOD) tidy -v

.PHONY: fmt build