# Commands Reference

## Global Flags

All commands support these global flags:

- `--config` - Specify config file location (default: `~/.searchctl/config.yaml`)
- `--context` - Override current context
- `--output, -o` - Output format: `table` (default), `json`, `yaml`, `wide`
- `--verbose, -v` - Enable verbose output
- `--dry-run` - Show what would be done without executing

## Core Commands

### get

List and display resources from the cluster.

#### get indices
```bash
searchctl get indices [INDEX_PATTERN] [flags]
```

**Aliases:** `index`, `idx`

**Examples:**
```bash
# List all indices
searchctl get indices

# List indices matching pattern
searchctl get indices logs-*

# Output as JSON
searchctl get indices -o json
```

#### get nodes
```bash
searchctl get nodes [flags]
```

**Aliases:** `node`, `no`

**Examples:**
```bash
# List all nodes
searchctl get nodes

# List nodes with wide output
searchctl get nodes -o wide
```

#### get datastreams
```bash
searchctl get datastreams [PATTERN] [flags]
```

**Aliases:** `datastream`, `ds`

**Examples:**
```bash
# List all data streams
searchctl get datastreams

# List data streams matching pattern
searchctl get datastreams logs-*

# Output as JSON
searchctl get datastreams -o json
```

### describe

Show detailed information about specific resources.

#### describe index
```bash
searchctl describe index INDEX_NAME [flags]
```

**Aliases:** `idx`

**Examples:**
```bash
# Describe specific index
searchctl describe index my-logs-index

# Output as YAML
searchctl describe index my-index -o yaml
```

### create

Create new resources in the cluster.

#### create index
```bash
searchctl create index INDEX_NAME [flags]
```

**Aliases:** `idx`

**Examples:**
```bash
# Create new index
searchctl create index new-logs

# Dry run creation
searchctl create index test-index --dry-run
```

#### create datastream
```bash
searchctl create datastream DATA_STREAM_NAME [flags]
```

**Aliases:** `ds`

**Examples:**
```bash
# Create new data stream
searchctl create datastream logs-nginx

# Dry run creation
searchctl create datastream logs-test --dry-run
```

### delete

Delete resources from the cluster.

#### delete index
```bash
searchctl delete index INDEX_NAME [flags]
```

**Aliases:** `idx`

**Examples:**
```bash
# Delete index
searchctl delete index old-logs

# Dry run deletion
searchctl delete index temp-index --dry-run
```

#### delete datastream
```bash
searchctl delete datastream DATA_STREAM_NAME [flags]
```

**Aliases:** `ds`

**Examples:**
```bash
# Delete data stream and all backing indices
searchctl delete datastream old-logs

# Dry run deletion
searchctl delete datastream temp-stream --dry-run
```

### apply

Apply configurations from files.

```bash
searchctl apply -f FILE [flags]
```

**Flags:**
- `-f, --filename` - Configuration file to apply (required)

**Examples:**
```bash
# Apply index template
searchctl apply -f index-template.yaml

# Dry run apply
searchctl apply -f config.yaml --dry-run
```

### rollover

Rollover data streams to create new backing indices.

#### rollover datastream
```bash
searchctl rollover datastream DATA_STREAM_NAME [flags]
```

**Aliases:** `ds`

**Flags:**
- `--max-age` - Maximum age before rollover (e.g., 30d, 1h)
- `--max-docs` - Maximum number of documents before rollover
- `--max-size` - Maximum index size before rollover (e.g., 50gb, 5gb)
- `--max-primary-shard-size` - Maximum primary shard size before rollover (e.g., 50gb)
- `--max-primary-shard-docs` - Maximum number of documents in primary shard before rollover
- `--lazy` - Only mark data stream for rollover at next write (data streams only)
- `-f, --conditions-file` - JSON file containing rollover conditions

**Examples:**
```bash
# Rollover based on age and document count
searchctl rollover datastream logs-nginx --max-age 7d --max-docs 1000000

# Rollover based on size
searchctl rollover datastream logs-app --max-size 50gb

# Rollover based on primary shard docs
searchctl rollover datastream logs-metrics --max-primary-shard-docs 500000

# Lazy rollover (mark for rollover at next write)
searchctl rollover datastream logs-system --lazy --max-age 1d

# Rollover using conditions file
searchctl rollover datastream logs-system -f rollover-conditions.json

# Dry run rollover
searchctl rollover datastream logs-test --dry-run --max-age 1d

# Output as JSON
searchctl rollover datastream logs-metrics --max-docs 500000 -o json
```

## Cluster Commands

### cluster health
```bash
searchctl cluster health [flags]
```

Display cluster health status including node counts, shard distribution, and overall status.

**Examples:**
```bash
# Show cluster health
searchctl cluster health

# Health as JSON
searchctl cluster health -o json
```

### cluster info
```bash
searchctl cluster info [flags]
```

Display general cluster information including version and cluster details.

**Examples:**
```bash
# Show cluster info
searchctl cluster info

# Info as YAML
searchctl cluster info -o yaml
```

## Configuration Commands

### config view
```bash
searchctl config view [flags]
```

Display the current configuration including contexts, clusters, and users.

**Examples:**
```bash
# View current config
searchctl config view

# View config as JSON
searchctl config view -o json
```

### config use-context
```bash
searchctl config use-context CONTEXT_NAME [flags]
```

Switch to a different context for subsequent operations.

**Examples:**
```bash
# Switch to production context
searchctl config use-context production

# Switch to development context
searchctl config use-context development
```

## Output Formats

### table (default)
Human-readable tabular format with aligned columns.

### json
Machine-readable JSON format suitable for scripting and automation.

### yaml
YAML format useful for configuration and human-readable structured data.

### wide
Extended table format with additional columns and details.

## Exit Codes

- `0` - Success
- `1` - General error
- `2` - Configuration error
- `3` - Connection error
- `4` - Resource not found
