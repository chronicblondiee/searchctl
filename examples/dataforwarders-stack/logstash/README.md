# Logstash Pipeline Examples (8.19)

This directory contains ready-to-use Logstash pipeline configurations for common inputs and for sending events to Elasticsearch or OpenSearch. They are designed to work well with the data stream and index templates in `examples/` and with clusters managed via `searchctl`. Examples have been validated against Logstash 8.19.

## What's Included

- `pipelines.yml` – Example multi-pipeline configuration
- `pipeline-filebeat.conf` – Input from Filebeat (Beats protocol)
- `pipeline-kafka.conf` – Input from Kafka topic(s)
- `pipeline-syslog.conf` – Syslog over UDP/TCP
- `pipeline-http.conf` – HTTP ingestion (JSON)
- `pipeline-stdin.conf` – Simple STDIN for testing

All pipelines default to Elasticsearch data streams. You can switch to index-based output or OpenSearch by adjusting the output sections or environment variables.

## Prerequisites

- Logstash 8.19 (bundled JDK 17 is recommended)
- Plugins: `logstash-output-elasticsearch` (and optionally `logstash-output-opensearch` if targeting OpenSearch)

## Environment Variables

These pipelines use environment variables with sensible defaults. Override them at runtime with `LS_JAVA_OPTS` or a Logstash `.env` file.

- Elasticsearch
  - `ES_HOSTS` (default: `https://localhost:9200`)
  - `ES_USERNAME` (default: `elastic`)
  - `ES_PASSWORD` (default: `changeme`)
  - `ES_API_KEY` (optional; use instead of username/password for Elastic Cloud)
  - `ES_CACERT` (optional; path to CA cert for HTTPS clusters)
  - `DATA_STREAM_NAMESPACE` (default: `default`)
  - `DATA_STREAM_TYPE` (default: `logs`)
  - `DATA_STREAM_DATASET` (pipeline-specific defaults; override as needed)

- Kafka
  - `KAFKA_BOOTSTRAP_SERVERS` (default: `localhost:9092`)
  - `KAFKA_TOPICS` (default: `logs`)
  - `KAFKA_GROUP_ID` (default: `logstash`)

- Beats
  - `BEATS_PORT` (default: `5044`)

- Syslog
  - `SYSLOG_UDP_PORT` (default: `5140`)

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
  --config.reload.automatic
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

## Connecting to Elasticsearch securely

- When using `https://` endpoints, provide the cluster CA certificate when required by your environment. One approach is to set an environment variable and reference it in your output (example shown, but not enabled by default in these configs):
  - `ES_CACERT=/path/to/http_ca.crt`
  - In `elasticsearch { ... }`, add `cacert => "${ES_CACERT}"`.
- For Elastic Cloud (including Serverless), prefer API keys and port `443`:
  - Set `ES_HOSTS=https://<your-deployment-id>.<region>.aws.elastic-cloud.com:443`
  - Set `ES_API_KEY=<id:api_key>` and remove username/password. Note: API key and basic auth are mutually exclusive.
- Data streams do not use ILM; DLM is managed on the Elasticsearch side. The examples set `ilm_enabled => false` explicitly for clarity.

## Notes

- The filters are conservative and aim to pass through JSON logs when possible. Extend with your own `grok`, `json`, and `date` processors as needed for your log sources.
- If events already contain `@timestamp`, the `date` filter will respect it; otherwise, ingestion time is used.
- Syslog example uses the `syslog` input on a non-privileged UDP port by default. If you require TCP syslog, add a `tcp { port => ... }` input and parse with grok/`syslog_pri` as needed.

## Production recommendations

- Enable persistent queues per pipeline when durability is required:
  - In `pipelines.yml`, set `queue.type: persisted` for output-heavy pipelines.
- Tune pipeline performance based on your host:
  - `pipeline.workers` defaults to CPU cores; increase if CPU is underutilized.
  - `pipeline.batch.size` defaults to 125; larger values may improve throughput at the cost of memory.
- Docker: mount pipeline configs to `/usr/share/logstash/pipeline/` and settings to `/usr/share/logstash/config/`. The default `http.host` in containers is `0.0.0.0`.


