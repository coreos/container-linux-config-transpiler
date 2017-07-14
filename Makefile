export CGO_ENABLED:=0

VERSION=$(shell git describe --dirty)
LD_FLAGS="-w -X github.com/coreos/container-linux-config-transpiler/internal/version.Raw=$(VERSION)"

REPO=github.com/coreos/container-linux-config-transpiler

all: build

build: bin/ct

bin/ct:
	@go build -o $@ -v -ldflags $(LD_FLAGS) $(REPO)/internal

test:
	@./test

.PHONY: vendor
vendor:
	@glide update --strip-vendor
	@glide-vc --use-lock-file --no-tests --only-code

clean:
	@rm -rf bin
	@rm -rf _output

.PHONY: release
release: \
	_output/unknown-linux-gnu/ct \
	_output/apple-darwin/ct \
	_output/pc-windows-gnu/ct

_output/unknown-linux-gnu/ct: GOARGS = GOOS=linux GOARCH=amd64
_output/apple-darwin/ct: GOARGS = GOOS=darwin GOARCH=amd64
_output/pc-windows-gnu/ct: GOARGS = GOOS=windows GOARCH=amd64

_output/%/ct: NAME=_output/ct-$(VERSION)-x86_64-$*
_output/%/ct:
	$(GOARGS) go build -o $(NAME) -ldflags $(LD_FLAGS) $(REPO)/internal

.PHONY: all build clean test
