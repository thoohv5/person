include $(shell dirname $(shell go env GOMOD))/Makefile
UPX := upx


.PHONY: wire
wire: $(GOWIRE)
	@echo ">> wire exec"
	@pushd $(PREFIX)/cmd/command && $(GOWIRE) && popd
	@pushd $(PREFIX)/internal/service && $(GOWIRE) && popd
	@echo "\033[32m >> ${PROJECT_NAME}模块注入已完成 \033[0m"

.PHONY: build
# build
build:
	@$(MAKE) yaml
	@$(MAKE) gen-all
	@$(MAKE) wire
	@$(MAKE) format
	@$(MAKE) swag
	@$(MAKE) mock
	@rm -rf build && mkdir -p build && $(build_env) $(GO) build -ldflags "-s -w -X main.version=$(VERSION)" -o build ./...
	@$(UPX) -9 -o build/$(PROJECT_NAME)_tmp build/$(PROJECT_NAME) && mv build/$(PROJECT_NAME)_tmp build/$(PROJECT_NAME)
	@echo "\033[32m >> 编译成功 \033[0m"


.PHONY: run
# run
run:
	@$(MAKE) gen-all
	@$(MAKE) wire
	@$(MAKE) format
	@$(MAKE) swag
	@echo "\033[32m >> $(PROJECT_NAME)项目已启动 \033[0m"
	@$(GO) run "$(PREFIX)/cmd/$(PROJECT_NAME)" server --conf "$(PREFIX)/configs/"


.PHONY: swag
swag: $(GOSWAG) $(REGISTERFIELD)
	@$(GOSWAG) init --instanceName $(PROJECT_NAME_VAR) -d $(PREFIX)/internal/controller,$(PREFIX)/internal/repository,$(PREFIX)/api/constant,$(PREFIX)/api/http/request,$(PREFIX)/api/http/response,$(APP_DIR)/internal/provide/http -g ../../cmd/$(PROJECT_NAME)/main.go -o $(PREFIX)/api/docs
	@$(GOSWAG) fmt
	@$(REGISTERFIELD) -d $(PREFIX)/api/http/request
	@$(MAKE) format
	@echo "\033[32m >> $(PROJECT_NAME)文档已生成 \033[0m"


.PHONY: test
test:
	@echo ">> running tests"
	@#$(GO) test $(test-flags) $(test-pkgs)