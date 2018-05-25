package reader

import (
	"io/ioutil"
	"net/http"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/prometheus/storage/remote"
	"github.com/yuriadams/prometheus-elasticsearch-adapter/config"
)

type reader interface {
	Read(req *remote.ReadRequest) (*remote.ReadResponse, error)
	Name() string
}

// Handle receives the payload from Elasticsearch, format and send to Prometheus
func Handle(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(config.ReadDuration)
	defer timer.ObserveDuration()

	compressed, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.With("err", err).Error("Failed to read body.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		log.With("err", err).Error("Failed to decompress body.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req remote.ReadRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		log.With("err", err).Error("Failed to unmarshal body.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Queries) != 1 {
		log.Error("More than one query sent.")
		http.Error(w, "Can only handle one query.", http.StatusBadRequest)
		return
	}

	log.Info("*******>>>>>KEURI<<<<<**********")
	log.Info(req.Queries[0])

	// result, err := ca.runQuery(req.Queries[0])
	// if err != nil {
	// 	log.With("err", err).Error("Failed to run select against Crate.")
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// resp := remote.ReadResponse{
	// 	Results: []*remote.QueryResult{
	// 		{Timeseries: result},
	// 	},
	// }
	// data, err := proto.Marshal(&resp)
	// if err != nil {
	// 	log.With("err", err).Error("Failed to marshal response.")
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// w.Header().Set("Content-Type", "application/x-protobuf")
	// if _, err := w.Write(snappy.Encode(nil, data)); err != nil {
	// 	log.With("err", err).Error("Failed to compress response.")
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}
