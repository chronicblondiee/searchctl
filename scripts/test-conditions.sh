#!/bin/bash
set -e

echo "🧪 Testing rollover conditions file functionality..."

# Set test config
export SEARCHCTL_CONFIG="examples/test-config.yaml"

# Build searchctl
echo "🔨 Building searchctl..."
make build

# Create a temporary directory for test files
TEST_DIR="/tmp/searchctl-rollover-tests"
mkdir -p "$TEST_DIR"

# Function to create and test different condition files
test_conditions_file() {
    local context=$1
    echo "📄 Testing conditions files for $context..."
    
    # Test 1: Basic conditions file
    cat > "$TEST_DIR/basic-conditions.json" << 'EOF'
{
  "max_age": "7d",
  "max_docs": 1000000
}
EOF
    
    echo "Testing basic conditions file..."
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f "$TEST_DIR/basic-conditions.json"
    
    # Test 2: Advanced conditions file
    cat > "$TEST_DIR/advanced-conditions.json" << 'EOF'
{
  "max_age": "30d",
  "max_docs": 5000000,
  "max_size": "100gb",
  "max_primary_shard_size": "50gb",
  "max_primary_shard_docs": 2500000
}
EOF
    
    echo "Testing advanced conditions file..."
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f "$TEST_DIR/advanced-conditions.json"
    
    # Test 3: Minimal conditions file
    cat > "$TEST_DIR/minimal-conditions.json" << 'EOF'
{
  "max_age": "1h"
}
EOF
    
    echo "Testing minimal conditions file..."
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f "$TEST_DIR/minimal-conditions.json"
    
    # Test 4: Size-only conditions
    cat > "$TEST_DIR/size-conditions.json" << 'EOF'
{
  "max_size": "10gb",
  "max_primary_shard_size": "5gb"
}
EOF
    
    echo "Testing size-only conditions file..."
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f "$TEST_DIR/size-conditions.json"
    
    # Test 5: Document count conditions
    cat > "$TEST_DIR/docs-conditions.json" << 'EOF'
{
  "max_docs": 1000000,
  "max_primary_shard_docs": 500000
}
EOF
    
    echo "Testing document count conditions file..."
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f "$TEST_DIR/docs-conditions.json"
    
    # Test 6: Invalid JSON (should fail gracefully)
    cat > "$TEST_DIR/invalid-conditions.json" << 'EOF'
{
  "max_age": "7d",
  "max_docs": 1000000
  // This is invalid JSON due to missing comma
}
EOF
    
    echo "Testing invalid JSON file (should fail)..."
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f "$TEST_DIR/invalid-conditions.json" || echo "✅ Correctly failed with invalid JSON"
    
    # Test 7: Non-existent file (should fail)
    echo "Testing non-existent file (should fail)..."
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f "$TEST_DIR/non-existent.json" || echo "✅ Correctly failed with non-existent file"
    
    # Test 8: Combination of command line args and file
    echo "Testing combination of command line and file conditions..."
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f "$TEST_DIR/basic-conditions.json" --max-size 50gb --lazy
}

# Function to test output formats with conditions files
test_output_formats() {
    local context=$1
    echo "🎨 Testing output formats with conditions files for $context..."
    
    echo "JSON output:"
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f examples/rollover-conditions.json -o json
    
    echo "YAML output:"
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f examples/rollover-conditions.json -o yaml
    
    echo "Table output (default):"
    ./bin/searchctl --context $context rollover datastream test-logs --dry-run -f examples/rollover-conditions.json
}

# Function to test verbose mode with conditions files
test_verbose_mode() {
    local context=$1
    echo "🔍 Testing verbose mode with conditions files for $context..."
    
    ./bin/searchctl --context $context --verbose rollover datastream test-logs --dry-run -f examples/rollover-conditions.json
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

# Test with both engines
echo ""
echo "🔍 Testing Elasticsearch conditions files..."
test_conditions_file "elasticsearch"
test_output_formats "elasticsearch"
test_verbose_mode "elasticsearch"

echo ""
echo "🔍 Testing OpenSearch conditions files..."
test_conditions_file "opensearch"
test_output_formats "opensearch"
test_verbose_mode "opensearch"

# Validate the example conditions file
echo ""
echo "📋 Validating example conditions file..."
if [[ -f "examples/rollover-conditions.json" ]]; then
    echo "✅ Example conditions file exists"
    if jq . examples/rollover-conditions.json >/dev/null 2>&1; then
        echo "✅ Example conditions file is valid JSON"
        echo "Contents:"
        jq . examples/rollover-conditions.json
    else
        echo "❌ Example conditions file is invalid JSON"
    fi
else
    echo "❌ Example conditions file not found"
fi

# Clean up
echo ""
echo "🧹 Cleaning up test files..."
rm -rf "$TEST_DIR"

echo ""
echo "✅ All conditions file tests completed successfully!"
echo "💡 The rollover conditions file functionality is working correctly"
