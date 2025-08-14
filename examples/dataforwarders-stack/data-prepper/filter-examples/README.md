Data Prepper filter examples

This directory mirrors the Logstash filter examples using Data Prepper processors. Each file contains a single pipeline showcasing parsing and enrichment for a specific log type.

How to run

Mount the desired file as `pipelines.yaml` when running Data Prepper. Example for Apache access:

```bash
docker run --rm \
  -p 2121:2121 \
  -v "$(pwd)/filters-apache-access.yaml:/usr/share/data-prepper/pipelines/pipelines.yaml:ro" \
  -e OPENSEARCH_HOSTS -e OPENSEARCH_USERNAME -e OPENSEARCH_PASSWORD \
  opensearchproject/data-prepper:latest
```

Examples
- `filters-apache-access.yaml`: Apache combined log → grok + rename `clientip` → user_agent → geoip
- `filters-json-app.yaml`: JSON parse from `message` → date normalize → lowercase level → rename host
- `filters-k8s-container.yaml`: dissect CRI line → date → JSON parse inner → remove temp fields
- `filters-nginx-access.yaml`: nginx access → grok + user_agent + geoip
- `filters-nginx-ingress.yaml`: nginx ingress → grok + user_agent + geoip
- `filters-syslog-auth.yaml`: syslog auth lines → grok → date
- `filters-log-type-detector.yaml`: detects common formats (apache/nginx/nginx-ingress/syslog/k8s/json) and applies per-type parsing, adds `log_type:*` tag

Combined file
- `pipelines.yaml` includes all `filters_*` pipelines for convenience. Expose all ports: 2121–2126, 2130.


