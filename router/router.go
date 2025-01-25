package router

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {

	r := gin.Default()

	auth := r.Group("/api/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			c.AbortWithStatusJSON(200, gin.H{
				"message": "login",
			})
		})
		auth.POST("/register", func(c *gin.Context) {
			c.AbortWithStatusJSON(200, gin.H{
				"message": "register",
			})
		})
	}
	return r
}
