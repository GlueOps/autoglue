# --- variables ---
GOCMD        ?= go
GOINSTALL    := $(GOCMD) install
BIN          ?= autoglue
MAIN         ?= main.go
UI_DIR       ?= ui
UI_DEST_DIR  ?= internal/web

# Module path (used for ldflags to internal/version)
GIT_HOST     ?= github.com
GIT_USER     ?= glueops
MODULE_PATH  ?= $(GIT_HOST)/$(GIT_USER)/$(BIN)

# SDK / module settings (Go)
SDK_REPO     ?= $(BIN)-sdk                  # repo name used for module path
SDK_OUTDIR   ?= sdk/go                      # output directory (inside repo)
SDK_PKG      ?= ${BIN}                      # package name inside the SDK

UI_SSG_ROUTES ?= /,/login,/docs,/pricing

# Go versioning (go.mod uses major.minor; you’re on 1.25.3)
GO_VERSION   ?= 1.25.3

# SDK / package settings (TypeScript)
SDK_TS_OUTDIR     ?= sdk/ts
SDK_TS_GEN        ?= typescript-fetch
SDK_TS_NPM_NAME   ?= @glueops/$(SDK_REPO)
SDK_TS_NPM_VER    ?= 0.1.0
SDK_TS_DIR        := $(abspath $(SDK_TS_OUTDIR))
SDK_TS_PROPS      ?= supportsES6=true,typescriptThreePlus=true,useSingleRequestParameter=true,withSeparateModelsAndApi=true,modelPropertyNaming=original,enumPropertyNaming=original,useUnionTypes=true
SDK_TS_PROPS_FLAGS := $(foreach p,$(subst , ,$(SDK_TS_PROPS)),-p $(p))

# Path for vendored UI SDK (absolute, path-safe)
SDK_TS_UI_OUTDIR ?= ui/src/sdk
SDK_TS_UI_DIR    := $(abspath $(SDK_TS_UI_OUTDIR))

SWAG         := $(shell command -v swag 2>/dev/null)
GMU          := $(shell command -v go-mod-upgrade 2>/dev/null)
YARN         := $(shell command -v yarn 2>/dev/null)
NPM          := $(shell command -v npm 2>/dev/null)
OGC          := $(shell command -v openapi-generator-cli 2>/dev/null || command -v openapi-generator 2>/dev/null)
BROTLI       := $(shell command -v brotli 2>/dev/null)
GZIP         := $(shell command -v gzip 2>/dev/null)
.DEFAULT_GOAL := build

# --- version metadata (ldflags) ---
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT   := $(shell git rev-parse HEAD 2>/dev/null || echo "none")
DATE     := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILT_BY := $(shell whoami)

LDFLAGS := -X '$(MODULE_PATH)/internal/version.Version=$(VERSION)' \
           -X '$(MODULE_PATH)/internal/version.Commit=$(COMMIT)' \
           -X '$(MODULE_PATH)/internal/version.Date=$(DATE)' \
           -X '$(MODULE_PATH)/internal/version.BuiltBy=$(BUILT_BY)'

# --- phony targets ---
.PHONY: all prepare ui-install ui-build ui swagger build clean fmt vet tidy upgrade \
        sdk sdk-go sdk-ts sdk-ts-ui sdk-all worksync wire-sdk-replace help dev ui-compress \
        print-version

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
	@rm -rf $(UI_DEST_DIR)/dist
	@if [ -n "$(YARN)" ]; then \
		cd $(UI_DIR) && yarn build; \
	else \
		cd $(UI_DIR) && npm run build; \
	fi

ui-compress: ui-build
	@echo ">> Precompressing assets (brotli + gzip) in $(UI_DEST_DIR)/dist"
	@if [ -n "$(BROTLI)" ]; then \
		find "$(UI_DEST_DIR)/dist" -type f \( -name '*.js' -o -name '*.css' -o -name '*.html' \) -print0 | \
		 xargs -0 -I{} brotli -f {}; \
	else echo "brotli not found; skipping .br"; fi
	@if [ -n "$(GZIP)" ]; then \
		find "$(UI_DEST_DIR)/dist" -type f \( -name '*.js' -o -name '*.css' -o -name '*.html' \) -print0 | \
		 xargs -0 -I{} gzip -kf {}; \
	else echo "gzip not found; skipping .gz"; fi

