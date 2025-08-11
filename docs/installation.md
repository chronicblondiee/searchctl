# Installation & Setup

This guide covers installing `searchctl` and configuring it to talk to your Elasticsearch or OpenSearch cluster.

## Prerequisites

- Go 1.24+ installed and on your PATH (`go version`)
- Git
- Optional: `make` (for building from source)

## Install Options

### Option A: Install with Go

```bash
# Installs the latest version to $GOBIN (or $GOPATH/bin if GOBIN is unset)
go install github.com/chronicblondiee/searchctl@latest

# Ensure the install directory is on PATH (Linux/macOS)
export PATH="$(go env GOPATH)/bin:$PATH"

# Verify
searchctl version
```

### Option B: Build from Source

```bash
git clone https://github.com/chronicblondiee/searchctl
cd searchctl

# Install to your Go bin (uses ldflags for version info)
make install

# Or build a local binary in ./bin
make build
./bin/searchctl version
```

## Initial Configuration

`searchctl` reads configuration from `~/.searchctl/config.yaml` by default.

Create the directory and file if they don’t exist:

```bash
mkdir -p ~/.searchctl
$EDITOR ~/.searchctl/config.yaml
```

Minimal local example:

```yaml
kind: Config
current-context: elasticsearch
contexts:
- name: elasticsearch
  context:
    cluster: elasticsearch-local
    user: default
clusters:
- name: elasticsearch-local
  cluster:
    server: http://localhost:9200
    insecure-skip-tls-verify: true
users:
- name: default
  user: {}
```

Security hardening:

```bash
chmod 600 ~/.searchctl/config.yaml
```

You can also point to a config file and/or force a context via environment variables:

- `SEARCHCTL_CONFIG` – override config file location
- `SEARCHCTL_CONTEXT` – override current context

Example:

```bash
export SEARCHCTL_CONFIG="examples/test-config.yaml"
export SEARCHCTL_CONTEXT="elasticsearch"
```

## Test the Connection

Using your config:

```bash
searchctl cluster health
```

Or by overriding the context just for a single command:

```bash
searchctl --context elasticsearch cluster info -o json
```

## Uninstall

- If installed via `go install`, remove the binary from `$GOBIN`/`$GOPATH/bin`.
- If built from source, delete the local `./bin/searchctl`.

## Troubleshooting

- Ensure `searchctl version` runs and prints version info
- Verify the config path and format (`~/.searchctl/config.yaml`)
- For TLS clusters, set `certificate-authority` and avoid `insecure-skip-tls-verify: true`
- Confirm network access to your cluster URL

For advanced configuration details, see `docs/configuration.md`.
