# Makefile for autoglue

# Go modules
GOCMD=go
GOINSTALL=$(GOCMD) install
SWAGBIN=$(shell which swag)

# Default target
.PHONY: all
all: build

.PHONY: ui
ui:
	rm -rf api/ui
	cd ../autoglue-ui && yarn build
	mkdir -p api/ui
	cp -r ../autoglue-ui/dist/* api/ui/

# Build the CLI binary
.PHONY: build
build: ui swagger
	$(GOCMD) build -o autoglue main.go

# Run the app with .env
.PHONY: run
run:
	AUTOGLUE_DB_DSN=$$(grep AUTOGLUE_DB_DSN .env | cut -d '=' -f2-) $(GOCMD) run .

# Generate Swagger docs
.PHONY: swagger
swagger:
ifndef SWAGBIN
	@echo "Installing swag..."
	$(GOINSTALL) github.com/swaggo/swag/cmd/swag@latest
endif
	@echo "Generating Swagger docs..."
	@rm -rf docs/swagger.* docs/docs.go
	@swag init

# Clean build artifacts
.PHONY: clean
clean:
	@rm -rf autoglue docs/swagger.* docs/docs.go

# Help message
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build     - Build the autoglue CLI"
	@echo "  make run       - Run the CLI (ensure .env is configured)"
	@echo "  make swagger   - Generate Swagger docs (docs/)"
	@echo "  make clean     - Clean up binaries and Swagger docs"