ui: ui-compress

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
build: prepare ui swagger sdk-all
	@echo ">> Building Go binary: $(BIN)"
	@$(GOCMD) build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN) $(MAIN)

# Handy: print resolved version metadata
print-version:
	@echo "VERSION  = $(VERSION)"
	@echo "COMMIT   = $(COMMIT)"
	@echo "DATE     = $(DATE)"
	@echo "BUILT_BY = $(BUILT_BY)"
	@echo "LDFLAGS  = $(LDFLAGS)"

# --- development ---
dev: ui-install swagger
	@echo ">> Starting Vite (frontend) and Go API (backend) with dev env..."
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
		$(GOCMD) run . serve & \
		wait \
	)

# --- sdk generation (Go) ---
sdk-go: swagger
	@echo ">> Generating Go SDK (module $(GIT_HOST)/$(GIT_USER)/$(SDK_REPO), Go $(GO_VERSION))..."
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
	rm -rf "$(SDK_OUTDIR)"; \
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
	$(GOCMD) mod tidy

# --- sdk generation (TypeScript) ---
sdk-ts: swagger
	@set -e; \
	if [ -z "$(OGC)" ]; then \
		if [ -z "$(NPM)" ]; then echo "Error: npm is required to install openapi-generator-cli." >&2; exit 1; fi; \
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
		-p npmName=$(SDK_TS_NPM_NAME) \
		-p npmVersion=$(SDK_TS_NPM_VER) \
		$(SDK_TS_PROPS_FLAGS); \
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
		-p npmName=$(SDK_TS_NPM_NAME) \
		-p npmVersion=$(SDK_TS_NPM_VER) \
		$(SDK_TS_PROPS_FLAGS); \
	# --- move src/* up one level ---
	@if [ -d "$(SDK_TS_UI_DIR)/src" ]; then \
    		mv "$(SDK_TS_UI_DIR)/src/"* "$(SDK_TS_UI_DIR)/"; \
    		rm -rf "$(SDK_TS_UI_DIR)/src"; \
    	fi; \
	rm -f "$(SDK_TS_UI_DIR)/package.json" "$(SDK_TS_UI_DIR)/tsconfig.json" "$(SDK_TS_UI_DIR)/README.md"

# convenience
sdk-all: sdk-go sdk-ts sdk-ts-ui
sdk: sdk-go

# --- clean/help ---
clean:
	@echo ">> Cleaning artifacts..."
	@rm -rf "$(BIN)" docs/swagger.* docs/docs.go $(UI_DEST_DIR)/dist $(UI_DIR)/dist $(UI_DIR)/node_modules "$(SDK_OUTDIR)" "$(SDK_TS_OUTDIR)"

help:
	@echo "Targets:"
	@echo "  build           - fmt, vet, tidy, upgrade, build UI, generate Swagger, build Go binary (with ldflags)"
	@echo "  ui              - build the Vite UI and copy to $(UI_DEST_DIR)/dist (with compression)"
	@echo "  swagger         - (re)generate Swagger docs using swag"
	@echo "  sdk-go (sdk)    - generate Go SDK with correct module path and Go version"
	@echo "  sdk-ts          - generate TypeScript SDK (typescript-fetch) with package.json"
	@echo "  sdk-ts-ui       - generate TypeScript SDK directly into ui/src for inline consumption"
	@echo "  sdk-all         - generate both Go and TypeScript SDKs"
	@echo "  dev             - run Vite UI dev server + Go API"
	@echo "  clean           - remove binary, Swagger outputs, UI dist, and SDKs"
	@echo "  prepare         - fmt, vet, tidy, upgrade deps"
	@echo "  print-version   - show computed ldflags values"
