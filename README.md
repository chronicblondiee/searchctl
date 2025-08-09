# searchctl

A kubectl-like CLI for Elasticsearch and OpenSearch cluster management.

## Quick Start

```bash
# Install and basic usage
git clone https://github.com/chronicblondiee/searchctl && cd searchctl && make build
searchctl cluster health
searchctl get indices --context production -o json
```

## Commands

### Basic Operations
```bash
# Get resources
searchctl get indices                           # List all indices
searchctl get indices logs-*                    # List indices matching pattern
searchctl get nodes                             # List cluster nodes
searchctl get datastreams                       # List data streams
searchctl get datastreams logs-* -o json        # List with JSON output
searchctl get index-templates                   # List index templates
searchctl get idx-templates                     # Same as above (alias)
searchctl get component-templates               # List component templates  
searchctl get ct                                # Same as above (short alias)
searchctl get component-templates base-* -o yaml # List matching component templates
searchctl get lifecycle-policies                # List lifecycle policies (ILM/ISM)
searchctl get ilm                               # Same as above (short alias)
searchctl get lp                                # Same as above (shortest alias)

# Create resources  
searchctl create index my-logs                  # Create new index
searchctl create datastream logs-nginx          # Create data stream
searchctl create index test-idx --dry-run       # Preview creation

# Delete resources
searchctl delete index old-logs                 # Delete index
searchctl delete datastream temp-logs           # Delete data stream + backing indices
searchctl delete index test-* --dry-run         # Preview deletion
searchctl delete index-template old-template    # Delete index template
searchctl delete template old-template          # Same as above (alias)
searchctl delete component-template old-ct -y   # Delete component template (auto-confirm)  
searchctl delete ct old-ct -y                   # Same as above (short alias)
searchctl delete lifecycle-policy old-policy    # Delete lifecycle policy (ILM/ISM)
searchctl delete ilm old-policy -y              # Same as above (auto-confirm)
searchctl delete lp old-policy -y               # Same as above (shortest alias)

# Describe resources
searchctl describe index my-logs-2024.01                    # Index details
searchctl describe index logs-nginx -o yaml                 # Index as YAML
searchctl describe lifecycle-policy my-policy --show-body   # ILM/ISM policy
searchctl describe lifecycle-policy my-policy -o yaml       # Policy as YAML
searchctl describe index-template logs-template --show-body # Composable template
searchctl describe component-template base-settings --show-body
searchctl describe datastream logs-app -o yaml              # Data stream details
searchctl describe node node-1 -o yaml                      # Node details
```

### Data Stream Management
```bash
# Rollover operations
searchctl rollover datastream logs --max-age 7d --dry-run         # Age-based rollover
searchctl rollover datastream logs --max-docs 1M                  # Document count rollover
searchctl rollover datastream logs --max-size 50gb                # Size-based rollover
searchctl rollover datastream logs --max-primary-shard-size 25gb  # Shard size rollover
searchctl rollover datastream logs --lazy --max-age 1d            # Lazy rollover
searchctl rollover datastream logs -f examples/rollover-conditions.json  # From JSON file
searchctl rollover datastream logs -f examples/rollover-conditions.yaml  # From YAML file
```

### Cluster Operations
```bash
# Cluster status
searchctl cluster health                        # Show cluster health
searchctl cluster info                          # Show cluster information
searchctl cluster health -o json                # Health as JSON
```

### Configuration Management
```bash
# Configuration
searchctl config view                           # Show current config
searchctl config use-context production         # Switch context
searchctl config view -o yaml                   # Config as YAML

# Apply configurations
searchctl apply -f template.yaml               # Apply index template from file
searchctl apply -f component-template.yaml     # Apply component template from file
searchctl apply -f lifecycle-policy.yaml       # Apply lifecycle policy from file
searchctl apply -f config.json --dry-run       # Preview apply

# Version info
searchctl version                               # Show version
searchctl version -o json                      # Version as JSON
```

**Global Flags:** `--config`, `--context`, `--output` (table|json|yaml|wide), `--dry-run`, `--verbose`

### Quick Reference - Template Aliases

**Index Templates:**
- `searchctl get idx-templates` or `searchctl get it` - List index templates
- `searchctl delete template <name>` or `searchctl delete it <name>` - Delete index template

**Component Templates:**
- `searchctl get ct` - List component templates  
- `searchctl delete ct <name>` - Delete component template

