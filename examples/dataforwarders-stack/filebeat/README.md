# Filebeat Example Configurations

This directory contains ready-to-use Filebeat configurations for common environments, designed to work with Elasticsearch/OpenSearch clusters managed by `searchctl` and the templates in `examples/`.

## What's Included

- `filebeat.yml` – Base config for local/file inputs with Elasticsearch output
- `filebeat-kubernetes-autodiscover.yml` – Kubernetes autodiscover with hints
- `filebeat-docker.yml` – Docker container logs on a single host
- `filebeat-syslog.yml` – Syslog receiver (TCP/UDP)
- `modules.d/system.yml` – Example system module
- `modules.d/nginx.yml` – Example Nginx module

All examples default to Elasticsearch data streams. You can switch to Logstash output or classic indices as needed.

## Prerequisites

- Filebeat 8.19+ (recommended). Examples use `filestream`, autodiscover hints, and HTTP monitoring compatible with 8.19.
- Apply example templates/policies with `searchctl` (recommended):

```bash
searchctl apply -f examples/component-templates/base-settings.yaml
searchctl apply -f examples/component-templates/observability-mappings.yaml
searchctl apply -f examples/datastream-template.yaml
```

## Environment Variables

- `ES_HOSTS` (default: `https://localhost:9200`)
- `ES_USERNAME` (default: `elastic`)
- `ES_PASSWORD` (default: `changeme`)
- `DATA_STREAM_TYPE` (default: `logs`)
- `DATA_STREAM_DATASET` (default: `filebeat.generic`)
- `DATA_STREAM_NAMESPACE` (default: `default`)
- `LOGSTASH_HOSTS` (default: empty; when set, you can switch output to Logstash)

## Usage

1) Pick the example that matches your environment

- `filebeat.yml`: Read local files on a single host
- `filebeat-docker.yml`: Read Docker container logs on a single host
- `filebeat-kubernetes-autodiscover.yml`: Discover pods and read container logs in Kubernetes
- `filebeat-syslog.yml`: Receive syslog over UDP/TCP (default port 5514)

2) Set environment variables (optional)

```bash
export ES_HOSTS="https://localhost:9200"
export ES_USERNAME="elastic"
export ES_PASSWORD="changeme"
# For Logstash instead, set LOGSTASH_HOSTS and switch the output section in the config
# export LOGSTASH_HOSTS="logstash.example:5044"
```

3) Validate configuration and connectivity

```bash
sudo filebeat test config -c <config>
sudo filebeat test output -c <config>
```

4) Run Filebeat

```bash
sudo filebeat -e -c <config>
```

Notes
- HTTP monitoring endpoint is enabled at `http://<bind>:5066` in these examples. Protect it in production.
- Example outputs target Elasticsearch. To use Logstash, comment out `output.elasticsearch` and uncomment `output.logstash` in the chosen config.
- Modules: example module configs are in `modules.d/`. Enable or adjust as needed.

## Run

Local file inputs:
```bash
sudo filebeat -e -c examples/filebeat/filebeat.yml
# Validate config and outputs first
sudo filebeat test config -c examples/filebeat/filebeat.yml
sudo filebeat test output -c examples/filebeat/filebeat.yml
```

Kubernetes autodiscover:
```bash
sudo filebeat -e -c examples/filebeat/filebeat-kubernetes-autodiscover.yml
sudo filebeat test config -c examples/filebeat/filebeat-kubernetes-autodiscover.yml
sudo filebeat test output -c examples/filebeat/filebeat-kubernetes-autodiscover.yml
```

Docker host:
```bash
sudo filebeat -e -c examples/filebeat/filebeat-docker.yml
sudo filebeat test config -c examples/filebeat/filebeat-docker.yml
sudo filebeat test output -c examples/filebeat/filebeat-docker.yml
```

Syslog receiver:
```bash
sudo filebeat -e -c examples/filebeat/filebeat-syslog.yml
sudo filebeat test config -c examples/filebeat/filebeat-syslog.yml
sudo filebeat test output -c examples/filebeat/filebeat-syslog.yml
```

Sending a quick test message to the syslog receiver:
```bash
echo "<14>1 $(date -u +%Y-%m-%dT%H:%M:%SZ) host app - - - hello from udp" | nc -u -w1 127.0.0.1 5514
```

## Data Streams

These configs set `data_stream.*` fields and write to indices of the form:
```
logs-%{[data_stream.dataset]}-%{[data_stream.namespace]}
```
With a data-stream-enabled index template applied, indexing to that name creates/uses a data stream automatically.

To use classic indices, change the `output.elasticsearch.index` accordingly and/or remove the `data_stream` fields.

## OpenSearch

Standard Filebeat ships with `elasticsearch` and `logstash` outputs. For OpenSearch, prefer the Logstash output (use the Logstash examples in `examples/logstash`) or use the OpenSearch Beats fork. Classic index output may also work when compatible.

## Notes for 8.19

- Inputs use `filestream` with `container` parser for Docker/Kubernetes for proper CRI Docker JSON decoding.
- Conditional `decode_json_fields` avoids corrupting plain-text logs.
- Syslog receiver supports UDP and TCP. Set host/ports via `SYSLOG_UDP_HOST`, `SYSLOG_UDP_PORT`, `SYSLOG_TCP_HOST`, `SYSLOG_TCP_PORT`.
- HTTP monitoring endpoint is enabled on port 5066. Protect it in production.
