# Component Templates Examples

This directory contains example component templates that can be applied using `searchctl apply -f`.

## Files

- **base-settings.yaml**: Basic index settings component template
- **observability-mappings.yaml**: Elasticsearch-compatible mappings for observability data
- **observability-mappings-opensearch.yaml**: OpenSearch-compatible version of observability mappings

## Engine Compatibility

### Elasticsearch vs OpenSearch Differences

Some field types are engine-specific:

- **flattened**: Elasticsearch-specific field type for dynamic object fields
- **object with enabled: false**: OpenSearch alternative to flattened fields

### Usage

For **Elasticsearch**:
```bash
searchctl apply -f examples/component-templates/base-settings.yaml
searchctl apply -f examples/component-templates/observability-mappings.yaml
```

For **OpenSearch**:
```bash
searchctl apply -f examples/component-templates/base-settings.yaml
searchctl apply -f examples/component-templates/observability-mappings-opensearch.yaml
```

## Component Template Structure

Component templates use the same YAML structure as index templates but only contain reusable configuration blocks:

```yaml
kind: ComponentTemplate
metadata:
  name: my-component-template
spec:
  template:
    settings:
      number_of_shards: 1
    mappings:
      properties:
        field_name:
          type: keyword
    aliases:
      my-alias: {}
  _meta:
    description: "Description of the component template"
    version: 1
```

## Commands

- **Apply**: `searchctl apply -f <component-template-file>`
- **List**: `searchctl get component-templates`
- **Get specific**: `searchctl get component-templates <name>`
- **Delete**: `searchctl delete component-template <name>`