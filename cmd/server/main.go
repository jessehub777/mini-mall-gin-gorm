package main

import (
	"fmt"
	"log"

	"mini-mall-gin-gorm/internal/config"
	"mini-mall-gin-gorm/internal/handler"
	"mini-mall-gin-gorm/internal/model"
	"mini-mall-gin-gorm/internal/repository"
	"mini-mall-gin-gorm/internal/router"
	"mini-mall-gin-gorm/internal/service"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 优先加载项目根目录 .env，若不存在则继续使用系统环境变量。
	_ = godotenv.Load()

	// 1. 加载配置
	cfg := config.Load()

	// 2. 初始化数据库连接
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 3. 自动迁移表结构（学习场景推荐保留）
	if err := db.AutoMigrate(&model.User{}, &model.Product{}, &model.Purchase{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// 4. 初始化仓储层
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)

	// 5. 初始化服务层
	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	purchaseService := service.NewPurchaseService(db, productRepo, purchaseRepo)

	// 6. 初始化处理器层
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	purchaseHandler := handler.NewPurchaseHandler(purchaseService)

	// 7. 注册路由并启动服务
	r := router.SetupRouter(cfg, authHandler, userHandler, productHandler, purchaseHandler)
	addr := fmt.Sprintf(":%s", cfg.App.Port)
	log.Printf("%s 启动成功，监听地址: %s", cfg.App.Name, addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(cfg.DSN()), gormCfg)
	if err != nil {
		return nil, err
	}
	return db, nil
}
