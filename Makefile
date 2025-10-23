# --- variables ---
GOCMD        ?= go
GOINSTALL    := $(GOCMD) install
BIN          ?= autoglue
MAIN         ?= main.go
UI_DIR       ?= ui
UI_DEST_DIR  ?= internal/ui

# SDK / module settings (Go)
GIT_HOST     ?= github.com
GIT_USER     ?= glueops
SDK_REPO     ?= autoglue-sdk                # repo name used for module path
SDK_OUTDIR   ?= sdk/go                      # output directory (inside repo)
SDK_PKG      ?= autoglue                    # package name inside the SDK

# Go versioning (go.mod uses major.minor; toolchain can include patch)
GO_VERSION   ?= 1.25
GO_TOOLCHAIN ?= go1.25.3

# SDK / package settings (TypeScript)
SDK_TS_OUTDIR     ?= sdk/ts
SDK_TS_GEN        ?= typescript-fetch
SDK_TS_NPM_NAME   ?= @glueops/autoglue-sdk
SDK_TS_NPM_VER    ?= 0.1.0
SDK_TS_DIR        := $(abspath $(SDK_TS_OUTDIR))
SDK_TS_PROPS      ?= supportsES6=true,typescriptThreePlus=true,useSingleRequestParameter=true,withSeparateModelsAndApi=true,modelPropertyNaming=original,enumPropertyNaming=original,useUnionTypes=true

# Path for vendored UI SDK (absolute, path-safe)
SDK_TS_UI_OUTDIR ?= ui/src/sdk
SDK_TS_UI_DIR    := $(abspath $(SDK_TS_UI_OUTDIR))

SWAG         := $(shell command -v swag 2>/dev/null)
GMU          := $(shell command -v go-mod-upgrade 2>/dev/null)
YARN         := $(shell command -v yarn 2>/dev/null)
NPM          := $(shell command -v npm 2>/dev/null)
OGC          := $(shell command -v openapi-generator-cli 2>/dev/null || command -v openapi-generator 2>/dev/null)

.DEFAULT_GOAL := build

.PHONY: all prepare ui-install ui-build ui swagger build clean fmt vet tidy upgrade \
        sdk sdk-go sdk-ts sdk-ts-ui sdk-all worksync wire-sdk-replace help dev

# --- meta targets ---
all: build
prepare: fmt vet tidy upgrade

# --- go hygiene ---
fmt:
	@$(GOCMD) fmt ./...

vet:
	@$(GOCMD) vet ./...

tidy:
	@$(GOCMD) mod tidy

upgrade:
	@echo ">> Checking go-mod-upgrade..."
	@if [ -z "$(GMU)" ]; then \
		echo "Installing go-mod-upgrade..."; \
		$(GOINSTALL) github.com/oligot/go-mod-upgrade@latest; \
	fi
	@go-mod-upgrade -f || true

# --- ui ---
ui-install:
	@echo ">> Installing UI deps in $(UI_DIR)..."
	@if [ -n "$(YARN)" ]; then \
		cd $(UI_DIR) && yarn install --frozen-lockfile; \
	elif [ -n "$(NPM)" ]; then \
		cd $(UI_DIR) && npm ci; \
	else \
		echo "Error: neither yarn nor npm is installed." >&2; exit 1; \
	fi

ui-build: ui-install
	@echo ">> Building UI in $(UI_DIR)..."
	@rm -rf $(UI_DIR)/dist
	@if [ -n "$(YARN)" ]; then \
		cd $(UI_DIR) && yarn build; \
	else \
		cd $(UI_DIR) && npm run build; \
	fi
	@echo ">> Copying UI dist -> $(UI_DEST_DIR)/dist"
	@rm -rf $(UI_DEST_DIR)/dist
	@mkdir -p $(UI_DEST_DIR)
	@cp -R $(UI_DIR)/dist $(UI_DEST_DIR)/dist

ui: ui-build

# --- swagger ---
swagger:
	@echo ">> Generating Swagger docs..."
	@if [ -z "$(SWAG)" ]; then \
		echo "Installing swag..."; \
		$(GOINSTALL) github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@rm -rf docs/swagger.* docs/docs.go
	@swag init -g $(MAIN) -o docs

# --- build ---
build: prepare ui swagger
	@echo ">> Building Go binary: $(BIN)"
	@$(GOCMD) build -o $(BIN) $(MAIN)

