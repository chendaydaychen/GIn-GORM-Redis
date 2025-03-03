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

// cachekey 定义了缓存中存储文章列表的键名
var cachekey = "articles"

// CreatArticle 处理文章创建逻辑
// 参数: ctx *gin.Context 上下文，用于处理HTTP请求和响应
// 该函数首先尝试解析请求体以获取文章数据，然后迁移数据库模型，最后创建文章记录并更新缓存
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

// GetArticle 处理获取文章列表逻辑
// 参数: ctx *gin.Context 上下文，用于处理HTTP请求和响应
// 该函数首先尝试从缓存中获取文章列表，如果缓存未命中，则从数据库中查询并更新缓存
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

// GetArticleById 处理根据ID获取文章详情逻辑
// 参数: ctx *gin.Context 上下文，用于处理HTTP请求和响应
// 该函数根据提供的ID查询数据库中的文章记录，如果找到则返回文章详情，否则返回错误信息
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
