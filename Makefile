APP_NAME = go-project-startup
APP_DESC = go project startup template
BASE_PKG = github.com/zunkk

APP_PKG = $(BASE_PKG)/$(APP_NAME)
CONFIG_PKG = $(APP_PKG)/pkg/config
APP_START_DIR = cmd/app
PROJECT_PATH := $(shell pwd)

GO_BIN = go
ifneq (${GO},)
	GO_BIN = ${GO}
endif

BUILD_TIME = $(shell date +%F-%Z/%T)
COMMIT_ID =
TAG =
ifneq ($(wildcard .git/config),)
    COMMIT_ID = $(shell git rev-parse HEAD)
    TAG = $(shell git describe --abbrev=0 --tag)
endif

ifeq ($(version),)
	# not specify version: make install
	APP_VERSION = $(TAG)
	ifeq ($(APP_VERSION),)
		APP_VERSION = dev
	endif
else
	# specify version: make install version=v0.6.1-dev
	APP_VERSION = $(version)
endif

LDFLAGS = -X "${CONFIG_PKG}.Version=${APP_VERSION}"
LDFLAGS += -X "${CONFIG_PKG}.BuildTime=${BUILD_TIME}"
LDFLAGS += -X "${CONFIG_PKG}.CommitID=${COMMIT_ID}"
LDFLAGS += -X "${CONFIG_PKG}.AppName=${APP_NAME}"
LDFLAGS += -X "${CONFIG_PKG}.AppDesc=${APP_DESC}"


ifeq ($(target),)
	COMPILE_TARGET=
else
    PARAMS=$(subst -, ,$(target))
    ifeq ($(words $(PARAMS)),2)
    	OS=$(word 1, $(PARAMS))
    	ARCH=$(word 2, $(PARAMS))
    	COMPILE_TARGET=CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH)
    else
        $(error error param: '$(target)'! example: 'target=darwin-amd64')
    endif
endif

.PHONY: init install build dev-build lint fmt test precommit compile-network-pb compile-grpc-pb

# Init subModule
init:
	${GO_BIN} install github.com/fsgo/go_fmt@v0.4.13
	${GO_BIN} install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2
	./scripts/update_go_pkg.sh $(PROJECT_PATH) github.com/zunkk/go-project-startup $(APP_PKG)


build:
	${GO_BIN} env -w CGO_LDFLAGS=""
	cd ${APP_START_DIR}  && $(COMPILE_TARGET) ${GO_BIN} build -ldflags '-s -w $(LDFLAGS)' -o app-${APP_VERSION}

# Check and print out style mistakes
lint:
	golangci-lint run --timeout=5m -v

# Formats go source code
fmt:
	go_fmt -local $(BASE_PKG) -mi

# Test unit tests of source code
test:
	${GO_BIN} test ./...

package:build
	cd ../../
	cp ./${APP_START_DIR}/app-${APP_VERSION} ./deploy/tools/bin/app
	tar czvf ./app-${APP_VERSION}.tar.gz -C ./deploy/ .

dev-package:build
	cd ../../
	cp ./${APP_START_DIR}/app-${APP_VERSION} ./deploy/tools/bin/app

reset-go-pkg:
	./scripts/update_go_pkg.sh $(PROJECT_PATH) $(APP_PKG) github.com/zunkk/go-project-startup