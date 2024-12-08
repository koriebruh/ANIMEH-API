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
	db := conf.InitDB()
	//insert.InsertDataCSVToElastic(es)

	animeService := service.NewAnimeService(es)
	userService := service.NewUserService(es, db)
	// Setup Gin router
	r := gin.Default()

	// Define route for auto-complete
	r.GET("/autocomplete", animeService.AutoComplete)
	r.GET("/search/anime", animeService.SearchAnime)
	r.GET("/anime/top", animeService.TopAnime)
	r.GET("/anime/:id", animeService.FindById)
	r.GET("/anime/:id/recommend", animeService.RecommendById)

	r.POST("/users", userService.Register)
	r.POST("/users/login", userService.Login)
	r.POST("/users/fav/:id", conf.JWTAuthMiddleware(), userService.AddFavAnime)
	r.DELETE("/users/fav/:id", conf.JWTAuthMiddleware(), userService.RemoveFavAnime)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "application/json", "WOI")
	})

	// Start server
	slog.Info("RUN IN PORT 8080")
	r.Run(":8080") // Listen on port 8080
}
