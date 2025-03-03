package controllers

import (
	"exchangeapp/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// LikeArticle 增加文章的点赞数
// 参数:
//
//	ctx *gin.Context: Gin框架的上下文对象，用于处理HTTP请求和响应
func LikeArticle(ctx *gin.Context) {
	// 获取URL参数中的文章ID
	articleID := ctx.Param("id")

	// 构造Redis中存储文章点赞数的键
	likekey := "article:" + articleID + ":likes"

	// 使用Redis的Incr命令增加点赞数，如果键不存在则创建并设置为1
	if err := global.RedisDB.Incr(likekey).Err(); err != nil {
		// 如果发生错误，返回HTTP 500错误响应
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to like article",
		})
	}

	// 点赞成功，返回HTTP 200成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Article liked successfully",
	})
}

// GetArticleLikes 获取文章的点赞数
// 参数:
//
//	ctx *gin.Context: Gin框架的上下文对象，用于处理HTTP请求和响应
func GetArticleLikes(ctx *gin.Context) {
	// 获取URL参数中的文章ID
	articleID := ctx.Param("id")
	// 构造Redis中存储文章点赞数的键
	likekey := "article:" + articleID + ":likes"
	// 从Redis获取点赞数
	likes, err := global.RedisDB.Get(likekey).Result()
	// 如果键不存在，返回默认值0
	if err == redis.Nil {
		likes = "0"
	} else if err != nil {
		// 如果发生其他错误，返回HTTP 500错误响应
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get article likes",
		})
	}

	// 返回点赞数
	ctx.JSON(http.StatusOK, gin.H{
		"likes": likes,
	})
}
