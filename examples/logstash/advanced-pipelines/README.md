# Advanced Multi-Pipeline Examples

These examples demonstrate Logstash pipeline-to-pipeline (P2P) routing to build a modular pipeline graph.

Topology:

```
           +---------+
           | errors  |  -> logs-ingest.errors data stream
           +----^----+
                |
                | (on parse failure)
+---------+   send_to   +-------+   send_to   +--------+   send_to   +--------+
| ingest  | ----------> | route | ---------> | enrich | --------->  | output |
+----^----+             +---^---+            +---^----+             +---^----+
     |                      |                    |                      |
  Beats/HTTP input     dataset routing        UA/GeoIP            ES data streams
```

Pipelines:
- `ingest`: Accepts events from inputs (Beats/HTTP), does light normalization, routes errors to `errors` and all other events to `route`.
- `route`: Classifies events and sets `data_stream.*` fields (e.g., `nginx.access`, `kubernetes.container`, `app.generic` fallback) based on content; forwards to `enrich`.
- `enrich`: Adds `useragent` and `geoip` enrichment when fields are present; forwards to `output`.
- `output`: Writes to Elasticsearch data streams using event `data_stream.*` fields.
- `errors`: Receives parsing failures and writes them to a dedicated `logs-ingest.errors` dataset.

Directory layout:
- `pipelines.yml`: wires the five pipelines together
- `pipeline-ingest.conf`: inputs + initial parsing + routing
- `pipeline-route.conf`: dataset classification
- `pipeline-enrich.conf`: enrichments (UA/GeoIP)
- `pipeline-output.conf`: Elasticsearch data stream output
- `pipeline-errors.conf`: error sink output

Prerequisites:
- Logstash 7.14+
- Plugins: `logstash-output-elasticsearch`, `logstash-filter-geoip`, `logstash-filter-useragent`
- Apply templates (recommended):
  - `searchctl apply -f examples/component-templates/base-settings.yaml`
  - `searchctl apply -f examples/component-templates/observability-mappings.yaml`
  - `searchctl apply -f examples/datastream-template.yaml`

Environment variables:
- `ES_HOSTS` (default: `https://localhost:9200`)
- `ES_USERNAME` (default: `elastic`)
- `ES_PASSWORD` (default: `changeme`)
- `DATA_STREAM_NAMESPACE` (default: `default`)
- `BEATS_PORT` (default: `5044`)
- `HTTP_HOST` (default: `0.0.0.0`), `HTTP_PORT` (default: `8081`)

Run:
```bash
logstash --path.settings examples/logstash/advanced-pipelines --config.reload.automatic
```

Quick test:
```bash
# Send an HTTP JSON log
curl -sS -X POST "http://localhost:8081/" \
  -H 'Content-Type: application/json' \
  -d '{"message":"127.0.0.1 - - [10/Jan/2025:13:55:36 +0000] \"GET / HTTP/1.1\" 200 123 \"-\" \"curl/8.0\"","verb":"GET","httpversion":"1.1"}'

# Or ship via Filebeat to Beats input at ${BEATS_PORT}
```

Notes:
- Outputs are data-stream based and derive `data_stream.*` from upstream routing. Adjust routing rules in `pipeline-route.conf`.
- On parse failures, events are routed to `pipeline-errors.conf` and indexed into `logs-ingest.errors`.
- Consider enabling persistent queues for durability in production.

## How this configuration works

1) Ingest stage (`pipeline-ingest.conf`)
- Accepts events from Beats and HTTP.
- Attempts to parse `message` as JSON (non-destructive on invalid JSON).
- Tags with `ingest_stage: ingest`.
- Conditional routing:
  - If `_jsonparsefailure` tag exists, send to `errors`.
  - Otherwise, send to `route`.

2) Routing stage (`pipeline-route.conf`)
- Tags with `route_stage: route`.
- Classifies the event and sets `data_stream.*` fields:
  - Kubernetes hints present → `kubernetes.container` dataset
  - HTTP-ish fields present (e.g., `verb`, `httpversion`) → `nginx.access` dataset
  - Fallback → `app.generic` dataset
- Forwards all events to `enrich`.

3) Enrichment stage (`pipeline-enrich.conf`)
- Tags with `enrich_stage: enrich`.
- User agent parsing if `agent` exists → stores in `user_agent`.
- GeoIP lookup if `source.ip` is present (or it renames `clientip` → `source.ip`).
- Forwards all events to `output`.

4) Output stage (`pipeline-output.conf`)
- Uses Elasticsearch output with data streams.
- Reads `data_stream.type`, `data_stream.dataset`, and `data_stream.namespace` from the event set upstream.

5) Error stage (`pipeline-errors.conf`)
- Receives events that failed early parsing.
- Sets `data_stream.*` to `logs / ingest.errors / ${DATA_STREAM_NAMESPACE}` and writes to Elasticsearch.

## Pipeline-to-pipeline mechanics

- Upstream pipelines use the `pipeline` output with `send_to` to forward events:
  - Example: `pipeline { send_to => "route" }`
- Downstream pipelines use the `pipeline` input with `address` to receive:
  - Example: `input { pipeline { address => "route" } }`
- Names must match exactly and be defined as `pipeline.id` entries in `pipelines.yml`.
- You can quote names (recommended) and send to multiple downstreams using arrays, e.g. `send_to => ["enrich", "audit"]`.

## Extending the topology

- Fan-out processing:
  - From `route`, send to multiple pipelines, e.g. `enrich` and `audit` for compliance redaction.
  - Create `pipeline-audit.conf` and wire it in `pipelines.yml` with `input { pipeline { address => "audit" } }`.
- Output isolation:
  - Add separate output pipelines (e.g., `output-es`, `output-s3`) and enable `queue.type: persisted` per pipeline in `pipelines.yml`.
  - This prevents a slow/broken output from blocking the rest of the graph.
- Dataset expansion:
  - Extend `pipeline-route.conf` with additional conditions (e.g., `apache.access`, `db.audit`).
  - Keep datasets stable and low-cardinality.

## Troubleshooting

- No events in downstream pipeline:
  - Ensure `send_to`/`address` names match and exist in `pipelines.yml`.
- Events not in expected data stream:
  - Verify `pipeline-route.conf` sets `data_stream.*` as intended and that later stages do not overwrite them.
- GeoIP or useragent missing:
  - Confirm the relevant fields exist (`agent`, `source.ip`) and required plugins are installed.
- Backpressure or blocking:
  - Consider isolating outputs and enabling persisted queues per-output pipeline in `pipelines.yml`.
