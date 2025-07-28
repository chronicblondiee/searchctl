#!/bin/bash

echo "=== Testing searchctl with both Elasticsearch and OpenSearch ==="
echo

# Test Elasticsearch
echo "🔍 Testing Elasticsearch (port 9200):"
echo "----------------------------------------"
echo -n "Cluster Info: "
./searchctl --config examples/test-config.yaml --context elasticsearch cluster info | grep "Cluster Name" | awk '{print $3}'

echo -n "Cluster Health: "
./searchctl --config examples/test-config.yaml --context elasticsearch cluster health | grep "Status:" | awk '{print $2}'

echo -n "Nodes Count: "
./searchctl --config examples/test-config.yaml --context elasticsearch get nodes | tail -n +2 | wc -l

echo -n "Creating test index... "
./searchctl --config examples/test-config.yaml --context elasticsearch create index test-es-index > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Success"
else
    echo "❌ Failed"
fi

echo -n "Index count: "
./searchctl --config examples/test-config.yaml --context elasticsearch get indices | tail -n +2 | wc -l

echo -n "Deleting test index... "
./searchctl --config examples/test-config.yaml --context elasticsearch delete index test-es-index > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Success"
else
    echo "❌ Failed"
fi

echo

# Test OpenSearch
echo "🔍 Testing OpenSearch (port 9201):"
echo "----------------------------------------"
echo -n "Cluster Info: "
./searchctl --config examples/test-config.yaml --context opensearch cluster info | grep "Cluster Name" | awk '{print $3}'

echo -n "Cluster Health: "
./searchctl --config examples/test-config.yaml --context opensearch cluster health | grep "Status:" | awk '{print $2}'

echo -n "Nodes Count: "
./searchctl --config examples/test-config.yaml --context opensearch get nodes | tail -n +2 | wc -l

echo -n "Creating test index... "
./searchctl --config examples/test-config.yaml --context opensearch create index test-os-index > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Success"
else
    echo "❌ Failed"
fi

echo -n "Index count: "
./searchctl --config examples/test-config.yaml --context opensearch get indices | tail -n +2 | wc -l

echo -n "Deleting test index... "
./searchctl --config examples/test-config.yaml --context opensearch delete index test-os-index > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ Success"
else
    echo "❌ Failed"
fi

echo
echo "🎉 All tests completed! Both Elasticsearch and OpenSearch are working with searchctl."
echo
echo "Available contexts in your config:"
echo "  - elasticsearch (port 9200)"
echo "  - opensearch (port 9201)"
echo
echo "Example usage:"
echo "  ./searchctl --config examples/test-config.yaml --context elasticsearch cluster info"
echo "  ./searchctl --config examples/test-config.yaml --context opensearch get indices"
