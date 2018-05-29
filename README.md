# Elasticsearch Remote storage adapter

Elastic write adapter for Prometheus remote storage, more details refer to: [Prometheus remote storage documentation](https://github.com/prometheus/prometheus/tree/master/documentation/examples/remote_storage/remote_storage_adapter)

It will receive prometheus samples and send batch requests to [Elastic](https://www.elastic.co/)

## Building

```
go build
```

## Running

```
./prometheus-elasticsearch-adapter
```

config.yaml file:

```yaml
elasticsearch.url: http://localhost:9200
elasticsearch.max.retries: 1
elasticsearch.index.perfix: prometheus
elasticsearch.type: prom-metric
web.listen.addr: :9201
web.telemetry.path: /metrics
```

## Configuring Prometheus

To configure Prometheus to send samples to this binary, add the following to your `prometheus.yml`:

```yaml
# Remote write configuration.
remote_write:
  - url: "http://localhost:9201/write"

# Remote read configuration
remote_read:
  - url: "http://localhost:9201/read"
```
