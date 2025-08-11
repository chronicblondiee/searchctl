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

#### describe lifecycle-policy
```bash
searchctl describe lifecycle-policy POLICY_NAME [flags]
```

**Aliases:** `lifecyclepolicy`, `ilm`, `ism`, `lp`, `lifecycle`

**Flags:**
- `--show-body` - Include full policy body when using table output

**Examples:**
```bash
# Describe lifecycle policy (table)
searchctl describe lifecycle-policy basic-log-rotation

# As YAML (full object)
searchctl describe lifecycle-policy basic-log-rotation -o yaml
```

#### describe index-template
```bash
searchctl describe index-template TEMPLATE_NAME [flags]
```

**Aliases:** `indextemplate`, `template`, `it`, `idx-template`

**Flags:**
- `--show-body` - Include full template body when using table output

**Examples:**
```bash
# Describe index template
searchctl describe index-template logs-template

# As JSON (full object)
searchctl describe index-template logs-template -o json
```

#### describe component-template
```bash
searchctl describe component-template NAME [flags]
```

**Aliases:** `componenttemplate`, `ct`

**Flags:**
- `--show-body` - Include full template body when using table output

**Examples:**
```bash
searchctl describe component-template base-settings
```

#### describe datastream
```bash
searchctl describe datastream NAME [flags]
```

**Aliases:** `datastreams`, `ds`

**Examples:**
```bash
searchctl describe datastream logs-nginx
searchctl describe datastream logs-nginx -o yaml
```

#### describe node
```bash
searchctl describe node NODE_ID [flags]
```

**Aliases:** `no`

**Examples:**
```bash
searchctl describe node node-1
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
searchctl delete index INDEX_NAME_OR_PATTERN [flags]
```

**Aliases:** `idx`

**Flags:**
- `-y, --yes` - Automatically confirm deletion without prompting

**Examples:**
```bash
# Delete specific index (with confirmation prompt)
searchctl delete index old-logs

# Delete multiple indices with wildcard (with confirmation prompt)
searchctl delete index logs-*

# Delete indices matching pattern without confirmation
searchctl delete index test-index-* -y

# Dry run deletion
searchctl delete index temp-index --dry-run

# Dry run with wildcard
searchctl delete index logs-2024-* --dry-run
```

#### delete datastream
```bash
searchctl delete datastream DATA_STREAM_NAME_OR_PATTERN [flags]
```

**Aliases:** `ds`

**Flags:**
- `-y, --yes` - Automatically confirm deletion without prompting

**Examples:**
```bash
# Delete specific data stream and all backing indices (with confirmation prompt)
searchctl delete datastream old-logs

# Delete multiple data streams with wildcard (with confirmation prompt)
searchctl delete datastream logs-*

# Delete data streams matching pattern without confirmation
searchctl delete datastream metrics-* -y

# Dry run deletion
searchctl delete datastream temp-stream --dry-run

# Dry run with wildcard
searchctl delete datastream logs-2024-* --dry-run
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

### clone

Export (clone) and import cluster configuration.

#### clone export
```bash
searchctl clone export [flags]
```

**Flags:**
- `-d, --dir` - Output directory (required)
- `--types` - Comma-separated list of resource types to export: `index-templates,component-templates,lifecycle-policies,ingest-pipelines,cluster-settings`
- `--names` - Optional names/patterns to filter (comma-separated). Empty exports all
- `--include-system` - Include system resources (names starting with `.`)

**Examples:**
```bash
# Export everything as YAML
searchctl clone export --dir /backup -o yaml

# Export only templates and ILM/ISM as JSON
searchctl clone export --types index-templates,component-templates,lifecycle-policies --dir /backup -o json

# Export only matching names/patterns
searchctl clone export --types index-templates --names logs-*,metrics-* --dir /backup -o yaml
```

#### clone import
```bash
searchctl clone import [flags]
```

Import resources by scanning subdirectories in the given directory. Import order is component-templates → index-templates → lifecycle-policies → ingest-pipelines → cluster-settings.

**Flags:**
- `-d, --dir` - Input directory (required)
- `--types` - Comma-separated list of resource types to import
- `--dry-run` - Show planned operations without applying
- `--continue-on-error` - Continue processing other files on errors

**Examples:**
```bash
# Import everything found in directory
searchctl clone import --dir /backup

# Dry run import of ILM and pipelines
searchctl clone import --types lifecycle-policies,ingest-pipelines --dir /backup --dry-run
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
- `-f, --conditions-file` - file containing rollover conditions (JSON or YAML)

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

# Rollover using conditions file (JSON or YAML)
searchctl rollover datastream logs-system -f examples/rollover-conditions.json
searchctl rollover datastream logs-system -f examples/rollover-conditions.yaml

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

# Wildcard Deletion Implementation Notes

## API Limitations Discovered

During testing, we discovered that both Elasticsearch and OpenSearch have security restrictions on wildcard deletion:

### Elasticsearch
- **Setting**: `action.destructive_requires_name` (defaults to `true`)
- **Error**: "Wildcard expressions or all indices are not allowed"
- **Documentation**: "When set to `true`, you must specify the index name to delete an index. It is not possible to delete all indices with `_all` or use wildcards."

### OpenSearch  
- **Similar restrictions**: Also inherits the `action.destructive_requires_name` setting
- **Default behavior**: Same as Elasticsearch - prevents wildcard deletion by default

## Solution Implemented

Instead of requiring cluster administrators to modify security settings, we implemented **client-side wildcard expansion**:

1. **Pattern Detection**: Detect wildcard patterns (containing `*`)
2. **List Matching Resources**: Use GET APIs to list all indices/datastreams 
3. **Filter Matches**: Apply prefix matching for patterns ending with `*`
4. **Confirm Actions**: Show user exactly what will be deleted
5. **Individual Deletion**: Delete each resource one by one

## Benefits

- ✅ Works with default cluster security settings
- ✅ No special permissions required
- ✅ Clear visibility of what will be deleted
- ✅ Safe confirmation prompts with -y bypass
- ✅ Maintains backwards compatibility

## Test Results

All wildcard deletion features are working correctly:
- Single index/datastream deletion with confirmation
- Wildcard pattern detection and expansion  
- Interactive confirmation prompts
- -y flag bypass for automation
- Error handling for non-existent patterns
