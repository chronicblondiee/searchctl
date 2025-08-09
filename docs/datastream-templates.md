# Data Stream Templates Examples

This directory contains example templates for working with Elasticsearch data streams using searchctl.

## Overview

Data streams are a convenient, scalable way to ingest, search, and manage time-series data in Elasticsearch. They require index templates that include the `data_stream` configuration.

## Template Files

### Basic Templates

1. **`simple-datastream-template.yaml`** - Minimal template for testing
2. **`datastream-template.yaml`** - Comprehensive logs template
3. **`metrics-datastream-template.yaml`** - Optimized for metrics data
4. **`traces-datastream-template.yaml`** - Distributed tracing template

### Component Templates

Component templates allow you to create reusable building blocks:

1. **`component-templates/base-settings.yaml`** - Common index settings
2. **`component-templates/observability-mappings.yaml`** - Shared field mappings
3. **`component-templates/composable-datastream-template.yaml`** - Template using components

## Usage Examples

### 1. Create a Simple Data Stream

```bash
# Apply the template first
searchctl apply -f examples/simple-datastream-template.yaml

# Create the data stream
searchctl create datastream test-logs

# Check the data stream
searchctl get datastreams test-logs
```

### 2. Create a Logs Data Stream with Full Configuration

```bash
# Apply the comprehensive logs template
searchctl apply -f examples/datastream-template.yaml

# Create the data stream
searchctl create datastream logs-application

# View the data stream details
searchctl get datastreams logs-* -o yaml
```

### 3. Create Metrics Data Stream

```bash
# Apply metrics template
searchctl apply -f examples/metrics-datastream-template.yaml

# Create metrics data stream
searchctl create datastream metrics-system

# Test rollover conditions
searchctl rollover datastream metrics-system --max-age 7d --max-docs 1000000
```

### 4. Using Component Templates (Composable Templates)

```bash
# Create component templates first
searchctl apply -f examples/component-templates/base-settings.yaml
searchctl apply -f examples/component-templates/observability-mappings.yaml

# Apply the composable template that uses the components
searchctl apply -f examples/component-templates/composable-datastream-template.yaml

# Create data stream
searchctl create datastream observability-logs
```

## Template Structure

### Required Fields for Data Streams

```yaml
apiVersion: v1
kind: IndexTemplate
metadata:
  name: your-template-name
spec:
  index_patterns:
  - "your-pattern-*"
  data_stream: {}  # This makes it a data stream template
  template:
    settings:
      # Your settings
    mappings:
      properties:
        "@timestamp":  # Required for data streams
          type: date
```

### Key Configuration Options

#### Settings
- **`number_of_shards`** - Number of primary shards (default: 1 for data streams)
- **`number_of_replicas`** - Number of replica shards
- **`refresh_interval`** - How often to refresh the index
- **`lifecycle.name`** - ILM policy name for data stream management

#### Common Field Types
- **`@timestamp`** - Required timestamp field (date type)
- **`keyword`** - Exact value fields (tags, IDs, enums)
- **`text`** - Full-text searchable fields
- **`flattened`** - Dynamic object fields
- **`nested`** - Complex object arrays

## Data Stream Management

### Rollover Operations

Data streams automatically create new backing indices when they rollover:

```bash
# Manual rollover with conditions
searchctl rollover datastream logs-app \
  --max-age 7d \
  --max-docs 1000000 \
  --max-size 5gb

# Rollover with conditions file
searchctl rollover datastream logs-app -f examples/rollover-conditions.json
```

### Lifecycle Management

Apply ILM policies to manage data stream lifecycle:

```bash
# View current lifecycle
searchctl describe datastream logs-app

# Data streams inherit ILM policies from their templates
```

## Best Practices

### Template Design
1. **Use meaningful index patterns** - `logs-*`, `metrics-*`, `traces-*`
2. **Set appropriate priority** - Higher numbers override lower ones
3. **Include `@timestamp`** - Required for data streams
4. **Plan shard strategy** - Usually 1 shard per data stream is sufficient

### Field Mapping
1. **Use `keyword` for filtering** - Status, service names, host names
2. **Use `text` for search** - Log messages, descriptions
3. **Use `flattened` for dynamic objects** - Labels, metadata
4. **Avoid mapping explosion** - Use flattened or disable dynamic mapping

### Performance
1. **Configure refresh interval** - 30s for metrics, 5s for logs
2. **Use compression** - `best_compression` for cold data
3. **Plan rollover strategy** - Balance between search performance and management overhead

## Testing Templates

### Dry Run
```bash
# Test template creation without applying
searchctl apply -f examples/datastream-template.yaml --dry-run
```

### Validation
```bash
# Create test data stream
searchctl create datastream test-validation

# Index sample document
curl -X POST "localhost:9200/test-validation/_doc" \
  -H "Content-Type: application/json" \
  -d '{"@timestamp": "2024-01-01T00:00:00Z", "message": "test"}'

# Verify mapping
searchctl get datastreams test-validation -o yaml

# Cleanup
searchctl delete datastream test-validation -y
```

## Troubleshooting

### Common Issues

1. **Template not applying**
   - Check template priority conflicts
   - Verify index pattern matching
   - Ensure `data_stream: {}` is present

2. **Data stream creation fails**
   - Verify template exists and matches pattern
   - Check that `@timestamp` field is mapped as date type
   - Ensure no conflicting regular index exists

3. **Mapping conflicts**
   - Use component templates for consistency
   - Check for conflicting field types across templates
   - Use explicit mappings instead of dynamic mapping

### Debug Commands

```bash
# List all templates
searchctl get templates

# Check template details
searchctl describe template your-template-name

# View data stream backing indices
searchctl get indices your-datastream-*

# Check data stream stats
searchctl get datastreams your-stream -o json
```