# Elasticsearch Remote storage adapter

Elastic write/read adapter for Prometheus remote storage, more details refer to: [Prometheus remote storage documentation](https://github.com/prometheus/prometheus/tree/master/documentation/examples/remote_storage/remote_storage_adapter)

It will receive prometheus samples and send batch requests to [Elastic](https://www.elastic.co/)

## dependencies using [dep](https://github.com/golang/dep)

```
dep init
```

## Building

```
go build
```

## Running

```
./prometheus-elasticsearch-adapter
```

## Running with Docker Compose

```
docker-compose up
```

config.yaml file:

```yaml
elasticsearch.url: http://localhost:9200
elasticsearch.max.retries: 1
elasticsearch.index.perfix: prometheus
elasticsearch.type: prom-metric
elasticsearch.aws.service: false #Indentifies if we are using AWS ElasticSearchService
web.listen.addr: :9201
web.telemetry.path: /metrics

# If we need to use AWS ElasticSearch Service, we must toggle the config 'elasticsearch.aws.service' to true e export your AWS Credentials: 
# export AWS_ACCESS_KEY=YourAccessKey
# export AWS_SECRET_KEY=YourSecretKey

```

## Configuring Prometheus

To configure Prometheus to send samples to this binary, add the following to your `prometheus.yml`:

```yaml
# Remote write configuration.
remote_write:
  - url: "http://app:9201/write"

# Remote read configuration
remote_read:
  - url: "http://app:9201/read"
```
