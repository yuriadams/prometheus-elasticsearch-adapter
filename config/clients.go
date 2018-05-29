package config

import (
	"log"
	"net/url"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/storage/remote"
	"github.com/yuriadams/prometheus-elasticsearch-adapter/elasticsearch"
)

type Writer interface {
	Write(samples model.Samples) error
	Name() string
}

type Reader interface {
	Read(req *remote.ReadRequest) ([]map[string]interface{}, error)
	Name() string
}

func BuildClient() (Writer, Reader) {
	cfg := GetConfig()
	var w Writer
	var r Reader

	if cfg.ElasticsearchURL != "" {
		url, err := url.Parse(cfg.ElasticsearchURL)
		if err != nil {
			log.Fatalf("Failed to parse Elasticsearch URL %q: %v", cfg.ElasticsearchURL, err)
		}
		c := elasticsearch.NewClient(url.String(), cfg.ElasticsearchMaxRetries,
			cfg.ElasticIndexPerfix, cfg.ElasticType, 30*time.Second)
		w = c
		r = c
	}
	return w, r
}
