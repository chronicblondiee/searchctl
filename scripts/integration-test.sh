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

echo "✅ Integration tests completed successfully!"
