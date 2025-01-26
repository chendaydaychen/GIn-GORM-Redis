package router

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://....com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Register)
	}

	api := r.Group("/api")
	api.GET("/exchangerates", controllers.GetExchangeRate)
	api.Use(middlewares.AuthMiddleware())
	{
		api.POST("/exchangerates", controllers.CreateExchangeRate)
		api.POST("/articles", controllers.CreatArticle)
		api.GET("/articles", controllers.GetArticle)
		api.GET("/articles/:id", controllers.GetArticleById)
		api.POST("/articles/:id/like", controllers.LikeArticle)
		api.GET("/articles/:id/like", controllers.GetArticleLikes)

	}

	return r
}
