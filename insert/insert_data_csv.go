package insert

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

/*
	DATA UJI KNN
	- SCORE
	- RANK
	- POPULARITY
	- MEMBER

*/

func InsertDataCSVToElastic() {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
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
	})

	err = AnimeIndex(es)
	if err != nil {
		log.Fatal(err)
	}

	//if err != nil {
	//	log.Fatal(err)
	//}

	// Buka file CSV
	file, err := os.Open("processed_anime_dataset.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Baca file CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Loop untuk memasukkan data ke Elasticsearch
	for i, record := range records {
		if i == 0 {
			continue // Lewati header
		}

		// Parsing data dari CSV
		animeID, _ := strconv.Atoi(record[0])
		score, _ := strconv.ParseFloat(record[4], 64)
		rank, _ := strconv.ParseFloat(record[18], 64)
		popularity, _ := strconv.Atoi(record[19])
		favorites, _ := strconv.Atoi(record[20])
		members, _ := strconv.Atoi(record[22])

		// Membuat data dalam format JSON untuk Elasticsearch
		doc := map[string]interface{}{
			"anime_id":     animeID,
			"name":         record[1],
			"english_name": record[2],
			"other_name":   record[3],
			"score":        score,
			"genres":       strings.Split(record[5], ","),
			"synopsis":     record[6],
			"type":         record[7],
			"episodes":     record[8],
			"aired":        record[9],
			"premiered":    record[10],
			"status":       record[11],
			"producers":    strings.Split(record[12], ","),
			"licensors":    strings.Split(record[13], ","),
			"studios":      strings.Split(record[14], ","),
			"source":       record[15],
			"duration":     record[16],
			"rating":       record[17],
			"rank":         rank,
			"popularity":   popularity,
			"favorites":    favorites,
			"scored_by":    record[21],
			"members":      members,
			"image_url":    record[23],
			"embedding":    []float64{score, rank, float64(popularity), float64(members)},
		}

		// Konversi data ke JSON
		jsonData, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Error marshalling document: %s", err)
			continue
		}

		// Index data ke Elasticsearch
		req := esapi.IndexRequest{
			Index:      "anime_info", // Ganti dengan nama index Anda
			DocumentID: strconv.Itoa(animeID),
			Body:       strings.NewReader(string(jsonData)),
			Refresh:    "true",
		}

		res, err := req.Do(context.Background(), es)
		if err != nil {
			log.Printf("Error indexing document %d: %s", animeID, err)
			continue
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Printf("Error response from Elasticsearch for document %d: %s", animeID, res.String())
		} else {
			log.Printf("Successfully indexed document %d", animeID)
		}

	}

	fmt.Println("Data successfully inserted into Elasticsearch!")
	time.Sleep(500 * time.Millisecond)
}
