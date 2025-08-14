# Logstash Filter Examples (8.19)

This directory contains Logstash pipelines that focus on filter examples for common log types. Inputs are minimal (often Beats or Syslog), and outputs default to Elasticsearch data streams to align with the templates in `examples/`.

Pipelines included:
- `pipeline-nginx-access.conf` – Nginx access logs (grok + useragent + geoip)
- `pipeline-apache-access.conf` – Apache access logs (grok + useragent + geoip)
- `pipeline-json-app.conf` – JSON app logs (json + mutate + date)
- `pipeline-syslog-auth.conf` – Syslog auth (grok for sshd)
- `pipeline-k8s-container.conf` – Kubernetes container logs (dissect + date)
- `pipeline-nginx-ingress.conf` – Nginx Ingress controller logs (grok + useragent + geoip)
- `pipeline-log-type-detector.conf` – Detects multiple formats (apache/nginx/ingress/syslog/k8s/json); applies parsing and tags with `log_type:*`
- `pipeline-log-level-detector.conf` – Detects JSON / key-value / line logs, extracts `level` from `message`, and drops DEBUG/debug
- `pipeline-unified-log-parser.conf` – Unified pipeline (JSON, KV, Serilog, line logs) with level extraction and normalization

Run:
```bash
logstash --path.settings examples/logstash/filter-examples --config.reload.automatic
```

Notes:
- Inputs use environment variables and can be swapped (e.g., Beats vs Kafka).
- Outputs are configured for data streams. Adjust `data_stream_dataset` to match your workload.
- For HTTPS outputs, set `ES_CACERT` and add `cacert => "${ES_CACERT}"` in the `elasticsearch {}` block if needed.
- For Elastic Cloud or Serverless, prefer API keys with `api_key => "${ES_API_KEY}"` on port 443.
