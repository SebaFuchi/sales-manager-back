#!/bin/bash

# Sales Manager Backend Setup Script
# This script initializes the backend, resolves dependencies, and verifies the build

set -e

echo "🚀 Sales Manager Backend Setup"
echo "================================"

# Change to backend directory
cd "$(dirname "$0")"

echo ""
echo "📦 Step 1: Resolving Go module dependencies..."
go mod tidy

echo ""
echo "🔨 Step 2: Building project..."
go build -o bin/sales-manager-api cmd/main.go

echo ""
echo "✓ Build successful! Binary created at: bin/sales-manager-api"

echo ""
echo "📋 Next steps:"
echo "  1. Set environment variables:"
echo "     export DB='user:password@tcp(host:3306)/sales_manager'"
echo "     export PORT='8080'"
echo ""
echo "  2. Run the server:"
echo "     ./bin/sales-manager-api"
echo ""
echo "  Or run directly with:"
echo "     go run cmd/main.go"
echo ""
echo "✨ Setup complete!"
