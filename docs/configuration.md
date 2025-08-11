# Configuration Guide

## Configuration File

`searchctl` uses a YAML configuration file located at `~/.searchctl/config.yaml` by default. The configuration supports multiple cluster contexts similar to kubectl.

## Configuration Structure

```yaml
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

## Configuration Sections

### Root Level
- `kind` - Always `Config`
- `current-context` - Default context to use

### Contexts
Define named combinations of cluster and user configurations.

```yaml
contexts:
- name: "context-name"
  context:
    cluster: "cluster-name"
    user: "user-name"
```

### Clusters
Define connection details for search clusters.

```yaml
clusters:
- name: "cluster-name"
  cluster:
    server: "https://cluster.example.com:9200"
    certificate-authority: "/path/to/ca.crt"         # Optional
    insecure-skip-tls-verify: false                  # Optional
```

**Cluster Options:**
- `server` - Elasticsearch/OpenSearch endpoint URL (required)
- `certificate-authority` - Path to CA certificate file
- `insecure-skip-tls-verify` - Skip TLS certificate verification (default: false)

### Users
Define authentication credentials.

```yaml
users:
- name: "user-name"
  user:
    username: "elastic"           # Basic auth username
    password: "password"          # Basic auth password
    api-key: "base64-key"         # API key authentication
```

**Authentication Methods:**

1. **Basic Authentication:**
   ```yaml
   user:
     username: "elastic"
     password: "changeme"
   ```

2. **API Key Authentication:**
   ```yaml
   user:
     api-key: "VnVhQ2ZHY0JDZGJrUW0tZTVhT3g6dWkybHAyYXhUTm1zeWFrdzl0dk5udw=="
   ```

3. **No Authentication:**
   ```yaml
   user: {}
   ```

## Environment Variables

Override configuration values using environment variables:

- `SEARCHCTL_CONFIG` - Override config file location
- `SEARCHCTL_CONTEXT` - Override current context
- `SEARCHCTL_SERVER` - Override cluster server URL
- `SEARCHCTL_USERNAME` - Override username
- `SEARCHCTL_PASSWORD` - Override password
- `SEARCHCTL_API_KEY` - Override API key

## Command-Line Overrides

Override configuration using command-line flags:

```bash
# Use different config file
searchctl --config /path/to/config.yaml get indices

# Use different context
searchctl --context staging get nodes

# Override server URL
searchctl --server https://other-cluster:9200 cluster health
```

## Multiple Clusters

Manage multiple clusters by defining separate contexts:

```yaml
current-context: "local"
contexts:
- name: "local"
  context:
    cluster: "local-es"
    user: "local-user"
- name: "staging"
  context:
    cluster: "staging-es"
    user: "staging-user"
- name: "production"
  context:
    cluster: "prod-es"
    user: "prod-user"

clusters:
- name: "local-es"
  cluster:
    server: "http://localhost:9200"
    insecure-skip-tls-verify: true
- name: "staging-es"
  cluster:
    server: "https://staging-es.company.com:9200"
- name: "prod-es"
  cluster:
    server: "https://prod-es.company.com:9200"
    certificate-authority: "/etc/ssl/certs/ca.crt"

users:
- name: "local-user"
  user: {}
- name: "staging-user"
  user:
    username: "elastic"
    password: "staging-password"
- name: "prod-user"
  user:
    api-key: "production-api-key"
```

## Security Best Practices

1. **File Permissions:** Ensure config file has restricted permissions:
   ```bash
   chmod 600 ~/.searchctl/config.yaml
   ```

2. **API Keys:** Prefer API keys over username/password when possible

3. **TLS Verification:** Enable TLS verification in production:
   ```yaml
   insecure-skip-tls-verify: false
   ```

4. **Certificate Authorities:** Use proper CA certificates:
   ```yaml
   certificate-authority: "/path/to/ca.crt"
   ```

5. **Environment Variables:** Use environment variables for sensitive data in CI/CD:
   ```bash
   export SEARCHCTL_API_KEY="your-api-key"
   searchctl get indices
   ```

## Configuration Examples

### Local Development
```yaml
kind: Config
current-context: "local"
contexts:
- name: "local"
  context:
    cluster: "local"
    user: "local"
clusters:
- name: "local"
  cluster:
    server: "http://localhost:9200"
    insecure-skip-tls-verify: true
users:
- name: "local"
  user: {}
```

### Production with TLS
```yaml
kind: Config
current-context: "production"
contexts:
- name: "production"
  context:
    cluster: "prod"
    user: "prod"
clusters:
- name: "prod"
  cluster:
    server: "https://elasticsearch.company.com:9200"
    certificate-authority: "/etc/ssl/certs/es-ca.crt"
users:
- name: "prod"
  user:
    api-key: "VnVhQ2ZHY0JDZGJrUW0tZTVhT3g6dWkybHAyYXhUTm1zeWFrdzl0dk5udw=="
```
