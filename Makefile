# supress output, run `make XXX V=` to be verbose
V := @

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

COSMOS_PKG_VERSION := $(shell go list -m github.com/cosmos/cosmos-sdk | sed 's:.* ::')

DOCKER := $(shell which docker)

LINTER_NAME := golangci-lint
LINTER_VERSION := v1.51.2

CONTAINER_PROTO_VERSION=0.13.0
CONTAINER_PROTO_IMAGE=ghcr.io/cosmos/proto-builder:$(CONTAINER_PROTO_VERSION)

GOSRC :=  $(shell go env GOPATH)/src
GOBIN := $(shell go env GOPATH)/bin

OUT_DIR := ./build
BIN_NAME := scorumd
BIN_MAIN_PKG := ./cmd/$(BIN_NAME)

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=scorum \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=scorumd \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

.PHONY: install
install: go.sum
	@echo "Install $(BIN_MAIN_PKG)"
	$(V)go install -mod=readonly $(BUILD_FLAGS) $(BIN_MAIN_PKG)
	@echo "DONE"


OS := $(shell go env | grep GOOS | sed -E 's:.*="(.*)":\1:')
ARCH := $(shell go env | grep GOARCH | sed -E 's:.*="(.*)":\1:')

.PHONY: build
build:
	@$(eval BIN_POSTFIX=$(if $(filter $(OS),windows),.exe,))
	@$(eval BIN_OUT := $(OUT_DIR)/$(BIN_NAME)-$(OS)-$(ARCH)$(BIN_POSTFIX))
	@echo "Build $(BIN_OUT)"
	$(V)GOOS=$(OS) GOARCH=$(ARCH) go build -mod=readonly $(BUILD_FLAGS) -o $(BIN_OUT) $(BIN_MAIN_PKG)
	@echo "DONE"

.PHONY: linux
linux: OS = linux
linux: ARCH = amd64
linux: go.sum
linux: build

.PHONE: darwin
darwin: OS = darwin
darwin: ARCH = amd64
darwin: go.sum
darwin: build

.PHONY: windows
windows: OS = windows
windows: ARCH = amd64
windows: go.sum
windows: build

.PHONY: clean
clean:
	@echo "Clean build"
	$(V)rm -rf build/

.PHONY: test
test:
	@echo "Running tests"
	$(V)go test -mod=readonly $(shell go list ./... | grep -v '/simulation')

.PHONY: test-determinism
test-determinism:
	@echo "Running simulation"
	$(V)go test -mod=readonly --tags=simulation -v -run TestAppStateDeterminism ./app

.PHONY: lint
lint:
	@echo "Running linter"
	$(V)$(LINTER_NAME) run --timeout 3m --config .golangci.yml

.PHONY: go.sum
go.sum:
	@echo "Ensure dependencies have not been modified"
	$(V)go mod verify

.PHONY: vendor
vendor:
	$(V)go mod tidy
	$(V)go mod vendor
	@echo "DONE"


.PHONY: generate
generate: generate-proto generate-proto-swagger

.PHONY: generate-proto
generate-proto:
	@echo "Generating Protobuf"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(CONTAINER_PROTO_IMAGE) \
		sh ./scripts/protocgen.sh;

.PHONY: generate-proto-swagger
generate-proto-swagger: SWAGGER_DIR := $(shell mktemp -d)
generate-proto-swagger: SWAGGER_FILES := $$(find $(SWAGGER_DIR) -name "*.swagger.json")
generate-proto-swagger: COSMOS_SDK_DIR := ${GOSRC}/github.com/cosmos/cosmos-sdk
generate-proto-swagger:
	@echo "Generating Protobuf Swagger"
	./scripts/protocgen-swagger-gen.sh

.PHONY: install-proto
install-proto: COSMOS_PROTO_REPO := $(shell mktemp -d)
install-proto:
#   replace when cosmos will be updated to 0.47
#	$(V) go install github.com/cosmos/gogoproto/protoc-gen-gocosmos@v1.4.10

	$(V)git clone git@github.com:regen-network/cosmos-proto.git $(COSMOS_PROTO_REPO) --depth 1 --branch v0.3.1
	$(V)cd $(COSMOS_PROTO_REPO) && go install ./protoc-gen-gocosmos/...

	$(V) go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0

	$(V) go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.16.0
	$(V) cd /tmp && go get -u github.com/g3co/go-swagger-merger


.PHONY: install-linter
install-linter: LINTER_INSTALL_PATH := $(GOBIN)/$(LINTER_NAME)
install-linter:
	@echo INSTALLING $(LINTER_INSTALL_PATH) $(LINTER_VERSION)
	$(V)curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		sh -s -- -b $(GOBIN) $(LINTER_VERSION)
	@echo DONE

.PHONY: check-linter-version
check-linter-version: ACTUAL_LINTER_VERSION := $(shell $(LINTER_NAME) --version 2>/dev/null | awk '{print $$4}')
check-linter-version:
	$(V)[ -z $(ACTUAL_LINTER_VERSION) ] && \
	 echo 'Linter is not installed, run `make install-linter`' && \
	 exit 1 || true

	$(V)if [ v$(ACTUAL_LINTER_VERSION) != $(LINTER_VERSION) ] ; then \
		echo $(LINTER_NAME) is version v$(ACTUAL_LINTER_VERSION), want $(LINTER_VERSION) ; \
		echo 'Make sure $$GOBIN has precedence in $$PATH and' \
		'run `make install-linter` to install the correct version' ; \
        exit 1 ; \
	fi


BIN_POSTFIX := $(if $(filter $(OS),windows),.exe,)
SCORUMD := $(OUT_DIR)/$(BIN_NAME)-$(OS)-$(ARCH)$(BIN_POSTFIX) --home test

.PHONY: local-init
local-init: build
local-init:
	@echo removing test directory
	$(V)rm -rf test

	@echo initalizing test chain
	$(V)$(SCORUMD) init --staking-bond-denom nsp local-network

	@echo adding key 'test'
	$(V)$(SCORUMD) keys add test --keyring-backend test --keyring-dir test

	@echo adding genesis account
	$(V)$(SCORUMD) add-genesis-account --keyring-backend test test 1000000000nscr,1000000000nsp

	@echo adding genesis supervisor
	$(V)$(SCORUMD) add-genesis-supervisor --keyring-backend test test

	@echo creating gentx
	$(V)$(SCORUMD) gentx --keyring-backend test test 1000000nsp

	@echo collecting gentx
	$(V)$(SCORUMD) collect-gentxs

	@echo replace stake with nsp
	sed -i -e 's/"stake"/"nsp"/g' test/config/genesis.json

	@echo validate genesis
	$(V)$(SCORUMD) validate-genesis

	@echo done
	@echo test node home is in ./test directory

.PHONY: local-start
local-start: build
local-start:
	$(V)if [ ! -d "test" ]; then \
            echo "error: run local-init before local-start"; exit 1; \
	fi

	@echo starting test node
	$(V)$(SCORUMD) start

.PHONY: local-reset
local-reset: build
local-reset:
	$(V)if [ ! -d "test" ]; then \
            echo "error: run local-init before local-start"; exit 1; \
	fi

	@echo resetting test node state
	$(V)$(SCORUMD) tendermint unsafe-reset-all


.PHONY: local-image
local-image:
	$(V)docker build -t $(BIN_NAME)-local -f Dockerfile .
