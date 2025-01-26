package controllers

import (
	"exchangeapp/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func LikeArticle(ctx *gin.Context) {
	articleID := ctx.Param("id")

	likekey := "article:" + articleID + ":likes"

	if err := global.RedisDB.Incr(likekey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to like article",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Article liked successfully",
	})
}

func GetArticleLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")
	likekey := "article:" + articleID + ":likes"
	likes, err := global.RedisDB.Get(likekey).Result()
	if err == redis.Nil {
		likes = "0"
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get article likes",
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"likes": likes,
	})
}
