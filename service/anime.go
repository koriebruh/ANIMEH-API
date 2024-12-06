package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"log"
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
	// Ambil query string yang mungkin dikirim oleh client
	name := c.DefaultQuery("name", "")
	genres := c.DefaultQuery("genres", "")
	status := c.DefaultQuery("status", "")
	minScore := c.DefaultQuery("min_score", "8")

	// Build query JSON
	query := map[string]interface{}{
		"from": 0,
		"size": 20,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{
					map[string]interface{}{
						"match": map[string]interface{}{
							"name": map[string]interface{}{
								"query":     name,
								"fuzziness": "AUTO",
							},
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"genres": genres, // Filter berdasarkan genres
						},
					},
					map[string]interface{}{
						"range": map[string]interface{}{
							"score": map[string]interface{}{
								"gte": minScore, // Filter berdasarkan score minimal
							},
						},
					},
					map[string]interface{}{
						"term": map[string]interface{}{
							"status": status, // Filter berdasarkan status
						},
					},
				},
				"minimum_should_match": 1, // Setidaknya satu kondisi harus terpenuhi
			},
		},
	}

	// Ubah query menjadi JSON
	queryBody, err := json.Marshal(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to build query: %v", err),
		})
		return
	}

	// Melakukan pencarian ke Elasticsearch
	res, err := s.Client.Search(
		s.Client.Search.WithContext(context.Background()),
		s.Client.Search.WithIndex("anime_info"),
		s.Client.Search.WithBody(strings.NewReader(string(queryBody))), // Menggunakan strings.NewReader
		s.Client.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		log.Printf("Error searching Elasticsearch: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to execute search query",
		})
		return
	}
	defer res.Body.Close()

	// Periksa apakah ada error dalam response
	if res.IsError() {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Elasticsearch query error: %s", res.String()),
		})
		return
	}

	// Parsing response body
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to parse response: %v", err),
		})
		return
	}

	// Mengirimkan hasil pencarian
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
