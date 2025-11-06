# ---- Provider build/dev settings ----
BIN ?= terraform-provider-autoglue
VER ?= 0.0.1
OS  ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH ?= $(shell uname -m | sed 's/x86_64/amd64/;s/arm64/arm64/')
PROVIDER_SRC ?= glueops/autoglue/autoglue

# ---- tfplugindocs settings ----
# Where Go places binaries
BIN_DIR := $(shell go env GOBIN)
ifeq ($(BIN_DIR),)
BIN_DIR := $(shell go env GOPATH)/bin
endif
DOCS_BIN := $(BIN_DIR)/tfplugindocs
DOCS_DIR ?= docs

.PHONY: build tidy dev clean tools docs docs-validate docs-clean docs-readme

# Build the provider binary in the repo root
build:
	go build -o $(BIN) .

# Tidy module deps
tidy:
	go mod tidy

# Install the provider locally for Terraform/OpenTofu as a dev provider
dev:
	@echo "Installing dev provider v$(VER) for $(OS)_$(ARCH)..."
	@DST="$${HOME}/.terraform.d/plugins/$(PROVIDER_SRC)/$(VER)/$(OS)_$(ARCH)"; \
	mkdir -p "$$DST"; \
	go build -o "$$DST/terraform-provider-autoglue_v$(VER)" .; \
	echo "Provider installed to $$DST"; \
	echo "Run: terraform init -upgrade"

# Remove build artifacts (and optionally generated docs if desired)
clean:
	rm -f $(BIN)

# ---------- Docs via tfplugindocs ----------

# Ensure tfplugindocs is available; install if missing
tools:
	@set -e; \
	if [ ! -x "$(DOCS_BIN)" ]; then \
		echo "tfplugindocs not found. Installing..."; \
		go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest; \
		echo "Installed tfplugindocs to $(DOCS_BIN)"; \
	else \
		echo "tfplugindocs found at $(DOCS_BIN)"; \
	fi

# Generate docs from your Go schemas into ./docs
docs: tools
	@echo "Generating provider docs into ./$(DOCS_DIR)..."
	@$(DOCS_BIN)
	@echo "Done. See ./$(DOCS_DIR)"

# Validate docs are up-to-date (useful in CI)
docs-validate: tools
	@$(DOCS_BIN) validate

# Clean generated docs
docs-clean:
	@rm -rf $(DOCS_DIR)

# OPTIONAL: copy the generated landing page to README.md
docs-readme: docs
	@[ -f "$(DOCS_DIR)/index.md" ] && cp "$(DOCS_DIR)/index.md" README.md || \
	  (echo "$(DOCS_DIR)/index.md not found. Did doc generation run?"; exit 1)
	@echo "README.md updated from $(DOCS_DIR)/index.md"
