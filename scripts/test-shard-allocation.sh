#!/bin/bash
set -e

# Integration tests for shard allocation features

source "$(dirname "$0")/common.sh"

log_test "Testing shard allocation features..."

setup_test_environment

check_environment

# Create and cleanup helpers
create_test_index() {
  local context=$1
  local name=$2
  local port
  if [[ "$context" == "elasticsearch" ]]; then
    port=9200
  else
    port=9201
  fi

  curl -s -X PUT "localhost:$port/$name" \
    -H 'Content-Type: application/json' \
    -d '{"settings":{"number_of_shards":1,"number_of_replicas":0}}' >/dev/null 2>&1 || true
}

delete_test_index() {
  local context=$1
  local name=$2
  local port
  if [[ "$context" == "elasticsearch" ]]; then
    port=9200
  else
    port=9201
  fi
  curl -s -X DELETE "localhost:$port/$name" >/dev/null 2>&1 || true
}

test_context() {
  local context=$1
  log_test "Testing $context shard allocation..."

  local idx="alloc-test-$context"
  delete_test_index "$context" "$idx"
  create_test_index "$context" "$idx"

  # Give cluster a moment
  sleep 1

  log_info "Listing shards for $idx"
  test_command "./bin/searchctl --context $context get shards $idx" true
  test_command "./bin/searchctl --context $context get shards $idx -o json" true
  test_command "./bin/searchctl --context $context get shards $idx -o yaml" true

  log_info "Explaining allocation for primary shard 0"
  # Ensure shard param is passed; index has 1 primary shard, so shard 0 exists
  test_command "./bin/searchctl --context $context describe allocation --index $idx --shard 0 --primary --include-yes -o json" true
  test_command "./bin/searchctl --context $context describe allocation --index $idx --shard 0 --primary --include-disk -o yaml" true

  log_info "Getting cluster allocation settings"
  test_command "./bin/searchctl --context $context cluster allocation-settings -o json" true

  log_info "Updating cluster allocation settings (enable=all, rebalance=all)"
  test_command "./bin/searchctl --context $context cluster allocation-settings --enable all" true
  test_command "./bin/searchctl --context $context cluster allocation-settings --rebalance all" true

  delete_test_index "$context" "$idx"

  log_success "$context shard allocation tests completed"
}

for ctx in elasticsearch opensearch; do
  test_context "$ctx"
done

log_success "Shard allocation feature tests completed successfully!"


