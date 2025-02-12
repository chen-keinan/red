GOCMD=go

all:
	$(info  "completed running make file for devcli")
fmt:
	@go fmt ./...
build:
	go build -o devcli
tidy:
	$(GOMOD) tidy -v

.PHONY: fmt build