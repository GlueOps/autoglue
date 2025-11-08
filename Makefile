# ========= Makefile (optimized + DRY + annotated) =========

# --- strict shell & make behavior ---
SHELL := $(shell command -v bash 2>/dev/null || echo /bin/sh)
.SHELLFLAGS := -eu -o pipefail -c

.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
.ONESHELL:

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
SDK_REPO     ?= $(BIN)-sdk-go               # repo name used for module path
SDK_OUTDIR   ?= ../autoglue-sdk-go          # output directory (inside repo)
SDK_PKG      ?= ${BIN}                      # package name inside the SDK

UI_SSG_ROUTES ?= /,/login,/docs,/pricing

# Go versioning (go.mod uses major.minor; you’re on 1.25.4)
GO_VERSION   ?= 1.25.4

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
BROTLI       := $(shell command -v brotli 2>/dev/null)
GZIP         := $(shell command -v gzip 2>/dev/null)

# OpenAPI Generator wrapper (npm) and *core* version pin
OGC_WRAPPER                ?= @openapitools/openapi-generator-cli@latest
OPENAPI_GENERATOR_VERSION  ?= 7.17.0
OGC_BIN                    := npx -y $(OGC_WRAPPER)

# Cache the core generator jar (faster CI/local)
export OPENAPI_GENERATOR_CLI_CACHE_DIR ?= $(HOME)/.openapi-generator

# Toggle alias-as-model (can trigger name collisions in some specs)
ALIAS_AS_MODEL ?= true
ALIAS_FLAG     := $(if $(filter true,$(ALIAS_AS_MODEL)),--generate-alias-as-model,)

# Post-process generated Go with gofmt (OpenAPI Generator honors this env var)
export GO_POST_PROCESS_FILE := gofmt -w

# Default goal
.DEFAULT_GOAL := help

# --- version metadata (ldflags) ---
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT   := $(shell git rev-parse HEAD 2>/dev/null || echo "none")
DATE     := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILT_BY := $(shell whoami)

LDFLAGS := -X '$(MODULE_PATH)/internal/version.Version=$(VERSION)' \
           -X '$(MODULE_PATH)/internal/version.Commit=$(COMMIT)' \
           -X '$(MODULE_PATH)/internal/version.Date=$(DATE)' \
           -X '$(MODULE_PATH)/internal/version.BuiltBy=$(BUILT_BY)'

# --- whitespace trimming helper ---
trim = $(strip $1)

# sanitized copies (use these everywhere in recipes)
SDK_OUTDIR_CLEAN    := $(call trim,$(SDK_OUTDIR))
SDK_TS_DIR_CLEAN    := $(call trim,$(SDK_TS_DIR))
SDK_TS_UI_DIR_CLEAN := $(call trim,$(SDK_TS_UI_DIR))
GIT_HOST_CLEAN      := $(call trim,$(GIT_HOST))
GIT_USER_CLEAN      := $(call trim,$(GIT_USER))
SDK_REPO_CLEAN      := $(call trim,$(SDK_REPO))
SDK_PKG_CLEAN       := $(call trim,$(SDK_PKG))

# --- phony targets ---
.PHONY: all prepare ui-install ui-build ui swagger build clean fmt vet tidy upgrade \
        sdk sdk-go sdk-ts sdk-ts-ui sdk-all help dev ui-compress print-version \
        validate-spec check-tags doctor diff-swagger

# --- inputs/outputs for swagger (incremental) ---
DOCS_JSON := docs/swagger.json
DOCS_YAML := docs/swagger.yaml
# Prefer git for speed; fall back to find. Exclude UI dir.
#GO_SRCS := $(shell (git ls-files '*.go' ':!$(UI_DIR)/**' 2>/dev/null || find . -name '*.go' -not -path './$(UI_DIR)/*' -type f))
GO_SRCS := $(shell ( \
  git ls-files '*.go' ':!$(UI_DIR)/**' ':!docs/**' ':!sdk/**' ':!terraform-provider-autoglue/**' 2>/dev/null \
  || find . -name '*.go' -not -path './$(UI_DIR)/*' -not -path './docs/*' -type f \
))

# Rebuild swagger when Go sources change
$(DOCS_JSON) $(DOCS_YAML): $(GO_SRCS)
	@echo ">> Generating Swagger docs..."
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "Installing swag/v2 CLI @v2.0.0-rc4..."; \
		$(GOINSTALL) github.com/swaggo/swag/v2/cmd/swag@v2.0.0-rc4; \
	fi
	@rm -rf docs/swagger.* docs/docs.go
	@swag init -g $(MAIN) -o docs

# --- spec validation + tag guard ---
validate-spec: $(DOCS_JSON) ## Validate docs/swagger.json and pin the core OpenAPI Generator version
	@$(OGC_BIN) version-manager set "$(OPENAPI_GENERATOR_VERSION)"
	@$(OGC_BIN) version
	@$(OGC_BIN) validate -i $(DOCS_JSON)
	@echo ">> Spec valid."

