APP_NAME = go-project-startup
APP_DESC = go project startup
BASE_PKG = github.com/zunkk
APP_PKG = $(BASE_PKG)/$(APP_NAME)
REPO_PKG = $(APP_PKG)/internal/pkg/config
APP_START_DIR = cmd/$(APP_NAME)
PROJECT_PATH := $(shell pwd)
MODELS_PATH := ${PROJECT_PATH}/internal/core/model
DB_TYPE = psql

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

LDFLAGS = -X "${REPO_PKG}.Version=${APP_VERSION}"
LDFLAGS += -X "${REPO_PKG}.BuildTime=${BUILD_TIME}"
LDFLAGS += -X "${REPO_PKG}.CommitID=${COMMIT_ID}"
LDFLAGS += -X "${REPO_PKG}.AppName=${APP_NAME}"
LDFLAGS += -X "${REPO_PKG}.AppDesc=${APP_DESC}"


RED=\033[0;31m
GREEN=\033[0;32m
BLUE=\033[0;34m
NC=\033[0m

.PHONY: help init lint fmt test test-coverage build package dev-package reset-project-info

help: Makefile
	@printf "${BLUE}Choose a command run:${NC}\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/    /'

## make init: Install dependencies
init:
	${GO_BIN} install go.uber.org/mock/mockgen@main
	${GO_BIN} install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
	${GO_BIN} install mvdan.cc/gofumpt@latest
	${GO_BIN} install golang.org/x/tools/cmd/goimports@latest
	${GO_BIN} install github.com/stephenafamo/bob/gen/bobgen-psql@latest

## make generate-models: Generate db models
generate-models:
	bobgen-psql -c ./build/bobgen.yaml

## make lint: Run golanci-lint
lint:
	golangci-lint run --timeout=5m -v

## make fmt: Formats source code
fmt:
	gofumpt -l -w .
	goimports -local $(BASE_PKG) -w .

## make test: Run go unittest
test:
	${GO_BIN} test -timeout 300s ./... -count=1

## make test-coverage: Test project with cover
test-coverage:
	${GO_BIN} test -timeout 300s -short -coverprofile cover.out -covermode=atomic ${COVERAGE_TEST_PKGS}
	cat cover.out | grep -v "pb.go" >> coverage.txt

## make build: Go build the project
build:
	${GO_BIN} env -w CGO_LDFLAGS=""
	cd ${APP_START_DIR}  && ${GO_BIN} build -ldflags '-s -w $(LDFLAGS)' -o ${APP_NAME}-${APP_VERSION}
	mv ./${APP_START_DIR}/${APP_NAME}-${APP_VERSION} ./

## make package: Package executable binaries and scripts
package:build
	cd ../../
	cp ./${APP_NAME}-${APP_VERSION} ./deploy/tools/bin/${APP_NAME}
	tar czvf ./${APP_NAME}-${APP_VERSION}.tar.gz -C ./deploy/ .

## make dev-package: Compile new executable binary under scripts
dev-package:build
	cd ../../
	rm -f ./deploy/tools/bin/${APP_NAME}
	cp ./${APP_NAME}-${APP_VERSION} ./deploy/tools/bin/${APP_NAME}

## make reset-project-info: Reset project info(name, go package name...)
reset-project-info:
	./scripts/reset_project_info.sh $(PROJECT_PATH) github.com/zunkk/go-project-startup $(APP_PKG) $(APP_NAME)