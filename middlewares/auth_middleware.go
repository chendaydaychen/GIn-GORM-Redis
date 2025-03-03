package middlewares

import (
	"exchangeapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 返回一个 Gin 中间件函数，用于验证请求中的 JWT 令牌。
//
// 该中间件执行以下步骤:
// 1. 从请求头中获取 "Authorization" 字段的值。
// 2. 如果令牌为空，返回 401 Unauthorized 状态码并中止请求。
// 3. 使用 utils.ParseJWT 函数解析令牌以获取用户名。
// 4. 如果解析失败，返回 401 Unauthorized 状态码并中止请求。
// 5. 如果解析成功，将用户名设置到 Gin 上下文中，并继续处理请求。
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "token Unauthorized",
			})
			ctx.Abort()
			return
		}

		username, err := utils.ParseJWT(token)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		ctx.Set("username", username)
		ctx.Next()
	}
}