check-tags: $(DOCS_JSON) ## Check that Swagger tags contain no spaces/slashes (jq optional)
	@echo ">> Checking tags for invalid characters (spaces or slashes)"
	@if command -v jq >/dev/null 2>&1; then \
		! jq -r '..|.tags? // empty | .[]' $(DOCS_JSON) | grep -Eq '[ /]'; \
	else \
		echo "jq not found; skipping tag check (install jq to enable)"; \
	fi

# Optional: quick diff between JSON & YAML swagger (nice for drift)
diff-swagger: $(DOCS_JSON) $(DOCS_YAML) ## Show diff between swagger.json and swagger.yaml (requires yq)
	@command -v yq >/dev/null 2>&1 || { echo "yq not found; brew install yq"; exit 1; }
	@diff -u <(jq -S . $(DOCS_JSON)) <(yq -o=json -S '.' $(DOCS_YAML)) || true

# --- meta targets ---
all: build ## Default meta-target: build everything
prepare: fmt vet tidy upgrade ## go fmt, vet, tidy, and upgrade module dependencies

# --- go hygiene ---
fmt: ## go fmt ./...
	@$(GOCMD) fmt ./...

vet: ## go vet ./...
	@$(GOCMD) vet ./...

tidy: ## go mod tidy
	@$(GOCMD) mod tidy

upgrade: ## Upgrade module requirements with go-mod-upgrade (best-effort)
	@echo ">> Checking go-mod-upgrade..."
	@if [ -z "$(GMU)" ]; then \
		echo "Installing go-mod-upgrade..."; \
		$(GOINSTALL) github.com/oligot/go-mod-upgrade@latest; \
	fi
	@go-mod-upgrade -f || true

# --- ui ---
ui-install: ## Install frontend dependencies (yarn or npm)
	@echo ">> Installing UI deps in $(UI_DIR)..."
	@if [ -n "$(YARN)" ]; then \
		cd $(UI_DIR) && yarn install --frozen-lockfile; \
	elif [ -n "$(NPM)" ]; then \
		cd $(UI_DIR) && npm ci; \
	else \
		echo "Error: neither yarn nor npm is installed." >&2; exit 1; \
	fi

ui-build: ui-install ## Build frontend (Vite)
	@echo ">> Building UI in $(UI_DIR)..."
	@rm -rf $(UI_DEST_DIR)/dist
	@if [ -n "$(YARN)" ]; then \
		cd $(UI_DIR) && yarn build; \
	else \
		cd $(UI_DIR) && npm run build; \
	fi

ui-compress: ui-build ## Precompress UI assets with brotli/gzip if available
	@echo ">> Precompressing assets (brotli + gzip) in $(UI_DEST_DIR)/dist"
	@if [ -n "$(BROTLI)" ]; then \
		find "$(UI_DEST_DIR)/dist" -type f \( -name '*.js' -o -name '*.css' -o -name '*.html' \) -print0 | \
		 xargs -0 -I{} brotli -f {}; \
	else echo "brotli not found; skipping .br"; fi
	@if [ -n "$(GZIP)" ]; then \
		find "$(UI_DEST_DIR)/dist" -type f \( -name '*.js' -o -name '*.css' -o -name '*.html' \) -print0 | \
		 xargs -0 -I{} gzip -kf {}; \
	else echo "gzip not found; skipping .gz"; fi

ui: ui-compress ## Build and precompress UI

# --- swagger convenience phony (kept for UX) ---
.PHONY: swagger
swagger: $(DOCS_JSON) ## Generate Swagger docs if stale
	@true

# --- build ---
build: prepare ui swagger sdk-all ## Build everything: Go hygiene, UI, Swagger, SDKs, then Go binary
	@echo ">> Building Go binary: $(BIN)"
	@$(GOCMD) build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN) $(MAIN)

# Handy: print resolved version metadata
print-version: ## Print ldflags/version metadata
	@echo "VERSION  = $(VERSION)"
	@echo "COMMIT   = $(COMMIT)"
	@echo "DATE     = $(DATE)"
	@echo "BUILT_BY = $(BUILT_BY)"
	@echo "LDFLAGS  = $(LDFLAGS)"

# --- development ---
dev: ui-install swagger ## Run Vite dev server and Go API (serve)
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

# --- shared generator flags ---
GEN_COMMON := --enable-post-process-file $(ALIAS_FLAG) -i $(DOCS_JSON)

# --- DRY macro for openapi-generator generate ---
# Usage: $(call OGC_GENERATE,<generator>,<outdir>,<extra flags>)
define OGC_GENERATE
	$(OGC_BIN) generate \
		$(GEN_COMMON) \
		-g $(1) \
		-o $(2) \
		$(3)
endef

# Convenience bundles
OAG_GIT_PROPS := --git-host "$(GIT_HOST_CLEAN)" --git-user-id "$(GIT_USER_CLEAN)" --git-repo-id "$(SDK_REPO_CLEAN)"
TS_PROPS      := -p npmName=$(SDK_TS_NPM_NAME) -p npmVersion=$(SDK_TS_NPM_VER) $(SDK_TS_PROPS_FLAGS)

