BIN ?= terraform-provider-autoglue
VER ?= 0.0.1
OS  ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH ?= $(shell uname -m | sed 's/x86_64/amd64/;s/arm64/arm64/')

.PHONY: build tidy dev clean

build:
	go build -o $(BIN) .

tidy:
	go mod tidy

dev:
	@echo "Installing dev provider v$(VER) for $(OS)_$(ARCH)..."
	@DST="$${HOME}/.terraform.d/plugins/glueops/autoglue/autoglue/$(VER)/$(OS)_$(ARCH)"; \
	mkdir -p "$$DST"; \
	go build -o "$$DST/terraform-provider-autoglue_v$(VER)" .; \
	echo "Provider installed to $$DST"; \
	echo "Run: terraform init -upgrade"

clean:
	rm -f $(BIN)
