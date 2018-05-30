package config

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Config for the app
type Config struct {
	ElasticsearchURL        string        `yaml:"elasticsearch.url"`
	ElasticsearchMaxRetries int           `yaml:"elasticsearch.max.retries"`
	ElasticIndexPerfix      string        `yaml:"elasticsearch.index.perfix"`
	ElasticType             string        `yaml:"elasticsearch.type"`
	RemoteTimeout           time.Duration `yaml:"web.timeout"`
	ListenAddr              string        `yaml:"web.listen.addr"`
	TelemetryPath           string        `yaml:"web.telemetry.path"`
	AwsElasticSearchService bool          `yaml:"elasticsearch.aws.service"`
}

// GetConfig returns the app's configuration described on config.yaml on root
func GetConfig() *Config {
	cfg := &Config{}
	yamlFile, err := ioutil.ReadFile(os.Getenv("CONFIG_PATH"))

	if err != nil {
		log.Printf("yamlFile.Get err  #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return cfg
}
