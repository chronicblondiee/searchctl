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

echo "✅ Integration tests completed successfully!"
