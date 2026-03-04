#
# Build/install settings
#
BINARY_NAME := terraform-provider-aws-appstream
BUILD_DIR := ./bin
GOBIN := $(shell go env GOBIN)
GOPATH := $(shell go env GOPATH)
INSTALL_DIR := $(if $(GOBIN),$(GOBIN),$(GOPATH)/bin)

default: fmt generate lint govulncheck test install

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

fmt: fmt-provider fmt-tool-provider-codegen fmt-tool-appstream-changelog-issues
	@echo "✅  Format completed"

fmt-provider:
	@echo "🧹  Formatting provider Go files..."
	@golangci-lint fmt ./...
	@echo "✅  Provider format completed"

fmt-tool-provider-codegen:
	@echo "🧹  Formatting tools/provider-codegen Go files..."
	@cd tools/provider-codegen && golangci-lint fmt --config ../../.golangci.yaml ./...
	@echo "✅  tools/provider-codegen format completed"

fmt-tool-appstream-changelog-issues:
	@echo "🧹  Formatting tools/appstream-changelog-issues Go files..."
	@cd tools/appstream-changelog-issues && golangci-lint fmt --config ../../.golangci.yaml ./...
	@echo "✅  tools/appstream-changelog-issues format completed"

lint: lint-provider lint-tool-provider-codegen lint-tool-appstream-changelog-issues
	@echo "✅  Lint completed"

lint-provider:
	@echo "🔍  Linting provider..."
	@golangci-lint run
	@echo "✅  Provider lint completed"

lint-tool-provider-codegen:
	@echo "🔍  Linting tools/provider-codegen..."
	@cd tools/provider-codegen && golangci-lint run ./...
	@echo "✅  tools/provider-codegen lint completed"

lint-tool-appstream-changelog-issues:
	@echo "🔍  Linting tools/appstream-changelog-issues..."
	@cd tools/appstream-changelog-issues && golangci-lint run ./...
	@echo "✅  tools/appstream-changelog-issues lint completed"

govulncheck: govulncheck-provider govulncheck-tool-provider-codegen govulncheck-tool-appstream-changelog-issues
	@echo "✅  govulncheck completed"

govulncheck-provider:
	@echo "🔐  Running govulncheck (provider)..."
	@go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@echo "✅  Provider govulncheck completed"

govulncheck-tool-provider-codegen:
	@echo "🔐  Running govulncheck (tools/provider-codegen)..."
	@cd tools/provider-codegen && go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@echo "✅  tools/provider-codegen govulncheck completed"

govulncheck-tool-appstream-changelog-issues:
	@echo "🔐  Running govulncheck (tools/appstream-changelog-issues)..."
	@cd tools/appstream-changelog-issues && go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	@echo "✅  tools/appstream-changelog-issues govulncheck completed"

generate:
	@echo "⚙️  Running code generation..."
	@cd tools/provider-codegen && go generate ./...
	@echo "✅  Generation completed"

test: test-provider test-tool-provider-codegen test-tool-appstream-changelog-issues
	@echo "✅  Tests completed"

test-provider:
	@echo "🧪  Running provider tests..."
	@go test -v -cover -timeout=10m ./...
	@echo "✅  Provider tests completed"

test-tool-provider-codegen:
	@echo "🧪  Running tools/provider-codegen tests..."
	@cd tools/provider-codegen && go test -v -cover -timeout=5m ./...
	@echo "✅  tools/provider-codegen tests completed"

test-tool-appstream-changelog-issues:
	@echo "🧪  Running tools/appstream-changelog-issues tests..."
	@cd tools/appstream-changelog-issues && go test -v -cover -timeout=5m ./...
	@echo "✅  tools/appstream-changelog-issues tests completed"

testacc:
	@echo "🧪  Running acceptance tests..."
	@TF_ACC=1 go test -v -cover -timeout 120m ./...
	@echo "✅  Acceptance tests completed"

clean:
	@echo "🧹  Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "✅  Clean complete"

.PHONY: \
	fmt \
	fmt-provider \
	fmt-tool-provider-codegen \
	fmt-tool-appstream-changelog-issues \
	lint \
	lint-provider \
	lint-tool-provider-codegen \
	lint-tool-appstream-changelog-issues \
	test \
	test-provider \
	test-tool-provider-codegen \
	test-tool-appstream-changelog-issues \
	govulncheck \
	govulncheck-provider \
	govulncheck-tool-provider-codegen \
	govulncheck-tool-appstream-changelog-issues \
	testacc \
	build \
	build-debug \
	install \
	generate \
	clean
