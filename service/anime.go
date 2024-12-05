package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AnimeService interface {
	AutoComplete(c *gin.Context)
	SearchAnime(c *gin.Context)
	recommendations(c *gin.Context)
	FindById(c *gin.Context)
	AnimeFilter(c *gin.Context)
	TopAnime(c *gin.Context)
}

type AnimeServiceImpl struct {
	*elasticsearch.Client
}

func NewAnimeService(client *elasticsearch.Client) *AnimeServiceImpl {
	return &AnimeServiceImpl{Client: client}
}

func (s AnimeServiceImpl) AutoComplete(c *gin.Context) {
	// Get the query parameter
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	// Elasticsearch query
	esQuery := fmt.Sprintf(`{
		"query": {
			"match_phrase_prefix": {
				"name": "%s"
			}
		}
	}`, query)

	// Execute search
	res, err := s.Client.Search(
		s.Client.Search.WithContext(context.Background()),
		s.Client.Search.WithIndex("anime_info"),
		s.Client.Search.WithBody(strings.NewReader(esQuery)),
		s.Client.Search.WithPretty(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error searching documents"})
		return
	}
	defer res.Body.Close()

	// Parse response
	if res.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{"error": res.String()})
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing response"})
		return
	}

	// Extract hits
	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	suggestions := []string{}
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"]
		name := source.(map[string]interface{})["name"].(string)
		suggestions = append(suggestions, name)
	}

	// Return suggestions
	c.JSON(http.StatusOK, gin.H{"suggestions": suggestions})
}

func (s AnimeServiceImpl) SearchAnime(c *gin.Context) {
	// Ambil parameter query dan lainnya
	query := c.DefaultQuery("query", "")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")
	genre := c.DefaultQuery("genre", "")
	status := c.DefaultQuery("status", "")
	minScore := c.DefaultQuery("min_score", "")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)
	from := (pageNum - 1) * limitNum

	// Membuat query pencarian secara langsung dalam string
	searchQuery := fmt.Sprintf(`
		{
			"from": %d,
			"size": %d,
			"query": {
				"bool": {
					"must": [
						{
							"match": {
								"name": {
									"query": "%s",
									"fuzziness": "AUTO"
								}
							}
						}
					],
					"filter": [%s]
				}
			}
		}
	`, from, limitNum, query, generateFilter(genre, status, minScore))

	// Inisialisasi client ElasticSearch
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // Ganti dengan alamat server Elasticsearch kamu
		},
	})
	if err != nil {
		log.Printf("Error creating ElasticSearch client: %s", err)
		c.JSON(500, gin.H{"error": "Failed to create ElasticSearch client"})
		return
	}

	// Mengonversi string query menjadi io.Reader
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("anime_info"),                  // Nama index anime di ElasticSearch
		es.Search.WithBody(strings.NewReader(searchQuery)), // Menggunakan strings.NewReader
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Printf("Error executing search: %s", err)
		c.JSON(500, gin.H{"error": "Failed to execute search"})
		return
	}
	defer res.Body.Close()

	// Jika berhasil, proses hasilnya dan kirimkan ke client
	var searchResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		log.Printf("Error decoding response: %s", err)
		c.JSON(500, gin.H{"error": "Failed to decode search response"})
		return
	}

	// Kirimkan hasil pencarian dalam format JSON
	c.JSON(200, searchResult)
}

func generateFilter(genre, status, minScore string) string {
	var filters []string

	if genre != "" {
		filters = append(filters, fmt.Sprintf(`{ "term": { "genres.keyword": "%s" } }`, genre))
	}
	if status != "" {
		filters = append(filters, fmt.Sprintf(`{ "term": { "status.keyword": "%s" } }`, status))
	}
	if minScore != "" {
		filters = append(filters, fmt.Sprintf(`{ "range": { "score": { "gte": %s } } }`, minScore))
	}

	// Jika tidak ada filter yang ditambahkan, kembalikan array kosong
	if len(filters) > 0 {
		return strings.Join(filters, ",")
	}
	return ""
}

func (s AnimeServiceImpl) recommendations(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s AnimeServiceImpl) FindById(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s AnimeServiceImpl) AnimeFilter(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (s AnimeServiceImpl) TopAnime(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
