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

- Filebeat 7.14+ (for stable data stream behavior)
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

## Run

Local file inputs:
```bash
sudo filebeat -e -c examples/filebeat/filebeat.yml
```

Kubernetes autodiscover:
```bash
sudo filebeat -e -c examples/filebeat/filebeat-kubernetes-autodiscover.yml
```

Docker host:
```bash
sudo filebeat -e -c examples/filebeat/filebeat-docker.yml
```

Syslog receiver:
```bash
sudo filebeat -e -c examples/filebeat/filebeat-syslog.yml
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
