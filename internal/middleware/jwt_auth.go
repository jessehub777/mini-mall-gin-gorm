package middleware

import (
	"net/http"
	"strings"

	"mini-mall-gin-gorm/pkg/jwtutil"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 校验请求头中的 JWT，并将用户信息写入上下文。
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": "未携带 Authorization"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": "Authorization 格式错误，应为 Bearer <token>"})
			c.Abort()
			return
		}

		claims, err := jwtutil.ParseToken(secret, parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": "Token 无效或已过期"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// GetCurrentUserID 从上下文中提取当前登录用户 ID。
func GetCurrentUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}