# --- development ---
dev: ui-install swagger
	@echo ">> Starting Vite (frontend) and Go API (backend)..."
	@cd $(UI_DIR) && \
	( \
		if command -v yarn >/dev/null 2>&1; then \
			yarn dev & \
		elif command -v npm >/dev/null 2>&1; then \
			npm run dev & \
		else \
			echo "Error: neither yarn nor npm is installed." >&2; exit 1; \
		fi; \
		cd .. && \
		UI_DEV=1 $(GOCMD) run . serve & \
		wait \
	)

# --- sdk generation (Go) ---
sdk-go: swagger
	@echo ">> Generating Go SDK (module $(GIT_HOST)/$(GIT_USER)/$(SDK_REPO), Go $(GO_VERSION) / $(GO_TOOLCHAIN))..."
	@set -e; \
	export GO_POST_PROCESS_FILE="gofmt -w"; \
	if [ -z "$(OGC)" ]; then \
		if [ -z "$(NPM)" ]; then \
			echo "Error: npm is required to install openapi-generator-cli." >&2; exit 1; \
		fi; \
		echo "Installing openapi-generator-cli..."; \
		$(NPM) i -g @openapitools/openapi-generator-cli; \
		OGC_BIN=openapi-generator-cli; \
	else \
		OGC_BIN="$(OGC)"; \
	fi; \
	mkdir -p "$(SDK_OUTDIR)"; \
	"$$OGC_BIN" generate \
		--enable-post-process-file \
		--generate-alias-as-model \
		-i docs/swagger.json \
		-g go \
		-o "$(SDK_OUTDIR)" \
		--additional-properties=packageName=$(SDK_PKG) \
		--git-host "$(GIT_HOST)" \
		--git-user-id "$(GIT_USER)" \
		--git-repo-id "$(SDK_REPO)"; \
	cd "$(SDK_OUTDIR)"; \
	$(GOCMD) mod edit -go=$(GO_VERSION); \
	$(GOCMD) mod edit -toolchain=$(GO_TOOLCHAIN); \
	$(GOCMD) mod tidy

# --- sdk generation (TypeScript) ---
sdk-ts: swagger
	@echo ">> Generating TypeScript SDK ($(SDK_TS_GEN)) -> $(SDK_TS_DIR) as $(SDK_TS_NPM_NAME)@$(SDK_TS_NPM_VER)"
	@set -e; \
	if [ -z "$(OGC)" ]; then \
		if [ -z "$(NPM)" ]; then \
			echo "Error: npm is required to install openapi-generator-cli." >&2; exit 1; \
		fi; \
		echo "Installing openapi-generator-cli..."; \
		$(NPM) i -g @openapitools/openapi-generator-cli; \
		OGC_BIN=openapi-generator-cli; \
	else \
		OGC_BIN="$(OGC)"; \
	fi; \
	rm -rf "$(SDK_TS_DIR)"; \
	mkdir -p "$(SDK_TS_DIR)"; \
	"$$OGC_BIN" generate \
		-i docs/swagger.json \
		-g "$(SDK_TS_GEN)" \
		-o "$(SDK_TS_DIR)" \
		--additional-properties=npmName=$(SDK_TS_NPM_NAME),npmVersion=$(SDK_TS_NPM_VER),$(SDK_TS_PROPS); \
	if [ ! -d "$(SDK_TS_DIR)" ]; then \
		echo "Generation failed: $(SDK_TS_DIR) not found." >&2; exit 1; \
	fi; \
	if command -v npx >/dev/null 2>&1; then \
		echo ">> Prettier: formatting generated TS SDK"; \
		cd "$(SDK_TS_DIR)" && npx --yes prettier -w . || true; \
	fi; \
	echo ">> Installing & building TS SDK in $(SDK_TS_DIR)"; \
	if command -v yarn >/dev/null 2>&1; then \
		cd "$(SDK_TS_DIR)" && yarn install --frozen-lockfile || true; \
		cd "$(SDK_TS_DIR)" && yarn build || true; \
	elif command -v npm >/dev/null 2>&1; then \
		cd "$(SDK_TS_DIR)" && npm ci || npm install || true; \
		cd "$(SDK_TS_DIR)" && npm run build || true; \
	else \
		echo "Warning: neither yarn nor npm is installed; skipping install/build for TS SDK."; \
	fi

