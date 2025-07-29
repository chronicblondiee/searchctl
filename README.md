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

# Create resources  
searchctl create index my-logs                  # Create new index
searchctl create datastream logs-nginx          # Create data stream
searchctl create index test-idx --dry-run       # Preview creation

# Delete resources
searchctl delete index old-logs                 # Delete index
searchctl delete datastream temp-logs           # Delete data stream + backing indices
searchctl delete index test-* --dry-run         # Preview deletion

# Describe resources
searchctl describe index my-logs-2024.01        # Show detailed index info
searchctl describe index logs-nginx -o yaml     # Describe with YAML output
```

### Data Stream Management
```bash
# Rollover operations
searchctl rollover datastream logs --max-age 7d --dry-run       # Age-based rollover
searchctl rollover datastream logs --max-docs 1M                # Document count rollover
searchctl rollover datastream logs --max-size 50gb              # Size-based rollover
searchctl rollover datastream logs --max-primary-shard-size 25gb # Shard size rollover
searchctl rollover datastream logs --lazy --max-age 1d          # Lazy rollover
searchctl rollover datastream logs -f conditions.json          # Conditions from file
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
searchctl apply -f template.yaml               # Apply from file
searchctl apply -f config.json --dry-run       # Preview apply

# Version info
searchctl version                               # Show version
searchctl version -o json                      # Version as JSON
```

**Global Flags:** `--config`, `--context`, `--output` (table|json|yaml|wide), `--dry-run`, `--verbose`

## Configuration

Context-based cluster management in `~/.searchctl/config.yaml`:

```yaml
apiVersion: v1
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

## Documentation

- [Commands](docs/commands.md) - Complete command reference
- [Configuration](docs/configuration.md) - Advanced setup options  
- [Test Scripts](scripts/README.md) - Development testing guide