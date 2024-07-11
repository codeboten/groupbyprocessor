# SRC_ROOT is the top of the source tree.
SRC_ROOT := $(shell git rev-parse --show-toplevel)
TOOLS_BIN_DIR    := $(SRC_ROOT)/../opentelemetry-collector-contrib/.tools
GO               := go
GOIMPORTS        := $(TOOLS_BIN_DIR)/goimports

.PHONY: generate fmt

fmt: $(GOIMPORTS)
	gofmt  -w -s ./
	$(GOIMPORTS) -w  -local github.com/codeboten/groupbyprocessor ./

generate:
	$(GO) generate ./...
	$(MAKE) fmt

tidy:
	$(GO) mod tidy -compat=1.21