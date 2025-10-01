# Makefile for gRPC Hello World with TLS and File Signing

# Variables
PROJECT := hellogrpc
PROTO_DIR := proto
PROTO_FILE := $(PROTO_DIR)/hello.proto
SERVER_DIR := server
CLIENT_DIR := client
CERTS_DIR := certs
BIN_DIR := bin

# Go variables
GO := go
PROTOC := protoc
GOPATH := $(shell go env GOPATH)
GOBIN := $(shell go env GOBIN)

# If GOBIN is not set, use GOPATH/bin
ifeq ($(GOBIN),)
	GOBIN := $(GOPATH)/bin
endif

# Add GOBIN to PATH for protoc
export PATH := $(GOBIN):$(PATH)

# Binaries
SERVER_BIN := $(BIN_DIR)/server
CLIENT_BIN := $(BIN_DIR)/client

# Default target
.DEFAULT_GOAL := help

# Complete setup and build (run this for first time setup)
all: init install-tools certs build
	@echo ""
	@echo "✅ Complete! Project is ready to use."
	@echo "Run 'make test' to verify everything works."

# Initialize project (run this first)
init:
	@echo "🚀 Initializing project..."
	@test -f go.mod || $(GO) mod init $(PROJECT)
	@$(GO) get google.golang.org/grpc
	@$(GO) get google.golang.org/protobuf
	@$(GO) mod tidy
	@mkdir -p $(BIN_DIR) $(CERTS_DIR)
	@echo "✅ Project initialized!"

# Install protobuf compiler plugins
install-tools:
	@echo "📦 Installing protobuf tools..."
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "✅ Tools installed!"

# Generate TLS certificates
certs:
	@echo "🔐 Generating TLS certificates..."
	@chmod +x certs/generate.sh
	@cd certs && ./generate.sh
	@echo "✅ Certificates generated in certs/"

# Generate protobuf code (clean output)
proto:
	@echo "⚙️  Generating protobuf code..."
	@rm -rf $(PROTO_DIR)/hellogrpc $(PROTO_DIR)/hello-grpc $(PROTO_DIR)/github.com
	@$(PROTOC) --go_out=. --go_opt=module=$(PROJECT) \
		--go-grpc_out=. --go-grpc_opt=module=$(PROJECT) \
		$(PROTO_FILE)
	@echo "✅ Protobuf code generated!"

# Build server
server: proto
	@echo "🔨 Building server..."
	@cd $(SERVER_DIR) && $(GO) build -o ../$(SERVER_BIN) .
	@echo "✅ Server built: $(SERVER_BIN)"

# Build client
client: proto
	@echo "🔨 Building client..."
	@cd $(CLIENT_DIR) && $(GO) build -o ../$(CLIENT_BIN) .
	@echo "✅ Client built: $(CLIENT_BIN)"

# Build both
build: server client

# Run server
run-server: server
	@echo "🚀 Starting server on :50051..."
	@./$(SERVER_BIN)

# Run client (just hello)
run-client: client
	@echo "📡 Running client..."
	@./$(CLIENT_BIN)

# Sign a file
sign: client
	@if [ -z "$(FILE)" ]; then \
		echo "❌ Usage: make sign FILE=path/to/file"; \
		exit 1; \
	fi
	@./$(CLIENT_BIN) "$(FILE)"

# Create test file and sign it
test: build
	@echo "🧪 Running test..."
	@echo "This is a test file for signing" > test.txt
	@./$(SERVER_BIN) & SERVER_PID=$$!; \
	sleep 2; \
	./$(CLIENT_BIN) test.txt; \
	kill $$SERVER_PID 2>/dev/null || true; \
	wait $$SERVER_PID 2>/dev/null || true; \
	echo "✅ Test completed!"

# Clean everything
clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(PROTO_DIR)/hellogrpc $(PROTO_DIR)/hello-grpc $(PROTO_DIR)/github.com
	@rm -f $(PROTO_DIR)/*.pb.go
	@rm -f *.pb.go
	@rm -f $(CERTS_DIR)/*.crt $(CERTS_DIR)/*.key $(CERTS_DIR)/*.srl
	@rm -f *.txt *.sig *.pub
	@$(GO) clean
	@echo "✅ Clean completed!"

# Deep clean (including go.mod)
distclean: clean
	@rm -f go.mod go.sum
	@echo "✅ Deep clean completed!"

# Help
help:
	@echo "📚 Available commands:"
	@echo ""
	@echo "  Setup:"
	@echo "    make all           - Complete setup (init + tools + certs + build)"
	@echo "    make init          - Initialize project"
	@echo "    make install-tools - Install protobuf compiler plugins"
	@echo "    make certs         - Generate TLS certificates"
	@echo ""
	@echo "  Build:"
	@echo "    make proto         - Generate protobuf code"
	@echo "    make server        - Build server"
	@echo "    make client        - Build client"
	@echo "    make build         - Build both server and client"
	@echo ""
	@echo "  Run:"
	@echo "    make run-server    - Start the server"
	@echo "    make run-client    - Run client (hello only)"
	@echo "    make sign FILE=... - Sign a file"
	@echo "    make test          - Run full test"
	@echo ""
	@echo "  Clean:"
	@echo "    make clean         - Remove generated files"
	@echo "    make distclean     - Remove everything including go.mod"
	@echo ""
	@echo "🚀 Quick start:"
	@echo "  make all           - Complete first-time setup"
	@echo "  make test          - Run full test"
	@echo "  make run-server    - Start server"
	@echo "  make sign FILE=... - Sign a file"

.PHONY: all init install-tools certs proto server client build \
        run-server run-client sign test clean distclean help