# Project Plan: searchctl

A kubectl-like CLI tool for OpenSearch and Elasticsearch cluster management.

## 1. Project Overview

`searchctl` will be a command-line interface for managing OpenSearch and Elasticsearch clusters, designed to provide a familiar kubectl-like experience. The tool will treat both OpenSearch and Elasticsearch APIs uniformly, allowing users to manage either type of cluster with the same commands.

### Key Technologies
- **Go**: Primary programming language
- **Cobra**: CLI framework for command structure and parsing
- **Viper**: Configuration management
- **go-elasticsearch**: Official Elasticsearch Go client (compatible with OpenSearch)

## 2. Project Structure

```
searchctl/
├── cmd/                    # Cobra commands
│   ├── root.go            # Root command setup
│   ├── get/               # Get commands
│   │   ├── get.go
│   │   ├── indices.go
│   │   ├── nodes.go
│   │   ├── shards.go
│   │   └── templates.go
│   ├── describe/          # Describe commands
│   │   ├── describe.go
│   │   ├── index.go
│   │   └── node.go
│   ├── create/            # Create commands
│   │   ├── create.go
│   │   └── index.go
│   ├── delete/            # Delete commands
│   │   ├── delete.go
│   │   └── index.go
│   ├── apply.go           # Apply configurations from files
│   ├── config.go          # Config management commands
│   └── cluster.go         # Cluster-wide operations
├── pkg/
│   ├── client/            # Elasticsearch/OpenSearch client wrapper
│   │   ├── client.go
│   │   ├── indices.go
│   │   ├── nodes.go
│   │   └── cluster.go
│   ├── config/            # Configuration management
│   │   ├── config.go
│   │   └── context.go
│   ├── output/            # Output formatting
│   │   ├── formatter.go
│   │   ├── table.go
│   │   ├── json.go
│   │   └── yaml.go
│   └── util/              # Utility functions
│       ├── validation.go
│       └── helpers.go
├── internal/              # Internal packages
│   └── version/
│       └── version.go
├── examples/              # Example configuration files
│   ├── index-template.yaml
│   └── config.yaml
├── docs/                  # Documentation
│   ├── commands.md
│   └── configuration.md
├── .goreleaser.yml        # GoReleaser configuration
├── Makefile              # Build automation
├── go.mod
├── go.sum
└── main.go
```

## 3. Configuration Management (Viper)

### Configuration File Structure
Location: `~/.searchctl/config.yaml`

```yaml
apiVersion: v1
kind: Config
current-context: "production"
contexts:
- name: "production"
  context:
    cluster: "prod-cluster"
    user: "prod-user"
- name: "development"
  context:
    cluster: "dev-cluster"
    user: "dev-user"
clusters:
- name: "prod-cluster"
  cluster:
    server: "https://es-prod.example.com:9200"
    certificate-authority: "/path/to/ca.crt"
    insecure-skip-tls-verify: false
- name: "dev-cluster"
  cluster:
    server: "http://localhost:9200"
    insecure-skip-tls-verify: true
users:
- name: "prod-user"
  user:
    username: "elastic"
    password: "password"
- name: "dev-user"
  user:
    api-key: "base64-encoded-api-key"
```

### Configuration Features
- Multiple cluster contexts
- Different authentication methods (basic auth, API keys, certificates)
- Environment variable overrides
- Command-line flag overrides

## 4. CLI Command Structure (Cobra)

### Core Commands

#### Get Commands
```bash
searchctl get indices [flags]
searchctl get nodes [flags]
searchctl get shards [index-name] [flags]
searchctl get templates [flags]
searchctl get aliases [flags]
```

#### Describe Commands
```bash
searchctl describe index <index-name> [flags]
searchctl describe node <node-name> [flags]
searchctl describe template <template-name> [flags]
```

#### Create Commands
```bash
searchctl create index <index-name> [flags]
searchctl create template <template-name> -f <file> [flags]
```

#### Delete Commands
```bash
searchctl delete index <index-name> [flags]
searchctl delete template <template-name> [flags]
```

#### Apply Commands
```bash
searchctl apply -f <file> [flags]
searchctl apply -f <directory> [flags]
```

#### Cluster Commands
```bash
searchctl cluster-info [flags]
searchctl cluster-health [flags]
searchctl cluster-settings [flags]
```

#### Config Commands
```bash
searchctl config view [flags]
searchctl config use-context <context-name>
searchctl config set-cluster <cluster-name> --server=<server-url>
searchctl config set-credentials <user-name> --username=<username> --password=<password>
```

