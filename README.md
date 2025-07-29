# searchctl

A kubectl-like CLI for Elasticsearch and OpenSearch cluster management.

## Quick Start

```bash
# Install from source
git clone https://github.com/chronicblondiee/searchctl
cd searchctl && make build

# Basic usage
searchctl cluster health
searchctl get indices
searchctl --context production get nodes -o json
```

## Core Features

- **Cluster Management**: Health checks, node information, configuration
- **Resource Operations**: Create, read, update, delete indices and data streams  
- **Rollover Management**: Advanced data stream rollover with multiple conditions
- **Multi-Engine Support**: Works with both Elasticsearch and OpenSearch
- **Output Formats**: Table, JSON, YAML with wide/verbose modes
- **Context Management**: Switch between clusters seamlessly

## Commands

| Command | Purpose | Examples |
|---------|---------|----------|
| `get` | List resources | `get indices`, `get nodes`, `get datastreams` |
| `create` | Create resources | `create index my-index`, `create datastream logs` |
| `delete` | Delete resources | `delete index old-data`, `delete datastream logs` |
| `describe` | Resource details | `describe index my-index` |
| `rollover` | Data stream rollover | `rollover ds logs --max-age 7d --max-docs 1M` |
| `cluster` | Cluster operations | `cluster health`, `cluster info` |
| `apply` | Apply configurations | `apply -f template.yaml` |
| `config` | Manage contexts | `config view`, `config use-context prod` |

## Rollover Features

Advanced data stream rollover with comprehensive condition support:

```bash
# Age-based rollover
searchctl rollover datastream logs --max-age 30d

# Size-based rollover  
searchctl rollover datastream logs --max-size 50gb --max-primary-shard-size 25gb

# Document count rollover
searchctl rollover datastream logs --max-docs 1000000 --max-primary-shard-docs 500000

# Lazy rollover (mark for rollover at next write)
searchctl rollover datastream logs --lazy --max-age 1d

# Multiple conditions with file
searchctl rollover datastream logs -f rollover-conditions.json

# Dry-run mode
searchctl rollover datastream logs --dry-run --max-age 7d
```

**Global Flags:** `--config`, `--context`, `--output` (table|json|yaml|wide), `--dry-run`, `--verbose`

## Configuration

Manages multiple clusters via contexts in `~/.searchctl/config.yaml`:

```yaml
apiVersion: v1
kind: Config
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
  cluster: {server: "http://localhost:9200", insecure-skip-tls-verify: true}
users:
- name: admin
  user: {username: elastic, password: changeme}
- name: default
  user: {}
```

Switch contexts: `searchctl config use-context development`

## Development

**Quick Setup:**
```bash
# Start test environment (Elasticsearch + OpenSearch)
./scripts/start-test-env.sh

# Run tests
./scripts/integration-test.sh       # Basic functionality
./scripts/test-rollover.sh          # Rollover features  
./scripts/test-rollover-real.sh     # Real operations
./scripts/test-performance.sh       # Performance tests

# Cleanup
./scripts/stop-test-env.sh
```

**Build & Test:**
```bash
make build test
go test ./...
```

## Documentation

- [Command Reference](docs/commands.md) - Complete command documentation
- [Configuration Guide](docs/configuration.md) - Advanced configuration options  
- [Test Scripts](scripts/README.md) - Development and testing guide