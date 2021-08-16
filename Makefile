BRANCH := $(shell git rev-parse --symbolic-full-name --abbrev-ref HEAD)
SUMMARY := $(shell bash .ci/version)
VERSION := $(shell cat VERSION)
OS := $(shell uname -s | awk '{print tolower($$0)}')
GOARCH = amd64
PACKAGE_NAME = $(shell head -n1 go.mod | awk '{print $$2}')
NAME := $(shell basename $(PACKAGE_NAME))
LDFLAGS = "-X $(PACKAGE_NAME)/pkg/common.SUMMARY=$(SUMMARY) -X $(PACKAGE_NAME)/pkg/common.BRANCH=$(BRANCH) -X $(PACKAGE_NAME)/pkg/common.VERSION=$(VERSION)"

vendor:
	go mod vendor

release: generate vendor
	@mkdir -p release/
	go build -mod=vendor -ldflags $(LDFLAGS) -o release/$(NAME)_linux_$(GOARCH) main.go

release-all: generate vendor
	@mkdir -p release/
	go build -mod=vendor -ldflags $(LDFLAGS) -o release/$(NAME)_linux_$(GOARCH) main.go
	GOARCH=arm64 go build -mod=vendor -ldflags $(LDFLAGS) -o release/$(NAME)_linux_arm64 main.go
	GOOS=darwin go build -mod=vendor -ldflags $(LDFLAGS) -o release/$(NAME)_darwin_$(GOARCH) main.go
	GOOS=windows go build -mod=vendor -ldflags $(LDFLAGS) -o release/$(NAME)_windows_$(GOARCH).exe main.go

pkg/box/blob.go:
	@go generate ./...

build-ui:
	$(MAKE) -C ui build

build: build-ui generate vendor
	@mkdir -p release/
	@go build -ldflags $(LDFLAGS) -o release/$(NAME)-$(OS)-$(GOARCH) main.go

build-static: generate
	@mkdir -p release/
	@go build -ldflags '-w -extldflags "-static"' -o release/pipeliner-$(OS)-$(GOARCH)-static main.go

generate: pkg/box/blob.go

docker:
	@docker build -t ekristen/pipeliner .

check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	GO111MODULE=on go mod vendor  && GO111MODULE=off swagger generate spec -o ./docs/swagger.yaml --scan-models

serve-swagger: check-swagger

docs-build:
	docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material build

docs-serve:
	docker run --rm -it -p 8000:8000 -v ${PWD}:/docs squidfunk/mkdocs-material

snapshot:
	SUMMARY=$(SUMMARY) VERSION=$(VERSION) BRANCH=$(BRANCH) goreleaser release --snapshot --skip-publish --rm-dist