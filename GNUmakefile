#
# Build/install settings
#
BINARY_NAME := terraform-provider-aws-appstream
BUILD_DIR := ./bin
GOBIN := $(shell go env GOBIN)
GOPATH := $(shell go env GOPATH)
INSTALL_DIR := $(if $(GOBIN),$(GOBIN),$(GOPATH)/bin)

default: generate fmt lint test install

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

generate:
	@echo "⚙️  Running code generation..."
	@cd tools/provider-codegen && go generate ./...
	@echo "✅  Generation completed"

fmt:
	@echo "🧹  Formatting Go files..."
	@gofmt -s -w -e .
	@echo "✅  Format completed"

test:
	@echo "🧪  Running tests..."
	@go test -v -cover -timeout=5m ./...
	@echo "✅  Tests completed"

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
	lint \
	lint-provider \
	lint-tool-provider-codegen \
	lint-tool-appstream-changelog-issues \
	test \
	testacc \
	build \
	build-debug \
	install \
	generate \
	clean
