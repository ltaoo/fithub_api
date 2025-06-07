package main

import (
	"log"
	"myapi/config"
	"myapi/internal/api/routes"
	"myapi/internal/db"
	"myapi/pkg/logger"
)

// var trans ut.Translator // 全局翻译器

func main() {
	// 初始化配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger := logger.NewLogger(cfg.LogLevel)
	defer logger.Sync()

	// 初始化数据库连接
	database, err := db.NewDatabase(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}

	// 运行数据库迁移
	migrator := db.NewMigrator(cfg, logger)
	if err := migrator.MigrateUp(); err != nil {
		logger.Fatal("Failed to run migrations", err)
	}

	// 设置路由
	r := routes.SetupRouter(database, logger, cfg)

	// 启动服务器
	logger.Info("Starting server on " + cfg.ServerAddress)
	if err := r.Run(cfg.ServerAddress); err != nil {
		logger.Fatal("Failed to start server", err)
	}
}
