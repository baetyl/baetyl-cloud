MODULE:=cloud
BIN:=baetyl-$(MODULE)
SRC_FILES:=$(shell find . -type f -name '*.go')
PLATFORM_ALL:=darwin/amd64 linux/amd64 linux/arm64 linux/arm/v7

HOMEDIR := $(shell pwd)
OUTDIR  := $(HOMEDIR)/output

GIT_TAG:=$(shell git tag --contains HEAD)
GIT_REV:=git-$(shell git rev-parse --short HEAD)
VERSION:=$(if $(GIT_TAG),$(GIT_TAG),$(GIT_REV))

GO       = go
GO_MOD   = $(GO) mod
GO_ENV   = env CGO_ENABLED=0
GO_FLAGX = -ldflags '-X "github.com/baetyl/baetyl-go/utils.REVISION=$(GIT_REV)" -X "github.com/baetyl/baetyl-go/utils.VERSION=$(VERSION)"'
GO_BUILD = $(GO_ENV) $(GO) build $(GO_FLAGX) $(GO_TAGS)
GOTEST   = $(GO) test
GOPKGS   = $$($(GO) list ./... | grep -vE "vendor")

ifndef PLATFORMS
	GO_OS:=$(shell go env GOOS)
	GO_ARCH:=$(shell go env GOARCH)
	GO_ARM:=$(shell go env GOARM)
	PLATFORMS:=$(if $(GO_ARM),$(GO_OS)/$(GO_ARCH)/$(GO_ARM),$(GO_OS)/$(GO_ARCH))
	ifeq ($(GO_OS),darwin)
		PLATFORMS+=linux/amd64
	endif
else ifeq ($(PLATFORMS),all)
	override PLATFORMS:=$(PLATFORM_ALL)
endif

REGISTRY?=
XFLAGS?=--load
XPLATFORMS:=$(shell echo $(filter-out darwin/amd64,$(PLATFORMS)) | sed 's: :,:g')

.PHONY: all
all: $(SRC_FILES)
	@echo "BUILD $(BIN)"
	@env GO111MODULE=on GOPROXY=https://goproxy.cn CGO_ENABLED=0 go build -o $(BIN) $(GO_FLAGS) .

test: fmt test-case
test-case: set-env
	$(GOTEST) -race -cover -coverprofile=coverage.out $(GOPKGS)

clean:
	rm -rf $(OUTDIR)
	rm -rf $(HOMEDIR)/$(MODULE)

image:
	@echo "BUILDX: $(REGISTRY)$(MODULE):$(VERSION)"
	@-docker buildx create --name baetyl
	@docker buildx use baetyl
	@docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	docker buildx build $(XFLAGS) --platform $(XPLATFORMS) -t $(REGISTRY)$(MODULE):$(VERSION) -f Dockerfile .

fmt:
	go fmt ./...

.PHONY: all test clean image
