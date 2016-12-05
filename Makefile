.PHONY: test clean

REPO	= github.com/akaspin/terraform-provider-ansible
BIN		= terraform-provider-ansible

BENCH	= .
TESTS	= .


CWD 		= $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VENDOR 		= $(CWD)/_vendor
SRC 		= $(shell find . -type f \( -iname '*.go' ! -iname "*_test.go" \) -not -path "./_vendor/*")
SRC_TEST 	= $(shell find . -type f -name '*_test.go' -not -path "./_vendor/*")
SRC_VENDOR 	= $(shell find ./_vendor -type f \( -iname '*.go' ! -iname "*_test.go" \))
PKGS 		= $(shell cd $(GOPATH)/src/$(REPO) && go list ./...)

V = $(shell git describe --always --tags --dirty)
GOOPTS = -installsuffix cgo -ldflags '-s -X $(REPO)/command.V=$(V)'

ifdef GOBIN
	INSTALL_DIR=$(GOBIN)
else
    INSTALL_DIR=$(GOPATH)/bin
endif

###
###	Install
###

install: $(INSTALL_DIR)/$(BIN)
install-debug: $(INSTALL_DIR)/$(BIN)-debug

$(INSTALL_DIR)/$(BIN): $(SRC) $(SRC_VENDOR)
	CGO_ENABLED=0 go build $(GOOPTS) -o $@ $(REPO)/command/$(BIN)

$(INSTALL_DIR)/$(BIN)-debug: $(SRC) $(SRC_VENDOR)
	GOPATH=$(VENDOR):$(GOPATH) CGO_ENABLED=0 go build $(GOOPTS) -tags debug -o $@ $(REPO)/command/$(BIN)

uninstall:
	rm -rf $(INSTALL_DIR)/$(BIN)
	rm -rf $(INSTALL_DIR)/$(BIN)-debug

###
### Deps
###

deps: $(VENDOR)/lock

$(VENDOR)/lock: vendor.conf
	trash --target _vendor/src
	touch $@

clean-deps:
	rm -rf $(VENDOR)

###
### clean
###

clean: clean-dist
