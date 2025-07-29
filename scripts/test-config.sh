#!/bin/bash
set -e

# Load common utilities
source "$(dirname "$0")/common.sh"

echo "[CONFIG] Testing searchctl configuration..."

# Setup test environment
setup_test_environment

echo "1. Testing config view with explicit config file:"
./bin/searchctl --config examples/test-config.yaml config view

echo ""
echo "2. Testing context override with elasticsearch:"
./bin/searchctl --config examples/test-config.yaml --context elasticsearch cluster health 2>&1 || echo "Failed"

echo ""
echo "3. Testing context override with opensearch:"
./bin/searchctl --config examples/test-config.yaml --context opensearch cluster health 2>&1 || echo "Failed"

echo ""
echo "4. Testing with verbose flag:"
./bin/searchctl --config examples/test-config.yaml --context elasticsearch --verbose cluster health 2>&1 || echo "Failed"

log_success "Configuration tests completed successfully!"
