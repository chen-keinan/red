GOCMD=go

all:
	$(info  "completed running make file for red")
fmt:
	@go fmt ./...
build:
	go build -o red
install:
	GOBIN=/usr/local/bin/ go install
tidy:
	$(GOMOD) tidy -v

.PHONY: fmt build