#!/bin/bash

# Release script for monorepo libraries
# Usage: ./scripts/release.sh <library-name> [version-type]
# Example: ./scripts/release.sh greetings patch

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}ℹ️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Validate arguments
if [ -z "$1" ]; then
    print_error "Library name is required"
    echo "Usage: ./scripts/release.sh <library-name> [version-type]"
    echo "Available libraries: greetings, math"
    echo "Version types: patch, minor, major (default: patch)"
    exit 1
fi

LIBRARY=$1
VERSION_TYPE=${2:-patch}

# Validate library exists
if [ ! -d "libs/$LIBRARY" ]; then
    print_error "Library '$LIBRARY' not found in libs/"
    exit 1
fi

# Validate release-it config exists
if [ ! -f "libs/$LIBRARY/.release-it.json" ]; then
    print_error "Release config not found: libs/$LIBRARY/.release-it.json"
    exit 1
fi

print_info "Starting release process for library: $LIBRARY"
print_info "Version bump type: $VERSION_TYPE"

# Run tests
print_info "Running tests for $LIBRARY..."
cd "libs/$LIBRARY"
if ! go test ./...; then
    print_error "Tests failed for $LIBRARY"
    exit 1
fi
cd ../..

print_info "All tests passed ✓"

# Run release-it
print_info "Running release-it..."
npm run release:$LIBRARY -- $VERSION_TYPE

print_info "Release completed successfully! 🎉"
