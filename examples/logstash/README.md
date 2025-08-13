# Logstash Pipeline Examples

This directory contains ready-to-use Logstash pipeline configurations for common inputs and for sending events to Elasticsearch or OpenSearch. They are designed to work well with the data stream and index templates in `examples/` and with clusters managed via `searchctl`.

## What's Included

- `pipelines.yml` – Example multi-pipeline configuration
- `pipeline-filebeat.conf` – Input from Filebeat (Beats protocol)
- `pipeline-kafka.conf` – Input from Kafka topic(s)
- `pipeline-syslog.conf` – Syslog over UDP/TCP
- `pipeline-http.conf` – HTTP ingestion (JSON)
- `pipeline-stdin.conf` – Simple STDIN for testing

All pipelines default to Elasticsearch data streams. You can switch to index-based output or OpenSearch by adjusting the output sections or environment variables.

## Prerequisites

- Logstash 7.14+ (for `elasticsearch` output data stream options)
- Plugins: `logstash-output-elasticsearch` (and optionally `logstash-output-opensearch` if targeting OpenSearch)

## Environment Variables

These pipelines use environment variables with sensible defaults. Override them at runtime with `LS_JAVA_OPTS` or a Logstash `.env` file.

- Elasticsearch
  - `ES_HOSTS` (default: `https://localhost:9200`)
  - `ES_USERNAME` (default: `elastic`)
  - `ES_PASSWORD` (default: `changeme`)
  - `DATA_STREAM_NAMESPACE` (default: `default`)

- Kafka
  - `KAFKA_BOOTSTRAP_SERVERS` (default: `localhost:9092`)
  - `KAFKA_TOPICS` (default: `logs`)
  - `KAFKA_GROUP_ID` (default: `logstash`)

- Beats
  - `BEATS_PORT` (default: `5044`)

- Syslog
  - `SYSLOG_UDP_PORT` (default: `5140`)
  - `SYSLOG_TCP_PORT` (default: `5140`)

- HTTP
  - `HTTP_HOST` (default: `0.0.0.0`)
  - `HTTP_PORT` (default: `8080`)

## Running

1) Apply templates and policies using `searchctl` (optional but recommended):

```bash
searchctl apply -f examples/component-templates/base-settings.yaml
searchctl apply -f examples/component-templates/observability-mappings.yaml
searchctl apply -f examples/datastream-template.yaml
```

2) Start Logstash with these pipelines:

```bash
logstash -f examples/logstash/pipeline-filebeat.conf \
  --path.settings examples/logstash \
  -r
```

Or use the multi-pipeline file:

```bash
logstash --path.settings examples/logstash --config.reload.automatic
```

3) Send data:

- Filebeat: Point Filebeat output to Logstash (`hosts: ["localhost:5044"]`).
- Kafka: Produce JSON records to the configured topic(s).
- Syslog: Send syslog messages to UDP/TCP `5140`.
- HTTP: `POST` JSON to `http://localhost:8080/`.

## Output Modes

These examples default to Elasticsearch data streams:

```text
data_stream.type: logs | metrics | traces
data_stream.dataset: e.g., app.generic
data_stream.namespace: default (override with DATA_STREAM_NAMESPACE)
```

To use classic indices instead, comment the data stream section and uncomment the `index =>` setting in the `elasticsearch` output blocks.

For OpenSearch, either install the Logstash OpenSearch output plugin and switch the output to `opensearch { ... }`, or rely on index-based ingestion.

## Notes

- The filters are conservative and aim to pass through JSON logs when possible. Extend with your own `grok`, `json`, and `date` processors as needed for your log sources.
- If events already contain `@timestamp`, the `date` filter will respect it; otherwise, ingestion time is used.


