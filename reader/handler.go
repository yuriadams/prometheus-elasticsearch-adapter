package reader

import "github.com/prometheus/prometheus/storage/remote"

type Reader interface {
	Read(req *remote.ReadRequest) (*remote.ReadResponse, error)
	Name() string
}

// http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
// 	compressed, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	reqBuf, err := snappy.Decode(nil, compressed)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	var req remote.ReadRequest
// 	if err := proto.Unmarshal(reqBuf, &req); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	// TODO: Support reading from more than one Reader and merging the results.
// 	if len(readers) != 1 {
// 		http.Error(w, fmt.Sprintf("expected exactly one Reader, found %d readers", len(readers)), http.StatusInternalServerError)
// 		return
// 	}
//
// 	reader := readers[0]
//
// 	var resp *remote.ReadResponse
// 	resp, err = reader.Read(&req)
// 	if err != nil {
// 		log.With("query", req).With("storage", reader.Name()).With("err", err).Warnf("Error executing query")
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	data, err := proto.Marshal(resp)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
//
// 	w.Header().Set("Content-Type", "application/x-protobuf")
// 	w.Header().Set("Content-Encoding", "snappy")
//
// 	compressed = snappy.Encode(nil, data)
// 	if _, err := w.Write(compressed); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// })