### Global Flags
- `--config`: Specify config file location
- `--context`: Override current context
- `--output, -o`: Output format (table, json, yaml, wide)
- `--namespace, -n`: Specify index pattern/namespace
- `--verbose, -v`: Verbose output
- `--dry-run`: Show what would be done without executing

## 5. API Client Abstraction

### Client Interface
```go
type SearchClient interface {
    // Cluster operations
    ClusterHealth() (*ClusterHealth, error)
    ClusterInfo() (*ClusterInfo, error)
    
    // Index operations
    GetIndices(pattern string) ([]Index, error)
    GetIndex(name string) (*Index, error)
    CreateIndex(name string, body map[string]interface{}) error
    DeleteIndex(name string) error
    
    // Node operations
    GetNodes() ([]Node, error)
    GetNode(nodeID string) (*Node, error)
    
    // Template operations
    GetTemplates() ([]Template, error)
    CreateTemplate(name string, template Template) error
    DeleteTemplate(name string) error
}
```

### Implementation Strategy
- Use `github.com/elastic/go-elasticsearch/v8` as the base client
- Create a unified wrapper that abstracts OpenSearch/Elasticsearch differences
- Handle version detection and API compatibility
- Implement connection pooling and retry logic
- Support multiple authentication methods

## 6. Output Formatting

### Format Types
- **Table**: Default human-readable format (similar to kubectl)
- **JSON**: Machine-readable JSON output
- **YAML**: YAML format for configuration files
- **Wide**: Extended table format with additional columns

### Implementation
- Use `text/tabwriter` for table formatting
- Implement custom JSON/YAML marshalers for clean output
- Support color coding for status indicators
- Implement pagination for large result sets

## 7. Error Handling and Validation

### Error Categories
- Configuration errors (missing config, invalid context)
- Connection errors (network, authentication)
- API errors (invalid requests, server errors)
- Validation errors (invalid resource names, malformed inputs)

### Implementation
- Standardized error messages with suggestions
- Exit codes following Unix conventions
- Contextual help messages
- Debug mode for troubleshooting

## 8. Development Phases

### Phase 1: Foundation (Week 1-2)
- Project setup and directory structure
- Basic Cobra command structure
- Viper configuration management
- Basic client wrapper
- Simple `get indices` and `cluster-health` commands

### Phase 2: Core Commands (Week 3-4)
- Complete `get` command family
- Implement `describe` commands
- Add output formatting options
- Error handling and validation

### Phase 3: Advanced Operations (Week 5-6)
- `create` and `delete` commands
- `apply` command with file support
- Template management
- Configuration management commands

### Phase 4: Polish and Distribution (Week 7-8)
- Comprehensive testing
- Documentation
- CI/CD pipeline
- Release automation with GoReleaser
- Performance optimization

## 9. Testing Strategy

### Unit Tests
- Client wrapper functionality
- Command parsing and validation
- Output formatting
- Configuration management

### Integration Tests
- End-to-end command execution
- Real cluster interactions (with test clusters)
- Authentication flows
- Error scenarios

### Test Structure
```
tests/
├── unit/
│   ├── client_test.go
│   ├── config_test.go
│   └── output_test.go
├── integration/
│   ├── commands_test.go
│   └── auth_test.go
└── fixtures/
    ├── responses/
    └── configs/
```

## 10. Build and Distribution

### Makefile Targets
```makefile
.PHONY: build test install clean release

build:
	go build -o bin/searchctl .

test:
	go test ./...

install:
	go install .

clean:
	rm -rf bin/

release:
	goreleaser release --rm-dist

dev-deps:
	go install github.com/goreleaser/goreleaser@latest
```

### CI/CD Pipeline
- GitHub Actions for automated testing
- Multi-platform builds (Linux, macOS, Windows)
- Automated releases with semantic versioning
- Docker image creation for containerized usage

### Distribution Methods
- GitHub Releases with binaries
- Homebrew formula for macOS
- APT/YUM repositories for Linux
- Docker Hub for containerized usage
- Go modules for library usage

## 11. Future Enhancements

### Advanced Features
- Index lifecycle management
- Snapshot and restore operations
- Performance monitoring and alerting
- Query execution and profiling
- Plugin system for custom commands

### Integration Possibilities
- Kubernetes operator integration
- Helm chart management
- Monitoring system integration (Prometheus, Grafana)
- CI/CD pipeline integration

This plan provides a comprehensive roadmap for building `searchctl` as a production-ready CLI tool for OpenSearch and Elasticsearch management.