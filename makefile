GOCMD=go

fmt:
	@go fmt ./...
build:
	go build -o red
install:
	GOBIN=/usr/local/bin/ go install
tidy:
	$(GOMOD) tidy -v

.PHONY: fmt build