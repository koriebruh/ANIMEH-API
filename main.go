package main

import (
	"github.com/gin-gonic/gin"
	"koriebruh/find/conf"
	"koriebruh/find/service"
	"log/slog"
	"net/http"
)

// Elasticsearch client

func main() {
	// Kemudian sisipkan data CSV
	es := conf.ElasticClient()
	//insert.InsertDataCSVToElastic(es)

	animeService := service.NewAnimeService(es)
	// Setup Gin router
	r := gin.Default()

	// Define route for auto-complete
	//r.GET("/autocomplete", handler.AutoCompleteHandler)
	r.GET("/autocomplete", animeService.AutoComplete)
	r.GET("/search/anime", animeService.SearchAnime)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "application/json", "WOI")
	})

	// Start server
	slog.Info("RUN IN PORT 8080")
	r.Run(":8080") // Listen on port 8080
}
