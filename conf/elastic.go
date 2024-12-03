package conf

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net"
	"net/http"
	"time"
)

func ElasticClient() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"},
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Optional: Cek koneksi dengan Info()
	_, err = esClient.Info()
	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %s", err)
	}

	fmt.Println("Elasticsearch connection established")
	return esClient
}
