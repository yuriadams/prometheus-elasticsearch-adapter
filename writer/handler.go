package writer

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/storage/remote"
	"github.com/yuriadams/prometheus-elasticsearch-adapter/config"
	"github.com/yuriadams/prometheus-elasticsearch-adapter/elasticsearch"
)

type writer interface {
	Write(samples model.Samples) error
	Name() string
}

// Handler receives the payload from Prometheus, format and send to Elasticsearch
func Handler(w http.ResponseWriter, r *http.Request) {
	compressed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req remote.WriteRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	samples := protoToSamples(&req)
	config.ReceivedSamples.Add(float64(len(samples)))

	go func(rw writer) {
		sendSamples(rw, samples)
	}(buildWriter())
}

func buildWriter() writer {
	cfg := config.GetConfig()
	var w writer

	if cfg.ElasticsearchURL != "" {
		url, err := url.Parse(cfg.ElasticsearchURL)
		if err != nil {
			log.Fatalf("Failed to parse Elasticsearch URL %q: %v", cfg.ElasticsearchURL, err)
		}
		c := elasticsearch.NewClient(url.String(), cfg.ElasticsearchMaxRetries,
			cfg.ElasticIndexPerfix, cfg.ElasticType, 30*time.Second)
		w = c
	}
	return w
}

func protoToSamples(req *remote.WriteRequest) model.Samples {
	var samples model.Samples
	for _, ts := range req.Timeseries {
		metric := make(model.Metric, len(ts.Labels))
		for _, l := range ts.Labels {
			metric[model.LabelName(l.Name)] = model.LabelValue(l.Value)
		}

		for _, s := range ts.Samples {
			samples = append(samples, &model.Sample{
				Metric:    metric,
				Value:     model.SampleValue(s.Value),
				Timestamp: model.Time(s.TimestampMs),
			})

		}
	}
	return samples
}

func sendSamples(w writer, samples model.Samples) {
	begin := time.Now()

	err := w.Write(samples)

	duration := time.Since(begin).Seconds()
	if err != nil {
		log.With("num_samples", len(samples)).With("storage", w.Name()).With("err", err).Warnf("Error sending samples to remote storage")
		config.FailedSamples.WithLabelValues(w.Name()).Add(float64(len(samples)))
	}
	config.SentSamples.WithLabelValues(w.Name()).Add(float64(len(samples)))
	config.SentBatchDuration.WithLabelValues(w.Name()).Observe(duration)
}
