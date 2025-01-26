package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateExchangeRate(ctx *gin.Context) {
	var exchangeRate models.ExchangeRate

	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	exchangeRate.Date = time.Now()

	if err := global.DB.AutoMigrate(&exchangeRate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to migrate database",
		})
		return
	}

	if err := global.DB.Create(&exchangeRate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create exchange rate",
		})
		return
	}

	ctx.JSON(http.StatusOK, exchangeRate)
}

func GetExchangeRate(ctx *gin.Context) {
	var exchangeRate []models.ExchangeRate

	if err := global.DB.Find(&exchangeRate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Exchange rate not found",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve exchange rate",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, exchangeRate)
}
