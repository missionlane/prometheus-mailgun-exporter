# prometheus-mailgun-exporter
Prometheus Mailgun Exporter

## Build
`make` will build both binary and Docker image

## Run
The exporter will serve metrics on `http://<ip>:9616/metrics`, and a healthz endpoint on `http://<ip>:9616/healthz` for use in Kubernetes

`export MG_API_KEY=<api_key>`

* Docker
  1. `docker run -ti --rm --name prometheus-mailgun-exporter -e MG_API_KEY -p 9616:9616 missionlane/prometheus-mailgun-exporter:latest`
* Binary
  1. `./prometheus-mailgun-exporter`

## Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `MG_API_KEY` | Mailgun API key (required) | |
| `API_BASE` | Mailgun API base URL | `https://api.mailgun.net` |

## Europe Users

`export API_BASE=https://api.eu.mailgun.net`

## Dashboard
A Grafana dashboard is included in [grafana/dashboard.json](grafana/dashboard.json).
