#!/bin/bash
set -e

echo "🔄 Testing SearchCtl Rollover Functionality..."

# Set test config
export SEARCHCTL_CONFIG="examples/test-config.yaml"

# Build searchctl
echo "🔨 Building searchctl..."
make build

# Function to test rollover with different conditions
test_rollover_conditions() {
    local context=$1
    echo "🧪 Testing rollover conditions for $context..."
    
    # Test individual conditions
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 7d
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-docs 1000000
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-size 50gb
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-primary-shard-size 25gb
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-primary-shard-docs 500000
    
    # Test lazy rollover
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run --lazy --max-age 1d
    
    # Test multiple conditions
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 30d --max-docs 1000000 --max-size 50gb
    
    # Test with conditions file
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f examples/rollover-conditions.json
    
    # Test different output formats
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 7d -o json
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age 7d -o yaml
    
    # Test using alias
    ./bin/searchctl --context $context rollover ds test-logs --dry-run --max-age 7d
}

# Function to test datastream operations
test_datastream_operations() {
    local context=$1
    echo "🗄️  Testing datastream operations for $context..."
    
    # Test create datastream (dry-run)
    ./bin/searchctl --context $context create datastream test-logs-new --dry-run
    ./bin/searchctl --context $context create ds test-logs-alias --dry-run
    
    # Test delete datastream (dry-run)
    ./bin/searchctl --context $context delete datastream test-logs-old --dry-run
    ./bin/searchctl --context $context delete ds test-logs-alias --dry-run
    
    # Test get datastreams
    ./bin/searchctl --context $context get datastreams
    ./bin/searchctl --context $context get ds
    ./bin/searchctl --context $context get datastreams "test-*"
    
    # Test different output formats for get
    ./bin/searchctl --context $context get datastreams -o json
    ./bin/searchctl --context $context get datastreams -o yaml
    ./bin/searchctl --context $context get datastreams -o wide
}

# Function to test error scenarios
test_error_scenarios() {
    local context=$1
    echo "❌ Testing error scenarios for $context..."
    
    # Test rollover without datastream name (should fail)
    ./bin/searchctl --context $context rollover datastream --dry-run || echo "Expected error: missing datastream name"
    
    # Test rollover with invalid conditions
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run --max-age invalid || echo "Expected error: invalid max-age format"
    
    # Test create datastream without name (should fail)
    ./bin/searchctl --context $context create datastream --dry-run || echo "Expected error: missing datastream name"
    
    # Test delete datastream without name (should fail)
    ./bin/searchctl --context $context delete datastream --dry-run || echo "Expected error: missing datastream name"
}

# Function to test verbose mode
test_verbose_mode() {
    local context=$1
    echo "🔍 Testing verbose mode for $context..."
    
    ./bin/searchctl --context $context --verbose rollover datastream test-logs --dry-run --max-age 7d
    ./bin/searchctl --context $context --verbose get datastreams
    ./bin/searchctl --context $context --verbose create datastream test-logs --dry-run
}

# Check if test environment is running
echo "🏥 Checking test environment..."
if ! curl -f http://localhost:9200/_cluster/health >/dev/null 2>&1; then
    echo "❌ Elasticsearch not running. Start test environment first:"
    echo "   ./scripts/start-test-env.sh"
    exit 1
fi

if ! curl -f http://localhost:9201/_cluster/health >/dev/null 2>&1; then
    echo "❌ OpenSearch not running. Start test environment first:"
    echo "   ./scripts/start-test-env.sh"
    exit 1
fi

echo "✅ Test environment is ready"

# Test Elasticsearch
echo ""
echo "🔍 Testing Elasticsearch rollover functionality..."
test_rollover_conditions "elasticsearch"
test_datastream_operations "elasticsearch"
test_error_scenarios "elasticsearch"
test_verbose_mode "elasticsearch"

# Test OpenSearch
echo ""
echo "🔍 Testing OpenSearch rollover functionality..."
test_rollover_conditions "opensearch"
test_datastream_operations "opensearch"
test_error_scenarios "opensearch"
test_verbose_mode "opensearch"

echo ""
echo "🎯 Testing help commands..."
./bin/searchctl rollover --help
./bin/searchctl rollover datastream --help
./bin/searchctl rollover ds --help
./bin/searchctl create datastream --help
./bin/searchctl delete datastream --help
./bin/searchctl get datastreams --help

echo ""
echo "✅ All rollover tests completed successfully!"
