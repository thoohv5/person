GO				:= go
PROTOC			:= protoc
PREFIX			:= $(shell pwd)
PROJECT_DIR		:= $(shell dirname $(PREFIX))
PROJECT_NAME	:= $(shell basename $(PREFIX))
PROJECT_NAME_VAR:= $(shell echo $(PROJECT_NAME) | tr '-' '_')

GOPATH	:= $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))

pkgs := ./...
test-pkgs = $(shell go list ./... | grep -v -E '/api/|/cmd/|/internal/server|/mocks')
lint-pkgs = $(shell go list -f '{{.Dir}}' ./... | grep -v -E '/mocks|/api/conf' | sed 's/$$/\/.../')
gofmt-files := $(shell go list -f '{{.Dir}}' ./... | grep -v -E '/api/docs|/api/conf|/api/translate')

# 区分系统
ifeq ($(OS),Windows_NT)
 PLATFORM=Windows
else
 ifeq ($(shell uname),Darwin)
  PLATFORM=Darwin
 else
  PLATFORM=Linux
 endif
endif

ifeq ($(findstring MINGW,$(shell uname -s)),MINGW)
	override gofmt-files := $(subst \,/,$(gofmt-files))
	override PREFIX := $(subst \,/,$(PREFIX))
	override GOPATH := $(subst \,/,$(shell $(GO) env GOPATH))
endif

GOTEST_DIR := test-results
test-flags := -v -parallel 1 -mod readonly
ifeq ($(GOHOSTARCH),amd64)
	ifeq ($(GOHOSTOS),$(filter $(GOHOSTOS),linux freebsd darwin windows))
		# Only supported on amd64
		test-flags := $(test-flags) -race
	endif
endif

GO_MOD  := $(shell $(GO) env GOMOD)
APP_DIR := $(shell dirname $(GO_MOD))
APP_PROJECT_DIR := $(shell dirname $(APP_DIR))
VERSION := $(shell git describe --tags --always)
API_PROTO_FILES      := $(shell cd $(PROJECT_DIR); [ -d "$(PREFIX)/api" ] && find $(PREFIX)/api -name *.proto)
API_PB_GO_FILES      := $(shell cd $(PROJECT_DIR); [ -d "$(PREFIX)/api" ] && find $(PREFIX)/api -name *.pb.go)
API_YAML_FILES       := $(shell cd $(PROJECT_DIR); [ -d "$(PREFIX)/api" ] && find $(PREFIX)/api -name *.yaml)
INTERNAL_GO_FILES    := $(shell cd $(PROJECT_DIR); [ -d "$(PREFIX)/internal" ] && find $(PREFIX)/internal -name *.go | grep -v -E '/conf|/server')
RPC_GO_FILES         := $(shell cd $(PROJECT_DIR); [ -d "$(PREFIX)/api/http/rpc" ] && (find $(PREFIX)/api/http/rpc -name *.go | grep -v -E '/*_test'))


GOLANGCI_LINT			?= $(GOPATH)/bin/golangci-lint
GOLANGCI_LINT_OPTS		?=
GOIMPORTS               ?= $(GOPATH)/bin/goimports
GOSWAG             		?= $(GOPATH)/bin/swag
GOWIRE             	    ?= $(GOPATH)/bin/wire
MOCKGEN             	?= $(GOPATH)/bin/mockgen
REGISTERFIELD           ?= $(GOPATH)/bin/field

$(GOLANGCI_LINT): $(GO_MOD)
	@echo "> installing golangci-lint"
	@$(GO) install "github.com/golangci/golangci-lint/cmd/golangci-lint"

$(GOIMPORTS): $(GO_MOD)
	@echo "> installing goimports"
	@$(GO) install "golang.org/x/tools/cmd/goimports"

$(GOSWAG): $(GO_MOD)
	@echo "> installing swag"
	@$(GO) install "github.com/swaggo/swag/cmd/swag@v1.8.9"

$(GOWIRE): $(GO_MOD)
	@echo "> installing wire"
	@$(GO) install "github.com/google/wire/cmd/wire@v0.5.0"

$(MOCKGEN): $(GO_MOD)
	@echo "> installing mockgen"
	@$(GO) install "github.com/golang/mock/mockgen"

$(REGISTERFIELD):
	@echo "> installing registerfield"
	@pushd $(APP_DIR)/pkg/cmd/field ;$(GO) install ./...; popd

# 提取文件路径的文件名
filename = $(notdir $(1))

# This rule is used to forward a target like "test" to "common-test".  This
# allows a new "build" target to be defined in a Makefile which includes this
# one and override "common-build" without override warnings.
%: common-% ;

.PHONY: common-all
common-all: lint test

