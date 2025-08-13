# Logstash Filter Examples

This directory contains Logstash pipelines that focus on filter examples for common log types. Inputs are minimal (often Beats or Syslog), and outputs default to Elasticsearch data streams to align with the templates in `examples/`.

Pipelines included:
- `pipeline-nginx-access.conf` – Nginx access logs (grok + useragent + geoip)
- `pipeline-apache-access.conf` – Apache access logs (grok + useragent + geoip)
- `pipeline-json-app.conf` – JSON app logs (json + mutate + date)
- `pipeline-syslog-auth.conf` – Syslog auth (grok for sshd)
- `pipeline-k8s-container.conf` – Kubernetes container logs (dissect + date)
- `pipeline-nginx-ingress.conf` – Nginx Ingress controller logs (grok + useragent + geoip)

Run:
```bash
logstash --path.settings examples/logstash/filter-examples --config.reload.automatic
```

Notes:
- Inputs use environment variables and can be swapped (e.g., Beats vs Kafka).
- Outputs are configured for data streams. Adjust `data_stream_dataset` to match your workload.