**Lifecycle Policies:**
- `searchctl get ilm` or `searchctl get lp` - List lifecycle policies
- `searchctl delete ilm <name>` or `searchctl delete lp <name>` - Delete lifecycle policy

## Configuration

### Template Management

**Index Templates** - Define index settings, mappings, and aliases for index patterns:
```bash
# Apply index template from file
searchctl apply -f examples/index-template.yaml
searchctl get idx-templates logs-*
searchctl delete index-template logs-template
```

**Component Templates** - Reusable building blocks for index templates:
```bash
# Apply component templates
searchctl apply -f examples/component-templates/base-settings.yaml
searchctl apply -f examples/component-templates/observability-mappings.yaml  # Elasticsearch
searchctl apply -f examples/component-templates/observability-mappings-opensearch.yaml  # OpenSearch

# Manage component templates
searchctl get component-templates
searchctl get ct base-settings -o yaml             # Short alias
searchctl delete component-template base-settings -y
searchctl delete ct base-settings -y               # Short alias
```

**Lifecycle Policy Management** - Automated index lifecycle management:
```bash
# Apply lifecycle policies  
searchctl apply -f examples/lifecycle-policies/basic-ilm-policy.yaml         # Elasticsearch ILM
searchctl apply -f examples/lifecycle-policies/basic-ism-policy.yaml         # OpenSearch ISM

# Manage lifecycle policies
searchctl get lifecycle-policies
searchctl get ilm my-policy -o yaml                    # Short alias
searchctl delete lifecycle-policy my-policy -y
searchctl delete lp my-policy -y                       # Shortest alias
```

**Engine Compatibility:**
- **Elasticsearch**: Uses Index Lifecycle Management (ILM) with `/ilm/policy` endpoints
- **OpenSearch**: Uses Index State Management (ISM) with `/_plugins/_ism/policies` endpoints  
- **Component Templates**: Elasticsearch supports `flattened` field types; OpenSearch uses `object` fields with `enabled: false`

Context-based cluster management in `~/.searchctl/config.yaml`:

```yaml
current-context: production
contexts:
- name: production
  context: {cluster: prod-es, user: admin}
- name: development  
  context: {cluster: dev-es, user: default}
clusters:
- name: prod-es
  cluster: {server: "https://prod.elastic.com:9200"}
- name: dev-es
  cluster: {server: "http://localhost:9200"}
users:
- name: admin
  user: {username: elastic, password: changeme}
```

Use: `searchctl config use-context development`

## Development

```bash
# Quick test environment setup
./scripts/start-test-env.sh
./scripts/integration-test.sh
./scripts/stop-test-env.sh

# Build and test
make build test
```

## Examples

The `examples/` directory contains ready-to-use templates:

- **Index Templates**: `index-template.yaml`, `datastream-template.yaml`, etc.
- **Component Templates**: `component-templates/base-settings.yaml`, `observability-mappings*.yaml`
- **Lifecycle Policies**: `lifecycle-policies/basic-ilm-policy.yaml`, `hot-warm-cold-policy.yaml`, etc.
- **Configuration**: Sample cluster configurations for different environments

Example workflow:
```bash
# Apply component templates (building blocks)
searchctl apply -f examples/component-templates/base-settings.yaml
searchctl apply -f examples/component-templates/observability-mappings.yaml

# Apply lifecycle policies for index management
searchctl apply -f examples/lifecycle-policies/basic-ilm-policy.yaml

# Apply composable index template that uses component templates
searchctl apply -f examples/component-templates/composable-datastream-template.yaml

# List templates using aliases
searchctl get idx-templates                     # List all index templates
searchctl get ct                                # List all component templates
searchctl get lp                                # List all lifecycle policies

# Create data stream using the template
searchctl create datastream observability-logs

# Clean up using aliases
searchctl delete template my-template -y        # Delete index template
searchctl delete ct base-settings -y            # Delete component template  
searchctl delete lp basic-log-rotation -y       # Delete lifecycle policy
```

## Documentation

- [Commands](docs/commands.md) - Complete command reference
- [Configuration](docs/configuration.md) - Advanced setup options  
- [Test Scripts](scripts/README.md) - Development testing guide
- [Component Templates](examples/component-templates/README.md) - Template examples and compatibility guide
- [Lifecycle Policies](examples/lifecycle-policies/README.md) - ILM and ISM policy examples