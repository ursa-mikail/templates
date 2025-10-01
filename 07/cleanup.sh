#!/bin/bash

# Cleanup script for gRPC project

echo "🧹 Cleaning up project structure..."

# Remove duplicate proto generated files
echo "Removing duplicate proto files..."
rm -rf proto/hellogrpc
rm -rf proto/hello-grpc
rm -rf proto/github.com

# Remove old generated pb.go files in proto dir
rm -f proto/*.pb.go

# Clean binaries
echo "Cleaning binaries..."
rm -rf bin/*

# Remove test artifacts
echo "Cleaning test files..."
rm -f *.txt *.sig *.pub

echo ""
echo "✅ Cleanup complete!"
echo ""
echo "📁 Your project structure should now look like:"
echo "./"
echo "├── bin/              (empty, will contain binaries)"
echo "├── certs/            (TLS certificates)"
echo "├── client/"
echo "│   └── main.go"
echo "├── server/"
echo "│   └── main.go"
echo "├── proto/"
echo "│   └── hello.proto"
echo "├── hellogrpc/        (generated proto code will go here)"
echo "├── Makefile"
echo "├── go.mod"
echo "└── go.sum"
echo ""
echo "🚀 Next steps:"
echo "  1. Run: make proto"
echo "  2. Run: make build"
echo "  3. Run: make test"