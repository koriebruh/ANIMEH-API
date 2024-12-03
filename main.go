package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"koriebruh/find/conf"
	"log/slog"
	"net/http"
	"strings"
)

// Elasticsearch client
var es = conf.ElasticClient()

// Auto-complete handler
func autoCompleteHandler(c *gin.Context) {
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
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("anime_info"),
		es.Search.WithBody(strings.NewReader(esQuery)),
		es.Search.WithPretty(),
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

func main() {

	// Kemudian sisipkan data CSV
	//insert.InsertDataCSVToElastic()

	// Setup Gin router
	r := gin.Default()

	// Define route for auto-complete
	r.GET("/autocomplete", autoCompleteHandler)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "application/json", "WOI")
	})

	// Start server
	slog.Info("RUN IN PORT 8080")
	r.Run(":8080") // Listen on port 8080
}
