# Makefile
# Usage:
#   make              # generate Markdown from schema
#   make schema       # just (re)generate schema.json
#   make clean        # remove generated files
#
# Notes:
# - CLI defaults to 'tofu' but can be switched to 'terraform' by doing: make CLI=terraform
# - Your cwd must contain a minimal config that requires the provider so the schema includes it.

CLI ?= tofu
SCHEMA ?= schema.json
OUT ?= autoglue_provider_reference.md
GEN ?= scripts/gen_md.py
PROVIDER ?= glueops/autoglue/autoglue

.PHONY: all schema clean

all: $(OUT)

schema: $(SCHEMA)

$(SCHEMA):
	$(CLI) providers schema -json > $(SCHEMA)

$(OUT): $(SCHEMA) $(GEN)
	python3 $(GEN) --provider "$(PROVIDER)" --schema "$(SCHEMA)" --out "$(OUT)"

clean:
	rm -f $(SCHEMA) $(OUT)
