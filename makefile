GOCMD=go

all:
	$(info  "completed running make file for devcli")
fmt:
	@go fmt ./...
build:
	go build -o devcli
install:
	GOBIN=/usr/local/bin/ go install
tidy:
	$(GOMOD) tidy -v

.PHONY: fmt build