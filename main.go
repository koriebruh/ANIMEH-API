package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"koriebruh/find/conf"
	"koriebruh/find/service"
	"log/slog"
	"net/http"
	"time"
)

func main() {

	cnf := conf.GetConfig()
	es := conf.ElasticClient(cnf)
	db := conf.InitDB(cnf)

	animeService := service.NewAnimeService(es)
	userService := service.NewUserService(es, db)
	// Setup Gin router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Define route for auto-complete
	r.GET("/autocomplete", animeService.AutoComplete)
	r.GET("/search/anime", animeService.SearchAnime)
	r.GET("/anime/top", animeService.TopAnime)
	r.GET("/anime/:id", animeService.FindById)
	r.GET("/anime/:id/recommend", animeService.RecommendById)

	r.POST("/users", userService.Register)
	r.POST("/users/login", userService.Login)
	r.POST("/users/change", conf.JWTAuthMiddleware(), userService.ChangePass)
	r.POST("/users/change-confirm", conf.JWTAuthMiddleware(), userService.ConfirmChangePass)

	r.POST("/users/fav/:id", conf.JWTAuthMiddleware(), userService.AddFavAnime)
	r.DELETE("/users/fav/:id", conf.JWTAuthMiddleware(), userService.RemoveFavAnime)
	r.GET("/users/fav", conf.JWTAuthMiddleware(), userService.FindAllFavAnime)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "application/json", "WOI")
	})

	serverRun := fmt.Sprintf("%s:%s", cnf.Server.Host, cnf.Server.Port)
	slog.Info(serverRun)
	r.Run(serverRun)
}
