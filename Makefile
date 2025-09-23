GOCMD        ?= go
GOINSTALL    := $(GOCMD) install
BIN          ?= autoglue
MAIN         ?= main.go
UI_DIR       ?= ui
UI_DEST_DIR  ?= internal/ui

SWAG         := $(shell command -v swag 2>/dev/null)
GMU          := $(shell command -v go-mod-upgrade 2>/dev/null)
YARN         := $(shell command -v yarn 2>/dev/null)
NPM          := $(shell command -v npm 2>/dev/null)

.PHONY: all prepare ui-install ui-build ui swagger build clean fmt vet tidy upgrade help

all: build

prepare: fmt vet tidy upgrade

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

ui: ui-build

swagger:
	@echo ">> Generating Swagger docs..."
	@if [ -z "$(SWAG)" ]; then \
		echo "Installing swag..."; \
		$(GOINSTALL) github.com/swaggo/swag/cmd/swag@latest; \
	fi
	@rm -rf docs/swagger.* docs/docs.go
	@swag init -g $(MAIN) -o docs

build: prepare ui swagger
	@echo ">> Building Go binary: $(BIN)"
	@$(GOCMD) build -o $(BIN) $(MAIN)

clean:
	@echo ">> Cleaning artifacts..."
	@rm -rf $(BIN) docs/swagger.* docs/docs.go $(UI_DEST_DIR)/dist $(UI_DIR)/dist $(UI_DIR)/node_modules

dev: swagger
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

help:
	@echo "Targets:"
	@echo "  build      - fmt, vet, tidy, upgrade, build UI, generate Swagger, build Go binary"
	@echo "  ui         - build the Vite UI (auto-detect yarn/npm)"
	@echo "  swagger    - (re)generate Swagger docs using swag"
	@echo "  clean      - remove binary, Swagger outputs, and UI dist"
	@echo "  prepare    - fmt, vet, tidy, upgrade deps"
	@echo "  dev        - run Vite UI dev server + Go API with UI_DEV=1"
