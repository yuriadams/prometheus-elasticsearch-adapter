package config

import (
	"io/ioutil"
	"log"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type config struct {
	ElasticsearchURL        string        `yaml:"elasticsearch.url"`
	ElasticsearchMaxRetries int           `yaml:"elasticsearch.max.retries"`
	ElasticIndexPerfix      string        `yaml:"elasticsearch.index.perfix"`
	ElasticType             string        `yaml:"elasticsearch.type"`
	RemoteTimeout           time.Duration `yaml:"web.timeout"`
	ListenAddr              string        `yaml:"web.listen.addr"`
	TelemetryPath           string        `yaml:"web.telemetry.path"`
}

func GetConfig() *config {
	cfg := &config{}
	yamlFile, err := ioutil.ReadFile("./config/config.yaml")

	if err != nil {
		log.Printf("yamlFile.Get err  #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return cfg
}
