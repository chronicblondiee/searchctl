# kubectl-style Architecture Implementation

This document describes the new kubectl-inspired architecture for the searchctl Elasticsearch client.

## Architecture Overview

The new architecture follows kubectl's clientset pattern with clear separation of concerns:

```
pkg/client/
├── clientset.go          # Main clientset interface
├── client.go             # Backward compatibility wrapper
├── factory.go            # Configuration factory
├── rest/                 # HTTP transport layer
│   └── client.go
├── cluster/              # Cluster operations
│   ├── interface.go
│   └── cluster.go
├── indices/              # Index operations
│   ├── interface.go
│   └── indices.go        # Includes template operations
├── datastreams/          # Data stream operations
│   ├── interface.go
│   └── datastreams.go
├── nodes/                # Node operations
│   ├── interface.go
│   └── nodes.go
└── types/                # Shared types
    └── types.go
```

## Usage Examples

### kubectl-style Resource-Centric API

```go
// Create clientset
clientset, err := client.NewClientset()
if err != nil {
    log.Fatal(err)
}

// Resource-first approach (like kubectl)
health, err := clientset.Cluster().Health()
indices, err := clientset.Indices().List("logs-*")
index, err := clientset.Indices().Get("specific-index")
templates, err := clientset.Indices().Templates().List("*")
dataStreams, err := clientset.DataStreams().List("metrics-*")
nodes, err := clientset.Nodes().List()
```

### Backward Compatibility

```go
// Old interface still works
oldClient, err := client.NewClient()
if err != nil {
    log.Fatal(err)
}

health, err := oldClient.ClusterHealth()
indices, err := oldClient.GetIndices("logs-*")
```

## Key Benefits

### 1. **Resource-Centric Design**
Operations are organized around Elasticsearch resources (cluster, indices, data streams, nodes) rather than mixed into a single client.

### 2. **Hierarchical Organization**
Similar to kubectl's structure:
- `clientset.Cluster().Health()`
- `clientset.Indices().Templates().List(pattern)`
- `clientset.DataStreams().Rollover(name, conditions, lazy)`

### 3. **Clean Separation of Concerns**
- **Factory**: Configuration and client creation
- **REST Client**: HTTP transport abstraction  
- **Resource Clients**: Domain-specific operations
- **Types**: Shared data structures

### 4. **Extensibility**
Easy to add new resource types:
```go
// Add new resource client
type SearchInterface interface {
    Query(index string, query map[string]interface{}) (*SearchResult, error)
}

// Extend clientset
type Interface interface {
    Cluster() cluster.Interface
    Indices() indices.Interface
    DataStreams() datastreams.Interface
    Nodes() nodes.Interface
    Search() search.Interface  // New addition
}
```

### 5. **Testability**
Each component can be unit tested independently:
- Mock REST client for transport testing
- Mock resource interfaces for business logic testing
- Separate concerns enable focused tests

## Design Patterns Used

### 1. **Factory Pattern**
`NewFactory()` handles configuration complexity and provides clean initialization.

### 2. **Interface Segregation**
Each resource has its own focused interface rather than one monolithic interface.

### 3. **Composition over Inheritance**
Clientset composes resource clients rather than inheriting everything.

### 4. **Adapter Pattern**
The old `Client` struct adapts the new clientset to maintain backward compatibility.

## Migration Guide

### For End Users
No changes required - existing code continues to work unchanged.

### For New Development
Use the new clientset API:
```go
// Old way
client, _ := client.NewClient()
indices, _ := client.GetIndices("*")

// New way
clientset, _ := client.NewClientset()
indices, _ := clientset.Indices().List("*")
```

### Adding New Features
1. Create new resource package (e.g., `pkg/client/security/`)
2. Define interface and implementation
3. Add to main clientset interface
4. Optionally add backward compatibility methods to old client

This architecture provides a solid foundation for future growth while maintaining backward compatibility and following proven patterns from the kubectl ecosystem.