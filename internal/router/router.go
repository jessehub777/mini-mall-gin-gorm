package router

import (
	"strings"

	"mini-mall-gin-gorm/internal/config"
	"mini-mall-gin-gorm/internal/handler"
	"mini-mall-gin-gorm/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 注册全量路由。
func SetupRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler,
	purchaseHandler *handler.PurchaseHandler,
) *gin.Engine {
	r := gin.Default()

	apiPrefix := normalizeAPIPrefix(cfg.App.APIPrefix)
	api := r.Group(apiPrefix)
	{
		// 健康检查接口
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// 商品查询接口可匿名访问。
		api.GET("/products", productHandler.List)
		api.GET("/products/:id", productHandler.GetByID)

		// 需要登录的接口。
		private := api.Group("")
		private.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			private.GET("/users/me", userHandler.GetMe)
			private.PUT("/users/me", userHandler.UpdateMe)

			private.POST("/products", productHandler.Create)
			private.PUT("/products/:id", productHandler.Update)
			private.DELETE("/products/:id", productHandler.Delete)

			private.POST("/purchases", purchaseHandler.Create)
			private.GET("/purchases/me", purchaseHandler.ListMine)
		}
	}

	return r
}

// normalizeAPIPrefix 规范化 API 前缀，允许以下形式：
// 1) "/api" -> "/api"
// 2) "api" -> "/api"
// 3) "" 或 "/" -> ""（表示无前缀）
func normalizeAPIPrefix(prefix string) string {
	p := strings.TrimSpace(prefix)
	if p == "" || p == "/" {
		return ""
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return strings.TrimRight(p, "/")
}
