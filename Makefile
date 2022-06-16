MODULE:=baetyl-cloud
SRC_FILES:=$(shell find . -type f -name '*.go')

export DOCKER_CLI_EXPERIMENTAL=enabled

GIT_TAG:=$(shell git tag --contains HEAD|awk 'END {print}')
GIT_REV:=git-$(shell git rev-parse --short HEAD)
VERSION:=$(if $(GIT_TAG),$(GIT_TAG),$(GIT_REV))

GO_FLAGS:=-s -w -X github.com/baetyl/baetyl-go/v2/utils.REVISION=$(GIT_REV) -X github.com/baetyl/baetyl-go/v2/utils.VERSION=$(VERSION)
GO_PROXY:=https://goproxy.cn
GO_TEST_FLAGS:=-race -short -covermode=atomic -coverprofile=coverage.txt
GO_TEST_PKGS:=$(shell go list ./...)

REGISTRY?=

.PHONY: all
all: build

.PHONY: build
build: $(SRC_FILES)
	env GO111MODULE=on GOPROXY=$(GO_PROXY) CGO_ENABLED=0 go build -o output/$(MODULE) -ldflags "$(GO_FLAGS)" .

.PHONY: image
image:
	@echo "BUILDX: $(REGISTRY)$(MODULE):$(VERSION)"
	@-docker buildx create --name $(MODULE)
	@docker buildx use $(MODULE)
	docker buildx build --push \
		--platform linux/amd64,linux/arm64,linux/arm/v7 \
		-t $(REGISTRY)$(MODULE):$(VERSION) \
		--build-arg GOPROXY="$(GO_PROXY)" \
		--build-arg GIT_REV="$(GIT_REV)" \
		--build-arg VERSION="$(VERSION)" \
		-f Dockerfile .

.PHONY: test
test: fmt
	@go test ${GO_TEST_FLAGS} ${GO_TEST_PKGS}
	@go tool cover -func=coverage.txt | grep total

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: clean
clean:
	@rm -rf output
