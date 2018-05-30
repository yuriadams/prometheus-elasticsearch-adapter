package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	elastic "github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/storage/remote"
	awsauth "github.com/smartystreets/go-aws-auth"
)

type awsSigningTransport struct {
	HTTPClient  *http.Client
	Credentials awsauth.Credentials
}

type errNoDataPointsFound struct {
	message string
}

func newErrNoDataPointsFound(message string) *errNoDataPointsFound {
	return &errNoDataPointsFound{
		message: message,
	}
}

func (e *errNoDataPointsFound) Error() string {
	return e.message
}

// Client allows sending batches of Prometheus samples to ElasticSearch.
type Client struct {
	client         *elastic.Client
	esIndex        string
	esType         string
	timeout        time.Duration
	ignoredSamples prometheus.Counter
}

// fieldsFromMetric extracts Elastic fields from a Prometheus metric.
// In elasticsearch, `__name__` could also be a simple field.
func fieldsFromMetric(m model.Metric) map[string]interface{} {
	fields := make(map[string]interface{}, len(m))
	for l, v := range m {
		fields[string(l)] = string(v)
	}
	return fields
}

func generateEsIndex(esIndexPerfix string) string {
	var separator = "-"
	dateSuffix := time.Now().Format("2006-01-02")
	return esIndexPerfix + separator + dateSuffix
}

// RoundTrip implementation
func (a awsSigningTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return a.HTTPClient.Do(awsauth.Sign4(req, a.Credentials))
}

// NewClient returns a Client which contains an elasticsearch client.
// Now, it generates the real esIndex formatted with `<esIndexPerfix>-YYYY-mm-dd`.
func NewClient(url string, maxRetries int, esIndexPerfix, esType string, timeout time.Duration, awsService bool) *Client {
	ctx := context.Background()
	var client *elastic.Client
	var err error
	if awsService {
		signingTransport := awsSigningTransport{
			Credentials: awsauth.Credentials{
				AccessKeyID:     os.Getenv("AWS_ACCESS_KEY"),
				SecretAccessKey: os.Getenv("AWS_SECRET_KEY"),
			},
			HTTPClient: http.DefaultClient,
		}

		signingClient := &http.Client{
			Transport: http.RoundTripper(signingTransport),
		}

		client, err = elastic.NewClient(
			elastic.SetURL(url),
			elastic.SetScheme("https"),
			elastic.SetHttpClient(signingClient),
			elastic.SetSniff(false),
		)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		client, err = elastic.NewClient(
			elastic.SetURL(url),
			elastic.SetMaxRetries(maxRetries),
			// TODO: add basic auth support.
		)

		if err != nil {
			log.Fatal(err)
		}
	}

	// Use the IndexExists service to check if a specified index exists.
	esIndex := generateEsIndex(esIndexPerfix)

	exists, err := client.IndexExists(esIndex).Do(ctx)
	if err != nil {
		log.Errorf("index %v is not found in Elastic.", esIndex)
	}

	// Create an index if it is not exist.
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(esIndex).Do(ctx)
		if err != nil {
			log.Fatalf("failed to create index %v, err: %v", esIndex, err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
			log.Fatalf("elasticsearch no acknowledged.")
		}
	}

	return &Client{
		client:  client,
		esIndex: esIndex,
		esType:  esType,
		timeout: timeout,
		ignoredSamples: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "prometheus_elasticsearch_ignored_samples_total",
				Help: "The total number of samples not sent to Elasticsearch due to unsupported float values (Inf, -Inf, NaN).",
			},
		),
	}
}

// Write sends a batch of samples to Elasticsearch.
func (c *Client) Write(samples model.Samples) error {
	ctx := context.Background()

	bulkRequest := c.client.Bulk().Timeout(c.timeout.String())

	for _, s := range samples {
		v := float64(s.Value)
		if math.IsNaN(v) || math.IsInf(v, 0) {
			log.Debugf("cannot send value %f to Elasticsearch, skipping sample %#v", v, s)
			c.ignoredSamples.Inc()
			continue
		}

		document := fieldsFromMetric(s.Metric)
		document["value"] = v
		document["timestamp"] = s.Timestamp.Time()
		documentJSON, err := json.Marshal(document)
		if err != nil {
			log.Debugf("error while marshaling document, err: %v", err)
			continue
		}

		indexRq := elastic.NewBulkIndexRequest().
			Index(c.esIndex).
			Type(c.esType).
			Doc(string(documentJSON))
		bulkRequest = bulkRequest.Add(indexRq)
	}

	bulkResponse, err := bulkRequest.Do(ctx)
	if err != nil {
		return err
	}

	// there are some failed requests, count it!
	failedResults := bulkResponse.Failed()
	log.Debugf("there are %d failed requests to Elasticsearch.", len(failedResults))
	for _ = range failedResults {
		c.ignoredSamples.Inc()
	}
	return nil
}

// Read queries metrics from Elasticsearch.
func (c *Client) Read(req *remote.ReadRequest) ([]map[string]interface{}, error) {
	ctx := context.Background()
	querier := req.Queries[0]
	query := elastic.NewBoolQuery()

	for _, matcher := range querier.Matchers {
		query = query.Must(elastic.NewTermQuery(matcher.Name, matcher.Value))
	}

	// building elasticsearch query
	query = query.Filter(elastic.
		NewRangeQuery("timestamp").
		From(querier.StartTimestampMs).
		To(querier.EndTimestampMs))

	searchResult, err := c.client.
		Search().
		Index(c.esIndex).
		Query(query).
		Size(1000).
		Sort("timestamp", true).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	var dataPoints []map[string]interface{}

	// parsing elasticsearch hits to array with datapoints
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d data points\n", searchResult.Hits.TotalHits)

		for _, hit := range searchResult.Hits.Hits {
			var dataPoint map[string]interface{}
			json.Unmarshal(*hit.Source, &dataPoint)
			dataPoints = append(dataPoints, dataPoint)
		}
	} else {
		return nil, newErrNoDataPointsFound("Found no metrics")
	}

	return dataPoints, nil
}

// Name identifies the client as an elasticsearch client.
func (c Client) Name() string {
	return "elasticsearch"
}

// Describe implements prometheus.Collector.
func (c *Client) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ignoredSamples.Desc()
}

// Collect implements prometheus.Collector.
func (c *Client) Collect(ch chan<- prometheus.Metric) {
	ch <- c.ignoredSamples
}
