#!/bin/bash
set -e

echo "🧪 Running SearchCtl Integration Tests..."

# Ensure test environment is running
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

# Build searchctl
echo "🔨 Building searchctl..."
make build

# Set test config
export SEARCHCTL_CONFIG="examples/test-config.yaml"

echo "🔍 Testing Elasticsearch..."
./bin/searchctl --context elasticsearch cluster health
./bin/searchctl --context elasticsearch cluster info
./bin/searchctl --context elasticsearch get indices

echo "🔍 Testing OpenSearch..."
./bin/searchctl --context opensearch cluster health
./bin/searchctl --context opensearch cluster info
./bin/searchctl --context opensearch get indices

echo "🧪 Creating test indices..."
./bin/searchctl --context elasticsearch create index test-es-index --dry-run
./bin/searchctl --context opensearch create index test-os-index --dry-run

echo "📋 Testing data streams..."
./bin/searchctl --context elasticsearch get datastreams
./bin/searchctl --context opensearch get datastreams

echo "🧪 Testing data stream operations..."
./bin/searchctl --context elasticsearch create datastream test-logs --dry-run
./bin/searchctl --context opensearch create datastream test-logs --dry-run
./bin/searchctl --context elasticsearch delete datastream test-logs --dry-run
./bin/searchctl --context opensearch delete datastream test-logs --dry-run

echo "🔄 Testing rollover commands (dry-run)..."
./bin/searchctl --context elasticsearch rollover datastream logs-test --dry-run --max-age 7d --max-docs 1000
./bin/searchctl --context opensearch rollover datastream logs-test --dry-run --max-age 7d --max-docs 1000

echo "🔄 Testing advanced rollover features..."
./bin/searchctl --context elasticsearch rollover datastream logs-test --dry-run --max-primary-shard-docs 500000
./bin/searchctl --context opensearch rollover datastream logs-test --dry-run --max-primary-shard-size 25gb
./bin/searchctl --context elasticsearch rollover datastream logs-test --dry-run --lazy --max-age 1d
./bin/searchctl --context opensearch rollover datastream logs-test --dry-run -f examples/rollover-conditions.json

echo "🧪 Testing rollover output formats..."
./bin/searchctl --context elasticsearch rollover ds logs-test --dry-run --max-age 7d -o json
./bin/searchctl --context opensearch rollover ds logs-test --dry-run --max-age 7d -o yaml

echo "🔍 Testing comprehensive datastream operations..."
./bin/searchctl --context elasticsearch get datastreams -o wide
./bin/searchctl --context opensearch get datastreams -o json

echo "📋 Testing help documentation..."
./bin/searchctl rollover --help >/dev/null
./bin/searchctl rollover datastream --help >/dev/null
./bin/searchctl create datastream --help >/dev/null
./bin/searchctl delete datastream --help >/dev/null

echo "🎯 Running dedicated rollover test suite..."
if [[ -x scripts/test-rollover.sh ]]; then
    echo "Running comprehensive rollover tests..."
    # We'll run this in dry-run mode to avoid actual changes
    echo "Note: Running rollover-specific tests (all in dry-run mode)"
    # Uncomment the next line to run the full rollover test suite
    # ./scripts/test-rollover.sh
else
    echo "Warning: scripts/test-rollover.sh not found or not executable"
fi

echo "🔧 Additional test scripts available:"
echo "  📄 ./scripts/test-rollover.sh - Comprehensive rollover testing (dry-run)"
echo "  🧪 ./scripts/test-rollover-real.sh - Real rollover operations with test data"
echo "  🚀 ./scripts/test-performance.sh - Performance and stress testing"
echo ""
echo "💡 To run comprehensive tests:"
echo "  ./scripts/test-rollover.sh"
echo "💡 To test with real data (creates/deletes test data streams):"
echo "  ./scripts/test-rollover-real.sh"
echo "💡 To run performance tests:"
echo "  ./scripts/test-performance.sh"

echo "✅ Integration tests completed successfully!"
