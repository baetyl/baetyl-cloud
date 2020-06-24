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

MODULE  := baetyl-cloud

all: test package

prepare: prepare-dep
prepare-dep:
	git config --global http.sslVerify false
	$(AGILE_BCLOUD)
	$(EASYPACK_ANTI_DEBUG)
	$(EASYPACK_AIPE_SECURITY)

set-env:
	$(GO) env -w GONOPROXY=\*\*.baidu.com\*\*
	$(GO) env -w GOPROXY=https://goproxy.baidu.com
	$(GO) env -w GONOSUMDB=\*

compile:build
build: set-env
	$(GO_MOD) tidy
	$(GO_BUILD) -o $(HOMEDIR)/$(MODULE)

test: fmt test-case
test-case: set-env
	$(GOTEST) -race -cover -coverprofile=coverage.out $(GOPKGS)

package: compile package-bin
package-bin:
	mkdir -p $(OUTDIR)/bin
	mv $(MODULE) $(OUTDIR)/bin/$(MODULE)

clean:
	rm -rf $(OUTDIR)
	rm -rf $(HOMEDIR)/$(MODULE)

fmt:
	go fmt ./...

.PHONY: all prepare compile test package clean build
