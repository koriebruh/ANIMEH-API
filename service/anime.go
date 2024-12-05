package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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

func (s *AnimeServiceImpl) SearchAnime(c *gin.Context) {
	// Ambil parameter dari query
	name := c.Query("name")
	from := c.DefaultQuery("from", "0")
	size := c.DefaultQuery("size", "20")

	// Buat query berdasarkan spesifikasi yang diberikan
	query := map[string]interface{}{
		"from": from,
		"size": size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"name": map[string]interface{}{
								"query":     name,
								"fuzziness": "AUTO",
							},
						},
					},
				},
				"should": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"genres": "Comedy",
						},
					},
					map[string]interface{}{
						"bool": map[string]interface{}{
							"must_not": map[string]interface{}{
								"exists": map[string]interface{}{
									"field": "genres",
								},
							},
						},
					},
					map[string]interface{}{
						"range": map[string]interface{}{
							"score": map[string]interface{}{
								"gte": 8, // Nilai score lebih besar atau sama dengan 8.0
							},
						},
					},
					map[string]interface{}{
						"bool": map[string]interface{}{
							"must_not": map[string]interface{}{
								"exists": map[string]interface{}{
									"field": "score",
								},
							},
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"status": "Finished Airing", // Memfilter berdasarkan status
						},
					},
				},
			},
		},
	}

	// Serialize query ke JSON
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error encoding query: %v", err)})
		return
	}

	// Kirim permintaan ke Elasticsearch
	res, err := s.Client.Search(
		s.Client.Search.WithContext(context.Background()),
		s.Client.Search.WithIndex("anime_info"),
		s.Client.Search.WithBody(&buf),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error performing search: %v", err)})
		return
	}
	defer res.Body.Close()

	// Periksa status Elasticsearch response
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Elasticsearch error: %v", string(body))})
		return
	}

	// Decode hasil Elasticsearch
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error decoding response: %v", err)})
		return
	}

	// Kirim hasil ke client
	c.JSON(http.StatusOK, result)
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
