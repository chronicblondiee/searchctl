# Lifecycle Policy Examples

This directory contains example lifecycle policy templates that demonstrate how to use `searchctl apply -f` with lifecycle policies.

## Compatibility

These examples work with both Elasticsearch (ILM) and OpenSearch (ISM) clusters. The tool automatically detects the cluster type and uses the appropriate API.

## Examples

- `basic-ilm-policy.yaml` - Basic Elasticsearch ILM policy for log rotation
- `basic-ism-policy.yaml` - Basic OpenSearch ISM policy for log rotation
- `hot-warm-cold-policy.yaml` - Multi-tier lifecycle policy with hot/warm/cold phases
- `delete-old-logs-policy.yaml` - Simple policy to delete old log indices

## Usage

```bash
# Apply a lifecycle policy
searchctl apply -f lifecycle-policies/basic-ilm-policy.yaml

# Get lifecycle policies  
searchctl get lifecycle-policies

# Delete a lifecycle policy
searchctl delete lifecycle-policy my-policy
```

## Aliases

The lifecycle policy commands support multiple aliases:
- `lifecycle-policies`, `lifecyclepolicies`
- `lifecycle-policy`, `lifecyclepolicy` 
- `ilm`, `ism`, `lp`, `lifecycle`