# --- sdk generation (TypeScript into UI/src) ---
sdk-ts-ui: swagger
	@echo ">> Generating TypeScript SDK directly into UI source: $(SDK_TS_UI_DIR)"
	@set -e; \
	# resolve generator binary (same as sdk-ts)
	if [ -z "$(OGC)" ]; then \
		if [ -z "$(NPM)" ]; then \
			echo "Error: npm is required to install openapi-generator-cli." >&2; exit 1; \
		fi; \
		echo "Installing openapi-generator-cli..."; \
		$(NPM) i -g @openapitools/openapi-generator-cli; \
		OGC_BIN=openapi-generator-cli; \
	else \
		OGC_BIN="$(OGC)"; \
	fi; \
	rm -rf "$(SDK_TS_UI_DIR)"; \
	mkdir -p "$(SDK_TS_UI_DIR)"; \
	"$$OGC_BIN" generate \
		-i docs/swagger.json \
		-g typescript-fetch \
		-o "$(SDK_TS_UI_DIR)" \
		--additional-properties=$(SDK_TS_PROPS),npmName=$(SDK_TS_NPM_NAME),npmVersion=$(SDK_TS_NPM_VER)

# convenience
sdk-all: sdk-go sdk-ts sdk-ts-ui
sdk: sdk-go

# --- workspace helper (optional with many consumers) ---
worksync:
	@echo ">> (Re)building go.work to include provider, SDK, and all consumers/* modules"
	@set -e; \
	rm -f go.work; \
	$(GOCMD) work init .; \
	if [ -d "$(SDK_OUTDIR)" ]; then $(GOCMD) work use "$(SDK_OUTDIR)"; fi; \
	if [ -d consumers ]; then \
		for d in consumers/*; do \
			if [ -d "$$d" ] && [ -f "$$d/go.mod" ]; then \
				echo "   - adding $$d"; \
				$(GOCMD) work use "$$d"; \
			fi; \
		done; \
	fi; \
	$(GOCMD) work sync

# --- replace helper (use if you don't want go.work) ---
wire-sdk-replace:
	@echo ">> Adding 'replace github.com/glueops/autoglue-sdk => ./$(SDK_OUTDIR)' to any consumers/*/go.mod"
	@set -e; \
	if [ -d consumers ]; then \
		for d in consumers/*; do \
			if [ -d "$$d" ] && [ -f "$$d/go.mod" ]; then \
				if ! grep -q "replace github.com/glueops/autoglue-sdk" "$$d/go.mod"; then \
					echo "replace github.com/glueops/autoglue-sdk => ./../../$(SDK_OUTDIR)" >> "$$d/go.mod"; \
					echo "   - wired $$d"; \
				else \
					echo "   - already wired: $$d"; \
				fi; \
			fi; \
		done; \
	else \
		echo "No consumers/ directory found; skipping."; \
	fi

# --- clean/help ---
clean:
	@echo ">> Cleaning artifacts..."
	@rm -rf "$(BIN)" docs/swagger.* docs/docs.go $(UI_DEST_DIR)/dist $(UI_DIR)/dist $(UI_DIR)/node_modules "$(SDK_OUTDIR)" "$(SDK_TS_OUTDIR)"

help:
	@echo "Targets:"
	@echo "  build           - fmt, vet, tidy, upgrade, build UI, generate Swagger, build Go binary"
	@echo "  ui              - build the Vite UI and copy to $(UI_DEST_DIR)/dist"
	@echo "  swagger         - (re)generate Swagger docs using swag"
	@echo "  sdk-go (sdk)    - generate Go SDK with correct module path and Go/toolchain versions"
	@echo "  sdk-ts          - generate TypeScript SDK (typescript-fetch) with package.json"
	@echo "  sdk-ts-ui       - generate TypeScript SDK directly into ui/src for inline consumption"
	@echo "  sdk-all         - generate both Go and TypeScript SDKs"
	@echo "  worksync        - create/update go.work including provider, SDK, and consumers/*"
	@echo "  wire-sdk-replace- add replace directives in each consumers/*/go.mod"
	@echo "  clean           - remove binary, Swagger outputs, UI dist, and SDKs"
	@echo "  prepare         - fmt, vet, tidy, upgrade deps"
	@echo "  dev             - run Vite UI dev server + Go API with UI_DEV=1"