.PHONY: common-setup
common-setup:
	@echo ">> Setup environments"
	@$(GO) env -w GO111MODULE=on
	@$(GO) env -w GOPRIVATE=github.com
	@$(GO) env -w GOPROXY=https://goproxy.cn,direct
	@$(GO) mod download

# show help
.PHONY: common-help
common-help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.PHONY: common-test
common-test:
	@echo ">> running tests"
	@$(GO) test $(test-flags) $(test-pkgs)


$(GOTEST_DIR):
	@mkdir -p $@


.PHONY: common-gen
common-gen:
	@#$(GO) generate ./...
	@rm -rf $(PREFIX)/api/constant/enum/*_string.go
	@$(GO) generate $(shell $(GO) list ./... | grep -v -E '/translate')
	@echo "\033[32m >> ${PROJECT_NAME}代码依赖已生成 \033[0m"

.PHONY: common-gen-all
common-gen-all:
	@$(foreach dir, $(shell ls -d $(APP_DIR)/app/interface/*), pushd $(dir) && make gen && popd;)


.PHONY: common-lint
common-lint: $(GOLANGCI_LINT)
	@echo ">> linting code"
# 'go list' needs to be executed before staticcheck to prepopulate the modules cache.
# Otherwise staticcheck might fail randomly for some reason not yet explained.
	@$(GO) list -e -compiled -test=true -export=false -deps=true -find=false -tags= -- ./... > /dev/null
	@$(GOLANGCI_LINT) run $(GOLANGCI_LINT_OPTS) $(lint-pkgs)

.PHONY: common-yaml
common-yaml:
	@$(foreach file, $(shell find $(PREFIX)/configs -name "*.yaml.example"), cp $(file) $(file:.example=);)
	@echo "\033[32m >> ${PROJECT_NAME}复制配置已完成 \033[0m"

.PHONY: common-format
common-format: $(GOIMPORTS)
	@echo ">> formatting code"
	@$(GOIMPORTS) -local "github.com" -w $(gofmt-files)
	@echo "\033[32m >> ${PROJECT_NAME}代码格式化已完成 \033[0m"

.PHONY: common-wire
common-wire: $(GOWIRE)
	@echo ">> wire exec"
	@pushd $(PREFIX)/cmd/command && $(GOWIRE) && popd
	@echo "\033[32m >> ${PROJECT_NAME}模块注入已完成 \033[0m"

.PHONY: common-mock
common-mock: $(MOCKGEN)
	@echo ">> interface mock"
	@$(MOCKGEN) -destination=mocks/pg_IDB.go -package=mocks github.com/go-pg/pg/v10 DBI
	@$(MOCKGEN) -destination=mocks/pg_orm_DB.go -package=mocks github.com/go-pg/pg/v10/orm DB
	@$(MOCKGEN) -destination=mocks/logger.go  -package=mocks -source=$(PREFIX)/../../../pkg/log/log.go
	@$(MOCKGEN) -destination=mocks/config.go  -package=mocks -source=$(PREFIX)/api/conf/conf_configs.pb.go
	@$(foreach file, ${RPC_GO_FILES}, $(MOCKGEN) -destination=$(shell echo $(file)|sed 's/$(PROJECT_NAME)\/api\/http/$(PROJECT_NAME)\/mocks/') -source=$(file);)
	@$(foreach file, ${INTERNAL_GO_FILES}, $(MOCKGEN) -destination=$(shell echo $(file)|sed 's/$(PROJECT_NAME)\/internal/$(PROJECT_NAME)\/mocks/') -source=$(file);)

.PHONY: common-mock-all
common-mock-all:
	@$(foreach dir, $(shell ls -d $(APP_DIR)/app/interface/*), pushd $(dir) && make mock && popd;)


.PHONY: common-swag
common-swag: $(GOSWAG) $(REGISTERFIELD)
	echo $PLATFORM;
	@$(GOSWAG) init --instanceName $(PROJECT_NAME_VAR) -d $(PREFIX)/internal/controller,$(PREFIX)/api/constant,$(PREFIX)/api/http/request,$(PREFIX)/api/http/response,$(APP_DIR)/internal/http -g ../../cmd/$(PROJECT_NAME)/main.go -o $(PREFIX)/api/docs
	@${GOSWAG} fmt
	@$(REGISTERFIELD) -d $(PREFIX)/api/http/request
	@$(MAKE) format
	@echo "\033[32m >> ${PROJECT_NAME}接口文档已生成 \033[0m"

.PHONY: common-migration
common-migration:
	@$(GO) run $(PREFIX)/cmd/$(PROJECT_NAME)/main.go migration reset
	@$(GO) run $(PREFIX)/cmd/$(PROJECT_NAME)/main.go migration up