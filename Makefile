default: build

build:
	CGO_ENABLED=0 go build -trimpath -o $(CURDIR)/bin/oauth2-authenticator $(CURDIR)/cmd/oauth2-authenticator/main.go

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o $(CURDIR)/bin/oauth2-authenticator $(CURDIR)/cmd/oauth2-authenticator/main.go

build_in_container:
	docker run --rm -v "$(CURDIR)":/project -w /project golang:1.14 make build

fmt:
	gofmt -l -w -s cmd/oauth2-authenticator
	gofmt -l -w -s pkg/authenticator

inspect:
	go vet ./...

.PHONY: default build build_linux build_in_container fmt inspect