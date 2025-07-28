# searchctl

A kubectl-like CLI tool for OpenSearch and Elasticsearch cluster management.

## Installation

```bash
# From source
git clone https://github.com/chronicblondiee/searchctl
cd searchctl && make build

# Using Go
go install github.com/chronicblondiee/searchctl@latest
```

## Quick Start

```bash
# Check cluster health
searchctl cluster health

# List indices
searchctl get indices

# Use different context
searchctl --context production get nodes -o json
```

## Commands

- `get` - List resources (indices, nodes)
- `describe` - Show detailed resource information  
- `create/delete` - Manage resources
- `cluster` - Cluster operations (health, info)
- `config` - Configuration management

**Global Flags:** `--config`, `--context`, `--output` (table|json|yaml), `--dry-run`

See [docs/commands.md](docs/commands.md) for complete reference.

## Configuration

Config stored in `~/.searchctl/config.yaml`. See [examples/test-config.yaml](examples/test-config.yaml).

## Documentation

- [Commands Reference](docs/commands.md)
- [Configuration Guide](docs/configuration.md)

## Development

**Test Environment:**
```bash
podman-compose up -d  # Start Elasticsearch + OpenSearch
./test-both-engines.sh  # Run tests
```

**Build & Test:**
```bash
make build test
```