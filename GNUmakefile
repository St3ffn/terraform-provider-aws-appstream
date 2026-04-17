#
# Build/install settings
#
BINARY_NAME := terraform-provider-aws-appstream
BUILD_DIR := ./bin
GOBIN := $(shell go env GOBIN)
GOPATH := $(shell go env GOPATH)
INSTALL_DIR := $(if $(GOBIN),$(GOBIN),$(GOPATH)/bin)
TESTFLAGS ?=
COMMITLINTFLAGS ?=

default: fmt generate lint govulncheck test build

commitlint:
	@echo "🧾  Validating commit messages..."
	@cd tools/commitlint && go run ./cmd/commitlint $(COMMITLINTFLAGS)
	@echo "✅  Commitlint completed"

build:
	@echo "🚧  Building provider..."
	@mkdir -p $(BUILD_DIR)
	@go build -v -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅  Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

build-debug:
	@echo "🐞  Building provider (debug)..."
	@mkdir -p $(BUILD_DIR)
	@go build -gcflags="all=-N -l" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅  Debug build completed: $(BUILD_DIR)/$(BINARY_NAME)"

install:
	@echo "📦  Installing provider..."
	@go install -v .
	@echo "✅  Install completed: $(INSTALL_DIR)/$(BINARY_NAME)"

install-tool-commitlint:
	@echo "📦  Installing tools/commitlint..."
	@cd tools/commitlint && go install -v ./cmd/commitlint
	@echo "✅  tools/commitlint install completed: $(INSTALL_DIR)/commitlint"

fmt: fmt-provider fmt-tool-provider-codegen fmt-tool-commitlint
	@echo "✅  Format completed"

fmt-provider:
	@echo "🧹  Formatting provider Go files..."
	@golangci-lint fmt ./...
	@echo "✅  Provider format completed"

fmt-tool-provider-codegen:
	@echo "🧹  Formatting tools/provider-codegen Go files..."
	@cd tools/provider-codegen && golangci-lint fmt --config ../../.golangci.yaml ./...
	@echo "✅  tools/provider-codegen format completed"

fmt-tool-commitlint:
	@echo "🧹  Formatting tools/commitlint Go files..."
	@cd tools/commitlint && golangci-lint fmt --config ../../.golangci.yaml ./...
	@echo "✅  tools/commitlint format completed"

lint: lint-provider lint-tool-provider-codegen lint-tool-commitlint
	@echo "✅  Lint completed"

lint-provider:
	@echo "🔍  Linting provider..."
	@golangci-lint run
	@echo "✅  Provider lint completed"

lint-tool-provider-codegen:
	@echo "🔍  Linting tools/provider-codegen..."
	@cd tools/provider-codegen && golangci-lint run
	@echo "✅  tools/provider-codegen lint completed"

lint-tool-commitlint:
	@echo "🔍  Linting tools/commitlint..."
	@cd tools/commitlint && golangci-lint run --config ../../.golangci.yaml
	@echo "✅  tools/commitlint lint completed"

govulncheck: govulncheck-provider govulncheck-tool-provider-codegen govulncheck-tool-commitlint
	@echo "✅  govulncheck completed"

govulncheck-provider:
	@echo "🔐  Running govulncheck (provider)..."
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@echo "✅  Provider govulncheck completed"

govulncheck-tool-provider-codegen:
	@echo "🔐  Running govulncheck (tools/provider-codegen)..."
	@cd tools/provider-codegen && go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@echo "✅  tools/provider-codegen govulncheck completed"

govulncheck-tool-commitlint:
	@echo "🔐  Running govulncheck (tools/commitlint)..."
	@cd tools/commitlint && go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@echo "✅  tools/commitlint govulncheck completed"

generate:
	@echo "⚙️  Running code generation..."
	@cd tools/provider-codegen && go generate ./...
	@echo "✅  Generation completed"

test: test-provider test-tool-provider-codegen test-tool-commitlint
	@echo "✅  Tests completed"

test-provider:
	@echo "🧪  Running provider tests..."
	@go test -v -cover -timeout=10m $(TESTFLAGS) ./...
	@echo "✅  Provider tests completed"

test-tool-provider-codegen:
	@echo "🧪  Running tools/provider-codegen tests..."
	@cd tools/provider-codegen && go test -v -cover -timeout=5m $(TESTFLAGS) ./...
	@echo "✅  tools/provider-codegen tests completed"

test-tool-commitlint:
	@echo "🧪  Running tools/commitlint tests..."
	@cd tools/commitlint && go test -v -timeout=5m $(TESTFLAGS) ./...
	@echo "✅  tools/commitlint tests completed"

testacc:
	@echo "🧪  Running acceptance tests..."
	@TF_ACC=1 go test -v -cover -timeout 120m $(TESTFLAGS) ./...
	@echo "✅  Acceptance tests completed"

clean:
	@echo "🧹  Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "✅  Clean complete"

.PHONY: \
	commitlint \
	fmt \
	fmt-provider \
	fmt-tool-provider-codegen \
	fmt-tool-commitlint \
	lint \
	lint-provider \
	lint-tool-provider-codegen \
	lint-tool-commitlint \
	test \
	test-provider \
	test-tool-provider-codegen \
	test-tool-commitlint \
	govulncheck \
	govulncheck-provider \
	govulncheck-tool-provider-codegen \
	govulncheck-tool-commitlint \
	testacc \
	build \
	build-debug \
	install \
	install-tool-commitlint \
	generate \
	clean
