package insert

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"strings"
)

func AnimeIndex(es *elasticsearch.Client) error {
	// Definisi mapping untuk indeks
	mapping := `{
        "mappings": {
            "properties": {
                "anime_id": {"type": "integer"},
                "name": {"type": "text"},
                "english_name": {"type": "text"},
                "other_name": {"type": "text"},
                "score": {"type": "float"},
                "genres": {"type": "text"},
                "synopsis": {"type": "text"},
                "type": {"type": "keyword"},
                "episodes": {"type": "text"},
                "aired": {"type": "text"},
                "premiered": {"type": "text"},
                "status": {"type": "keyword"},
                "producers": {"type": "keyword"},
                "licensors": {"type": "keyword"},
                "studios": {"type": "keyword"},
                "source": {"type": "keyword"},
                "duration": {"type": "text"},
                "rating": {"type": "keyword"},
                "rank": {"type": "float"},
                "popularity": {"type": "integer"},
                "favorites": {"type": "integer"},
                "scored_by": {"type": "text"},
                "members": {"type": "integer"},
                "image_url": {"type": "keyword"},
                "embedding": {
                    "type": "dense_vector", 
                    "dims": 4
                }
            }
        }
    }`

	// Hapus index jika sudah ada
	es.Indices.Delete([]string{"anime_info"}, es.Indices.Delete.WithIgnoreUnavailable(true))

	// Buat index baru dengan mapping
	res, err := es.Indices.Create(
		"anime_info",
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return fmt.Errorf("gagal membuat index: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error response saat membuat index: %s", res.String())
	}

	return nil
}
