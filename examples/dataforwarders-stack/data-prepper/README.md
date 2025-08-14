Data Prepper examples (OpenSearch 3.1)

This directory mirrors the Logstash examples using OpenSearch Data Prepper pipelines. It includes HTTP, Kafka, syslog-like, and stdin/beats-like pipelines that send events to an OpenSearch index. These examples favor secure defaults, environment-driven configuration, and simple JSON payloads.

Requirements
- Docker (recommended) or a local Data Prepper install
- OpenSearch 2.x/3.x endpoint reachable from Data Prepper

Files
- `pipelines.yaml`: Multiple named pipelines (`http_ingest`, `kafka`, `syslog_like`, `stdin_like`, `beats_like`)

Run Data Prepper with Docker

Use environment variables for OpenSearch connection details:

```bash
export OPENSEARCH_HOSTS="https://localhost:9200"
export OPENSEARCH_USERNAME="admin"
export OPENSEARCH_PASSWORD="admin"

docker run --rm \
  -p 2021:2021 \# http_ingest
  -p 2022:2022 \# beats_like
  -p 2023:2023 \# stdin_like
  -p 2024:2024 \# syslog_like
  -v "$(pwd)/pipelines.yaml:/usr/share/data-prepper/pipelines/pipelines.yaml:ro" \
  -e OPENSEARCH_HOSTS -e OPENSEARCH_USERNAME -e OPENSEARCH_PASSWORD \
  opensearchproject/data-prepper:latest
```

Notes
- The container image looks for `pipelines/pipelines.yaml` by default; we mount our file there.
- These examples use HTTP sources for easy local testing. For production, enable TLS and auth on sources and avoid exposing them publicly.
- If you use AWS/OpenSearch Service, prefer secrets managers/IAM for credentials and private networking.

Testing the HTTP pipelines

- http_ingest pipeline (POST JSON):
```bash
curl -sS -X POST "http://localhost:2021" \
  -H 'Content-Type: application/json' \
  -d '{"message":"hello http","level":"info"}'
```

- stdin_like pipeline (POST JSON):
```bash
curl -sS -X POST "http://localhost:2023" \
  -H 'Content-Type: application/json' \
  -d '{"message":"hello stdin","src":"cli"}'
```

- syslog_like pipeline (send plain text syslog line):
```bash
curl -sS -X POST "http://localhost:2024" \
  -H 'Content-Type: text/plain' \
  --data-binary "<34>1 2024-01-01T12:00:00Z host app 1234 ID47 [exampleSDID@32473 iut=3 eventSource=Application] starting"
```

- beats_like pipeline (Filebeat over HTTP): configure Filebeat to use `output.http` and point it at `http://localhost:2022` with `Content-Type: application/json`.

Example Filebeat output (beats-like via HTTP)

```yaml
output.http:
  hosts: ["http://localhost:2022"]
  # Optional auth/TLS here
  # username: ${OP_HTTP_USER}
  # password: ${OP_HTTP_PASS}
  headers:
    Content-Type: application/json
```

Kafka pipeline

Set environment values and run a Kafka broker locally or point to your cluster.

```bash
export KAFKA_BOOTSTRAP_SERVERS=localhost:9092
export KAFKA_TOPICS=logs
export KAFKA_GROUP_ID=data-prepper
```

OpenSearch index naming

Each example writes to a simple index (e.g., `logs-http-ingest`, `logs-kafka-generic`). Adjust to your naming conventions. If using OpenSearch data streams, adapt the sink to your template/index strategy.

Security & best practices
- Do not hardcode credentials. Use environment variables or secrets.
- Prefer TLS for sources and OpenSearch connections; validate certificates.
- Lock down source listeners (bind to private addresses, use auth).
- Apply least-privilege roles for the OpenSearch user.

References
- Data Prepper: https://docs.opensearch.org/latest/data-prepper/
- OpenSearch sink config: https://docs.opensearch.org/latest/data-prepper/pipelines/configuration/sinks/opensearch/