# --- sdk generation (Go) ---
sdk-go: $(DOCS_JSON) validate-spec check-tags ## Generate Go SDK + tidy module
	@echo ">> Generating Go SDK (module $(GIT_HOST_CLEAN)/$(GIT_USER_CLEAN)/$(SDK_REPO_CLEAN), Go $(GO_VERSION))..."
	@$(call OGC_GENERATE,go,$(SDK_OUTDIR_CLEAN),--additional-properties=packageName=$(SDK_PKG_CLEAN) $(OAG_GIT_PROPS))
	@cd "$(SDK_OUTDIR_CLEAN)"; \
	$(GOCMD) mod edit -go=$(GO_VERSION); \
	$(GOCMD) mod tidy

# --- sdk generation (TypeScript) ---
sdk-ts: $(DOCS_JSON) validate-spec check-tags ## Generate TypeScript SDK, format, and build
	@echo ">> Generating TypeScript SDK in $(SDK_TS_DIR_CLEAN)"
	@rm -rf "$(SDK_TS_DIR_CLEAN)"; mkdir -p "$(SDK_TS_DIR_CLEAN)"
	@$(call OGC_GENERATE,$(SDK_TS_GEN),$(SDK_TS_DIR_CLEAN),$(TS_PROPS))
	@if command -v npx >/dev/null 2>&1; then \
		echo ">> Prettier: formatting generated TS SDK"; \
		cd "$(SDK_TS_DIR_CLEAN)" && npx --yes prettier -w . || true; \
	fi
	@echo ">> Installing & building TS SDK in $(SDK_TS_DIR_CLEAN)"
	@if command -v yarn >/dev/null 2>&1; then \
		cd "$(SDK_TS_DIR_CLEAN)" && yarn install --frozen-lockfile || true; \
		cd "$(SDK_TS_DIR_CLEAN)" && yarn build || true; \
	elif command -v npm >/dev/null 2>&1; then \
		cd "$(SDK_TS_DIR_CLEAN)" && npm ci || npm install || true; \
		cd "$(SDK_TS_DIR_CLEAN)" && npm run build || true; \
	else \
		echo "Warning: neither yarn nor npm is installed; skipping install/build for TS SDK."; \
	fi

# --- sdk generation (TypeScript into UI/src) ---
sdk-ts-ui: $(DOCS_JSON) validate-spec check-tags ## Generate TS SDK directly into UI src (no package files)
	@echo ">> Generating TypeScript SDK directly into UI source: $(SDK_TS_UI_DIR_CLEAN)"
	@rm -rf "$(SDK_TS_UI_DIR_CLEAN)"; mkdir -p "$(SDK_TS_UI_DIR_CLEAN)"
	@$(call OGC_GENERATE,typescript-fetch,$(SDK_TS_UI_DIR_CLEAN),$(TS_PROPS))
	@if [ -d "$(SDK_TS_UI_DIR_CLEAN)/src" ]; then \
		mv "$(SDK_TS_UI_DIR_CLEAN)/src/"* "$(SDK_TS_UI_DIR_CLEAN)/"; \
		rm -rf "$(SDK_TS_UI_DIR_CLEAN)/src"; \
	fi
	@rm -f "$(SDK_TS_UI_DIR_CLEAN)/package.json" "$(SDK_TS_UI_DIR_CLEAN)/tsconfig.json" "$(SDK_TS_UI_DIR_CLEAN)/README.md"

# --- convenience ---
sdk-all: sdk-go sdk-ts sdk-ts-ui ## Generate Go + TS SDKs (tip: run with "make -j sdk-all" for parallel)
sdk: sdk-go ## Alias for sdk-go

# --- clean/help ---
clean: ## Clean build artifacts, Swagger outputs, UI dist, and SDKs
	@echo ">> Cleaning artifacts..."
	@rm -rf "$(BIN)" docs/swagger.* docs/docs.go \
		"$(UI_DEST_DIR)/dist" "$(UI_DIR)/dist" "$(UI_DIR)/node_modules" \
		"$(SDK_OUTDIR_CLEAN)" "$(SDK_TS_DIR_CLEAN)" "$(SDK_TS_UI_DIR_CLEAN)"

doctor: ## Print environment diagnostics (shell, versions, generator availability)
	@echo ">> Make is using shell: $(SHELL)"
	@{ \
		echo ">> Detected runtime shell process: $$0"; \
		command -v bash >/dev/null 2>&1 || { echo "bash not found. On macOS: brew install bash"; exit 1; }; \
		echo ">> Versions:"; \
		bash --version | head -1 || true; \
		go version || true; \
		node -v || true; \
		npm -v || true; \
		npx -v || true; \
		jq --version || true; \
		echo ">> OpenAPI Generator (wrapper) available via npx:"; \
		$(OGC_BIN) version || true; \
	}

help: ## Show this help
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# ========= end Makefile =========
