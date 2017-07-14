export CGO_ENABLED:=0

# kernel-style V=1 build verbosity
ifeq ("$(origin V)", "command line")
       BUILD_VERBOSE = $(V)
endif

ifeq ($(BUILD_VERBOSE),1)
       Q =
else
       Q = @
endif

VERSION=$(shell git describe --dirty)
LD_FLAGS="-w -X github.com/coreos/container-linux-config-transpiler/internal/version.Raw=$(VERSION)"

REPO=github.com/coreos/container-linux-config-transpiler

all: build

build: bin/ct

bin/ct:
	$(Q)go build -o $@ -v -ldflags $(LD_FLAGS) $(REPO)/internal

test:
	$(Q)./test

.PHONY: vendor
vendor:
	$(Q)glide update --strip-vendor
	$(Q)glide-vc --use-lock-file --no-tests --only-code

clean:
	$(Q)rm -rf bin
	$(Q)rm -rf _output

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
	$(Q)$(GOARGS) go build -o $(NAME) -ldflags $(LD_FLAGS) $(REPO)/internal

.PHONY: all build clean test
