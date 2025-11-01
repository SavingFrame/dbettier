#!/bin/bash

# Database testing script for dbettier
# 
# Usage:
#   ./test.sh                    # Use testcontainers (default)
#   ./test.sh --existing         # Try to use existing postgres container
#   ./test.sh --help             # Show this help

set -e

show_help() {
    echo "Usage: ./test.sh [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --existing    Use existing postgres container if available"
    echo "  --help        Show this help message"
    echo ""
    echo "Environment variables (for custom database):"
    echo "  TEST_POSTGRES_HOST      (default: localhost)"
    echo "  TEST_POSTGRES_PORT      (default: 5432)"
    echo "  TEST_POSTGRES_USER      (default: postgres)"
    echo "  TEST_POSTGRES_PASSWORD  (default: postgres)"
    echo "  TEST_POSTGRES_DB        (default: postgres)"
    echo ""
    echo "Examples:"
    echo "  ./test.sh"
    echo "  ./test.sh --existing"
    echo "  TEST_POSTGRES_HOST=myhost TEST_POSTGRES_PORT=5433 ./test.sh"
}

if [ "$1" == "--help" ]; then
    show_help
    exit 0
fi

# Parse custom flags and collect remaining arguments for go test
GO_TEST_ARGS=()
while [[ $# -gt 0 ]]; do
    case $1 in
        --existing)
            # Check if postgres container is running
                echo "✓ Using local postgres"
                export TEST_POSTGRES_HOST=localhost
                export TEST_POSTGRES_PORT=5432
                export TEST_POSTGRES_USER=postgres
                export TEST_POSTGRES_PASSWORD=password
                export TEST_POSTGRES_DB=postgres
                echo "→ Using existing PostgreSQL instance"
            shift
            ;;
        *)
            # Pass unknown flags to go test
            GO_TEST_ARGS+=("$1")
            shift
            ;;
    esac
done

# Run tests
go test -v ./internal/database/ "${GO_TEST_ARGS[@]}"
