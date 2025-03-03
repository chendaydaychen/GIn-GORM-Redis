package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register 处理用户注册请求。
//
// 参数:
// - ctx: Gin 上下文，用于绑定请求体和发送 JSON 响应。
//
// 该函数执行以下步骤:
// 1. 将传入的 JSON 数据绑定到 User 结构体。
// 2. 使用 utils.HashPassword 函数对提供的密码进行哈希处理。
// 3. 为新注册的用户生成 JWT 令牌。
// 4. 使用 AutoMigrate 确保数据库模式是最新的。
// 5. 将新用户创建到数据库中。
// 6. 注册成功后返回生成的令牌，如果任何步骤失败则返回错误信息。
func Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求体",
		})
		return
	}

	hashedpassword, err := utils.HashPassword(user.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "密码哈希失败",
		})
		return
	}

	user.Password = hashedpassword

	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成令牌失败",
		})
		return
	}

	if err := global.DB.AutoMigrate(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "数据库迁移失败",
		})
		return
	}

	if err := global.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建用户失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// Login 处理用户登录请求。
//
// 参数:
// - ctx: Gin 上下文，用于绑定请求体和发送 JSON 响应。
//
// 该函数执行以下步骤:
// 1. 将传入的 JSON 数据绑定到包含用户名和密码的结构体。
// 2. 根据用户名从数据库中查找用户。
// 3. 使用 utils.CheckPassword 函数验证提供的密码是否正确。
// 4. 为登录成功的用户生成 JWT 令牌。
// 5. 登录成功后返回生成的令牌，如果任何步骤失败则返回错误信息。
func Login(ctx *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求体",
		})
		return
	}

	var user models.User

	if err := global.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户名",
		})
		return
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的密码",
		})
		return
	}

	token, err := utils.GenerateJWT(user.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成令牌失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
