package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

	ctx.JSON(http.StatusOK, article)
}

func GetArticle(ctx *gin.Context) {
	var articles []models.Article

	if err := global.DB.Find(&articles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Article not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, articles)
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
