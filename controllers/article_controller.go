package controllers

import (
	"encoding/json"
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var cachekey = "articles"

func CreatArticle(ctx *gin.Context) {
	var article models.Article

	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := global.DB.AutoMigrate(&article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": " Failed to migrate database",
		})
		return
	}

	if err := global.DB.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": " Failed to create article",
		})
		return
	}

	if err := global.RedisDB.Del(cachekey).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete cache",
		})
	}

	ctx.JSON(http.StatusOK, article)
}

func GetArticle(ctx *gin.Context) {

	cacheData, err := global.RedisDB.Get(cachekey).Result()

	if err == redis.Nil {

		var articles []models.Article

		if err := global.DB.Find(&articles).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Article not found",
			})
			return
		}

		articleJSON, err := json.Marshal(articles)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to marshal articles",
			})
			return
		}

		if err := global.RedisDB.Set(cachekey, articleJSON, 10*time.Minute).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to set cache",
			})
		}

		ctx.JSON(http.StatusOK, articles)

	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get cache",
		})
		return
	} else {
		var articles []models.Article
		if err := json.Unmarshal([]byte(cacheData), &articles); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to unmarshal cache data",
			})
			return
		}
		ctx.JSON(http.StatusOK, articles)
	}

}

func GetArticleById(ctx *gin.Context) {
	id := ctx.Param("id")
	var article models.Article
	if err := global.DB.Where("id = ?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Article not found",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Article not found",
			})
		}
		return
	}
	ctx.JSON(http.StatusOK, article)
}
