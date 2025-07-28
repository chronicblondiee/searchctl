# searchctl

A kubectl-like CLI tool for OpenSearch and Elasticsearch cluster management.

## Overview

`searchctl` provides a familiar kubectl-style interface for managing OpenSearch and Elasticsearch clusters. It supports multiple cluster contexts, various authentication methods, and comprehensive cluster operations.

## Installation

### From Source
```bash
git clone https://github.com/chronicblondiee/searchctl
cd searchctl
make build
sudo cp bin/searchctl /usr/local/bin/
```

### Using Go
```bash
go install github.com/chronicblondiee/searchctl@latest
```

## Quick Start

1. **Initialize configuration:**
   ```bash
   searchctl config view
   ```

2. **Check cluster health:**
   ```bash
   searchctl cluster health
   ```

3. **List indices:**
   ```bash
   searchctl get indices
   ```

## Commands

### Core Operations
- `searchctl get` - List resources (indices, nodes)
- `searchctl describe` - Show detailed resource information
- `searchctl create` - Create resources
- `searchctl delete` - Delete resources
- `searchctl apply -f` - Apply configurations from files

### Cluster Operations
- `searchctl cluster health` - Show cluster health
- `searchctl cluster info` - Show cluster information

### Configuration
- `searchctl config view` - Display current configuration
- `searchctl config use-context` - Switch contexts

## Global Flags

- `--config` - Specify config file location
- `--context` - Override current context  
- `--output, -o` - Output format (table|json|yaml|wide)
- `--verbose, -v` - Verbose output
- `--dry-run` - Show what would be done without executing

## Configuration

Configuration is stored in `~/.searchctl/config.yaml`. See [examples/config.yaml](examples/config.yaml) for a complete example.

## Examples

```bash
# List all indices
searchctl get indices

# Get specific index pattern
searchctl get indices logs-*

# Describe an index
searchctl describe index my-index

# Create an index
searchctl create index new-index

# Delete an index (with dry-run)
searchctl delete index old-index --dry-run

# Output as JSON
searchctl get nodes -o json

# Use different context
searchctl get indices --context production
```

## Documentation

- [Commands Reference](docs/commands.md)
- [Configuration Guide](docs/configuration.md)

## License

See [LICENSE](LICENSE) file.