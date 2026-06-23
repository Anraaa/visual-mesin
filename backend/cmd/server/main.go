package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/anraaa/visual-mesin/internal/ai"
	"github.com/anraaa/visual-mesin/internal/config"
	"github.com/anraaa/visual-mesin/internal/db"
	"github.com/anraaa/visual-mesin/internal/handlers"
	"github.com/anraaa/visual-mesin/internal/repository"
	"github.com/anraaa/visual-mesin/internal/routes"
	"github.com/anraaa/visual-mesin/internal/services"
	"github.com/anraaa/visual-mesin/internal/ws"

	_ "github.com/anraaa/visual-mesin/docs"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	zerolog.New(os.Stderr).With().Timestamp().Logger()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	if err := db.RunMigrations(dsn, "migrations"); err != nil {
		log.Printf("Migration warning: %v", err)
	}

	db.SeedDefaultUsers(gormDB)

	dbMgr := db.NewManager()

	redisAddr := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Printf("Redis not available: %v", err)
	}

	jwtExpiry, _ := time.ParseDuration(cfg.JWTExpiry)
	jwtSvc := services.NewJWTService(cfg.JWTSecret, jwtExpiry)

	userRepo := repository.NewUserRepository(gormDB)
	roleRepo := repository.NewRoleRepository(gormDB)
	dbConnRepo := repository.NewDBConnectionRepository(gormDB)
	resourceDBConfigRepo := repository.NewResourceDBConfigRepository(gormDB)
	exportJobRepo := repository.NewExportJobRepository(gormDB)

	authSvc := services.NewAuthService(userRepo, jwtSvc)
	dbConnSvc := services.NewDBConnectionService(dbConnRepo, dbMgr, gormDB)
	resourceDBConfigSvc := services.NewResourceDBConfigService(resourceDBConfigRepo, dbMgr, gormDB)
	resourceQuerySvc := services.NewResourceQueryService(resourceDBConfigRepo, dbMgr)
	activityLogSvc := services.NewActivityLogService(gormDB)

	// WebSocket Hub
	wsHub := ws.NewHub()
	go wsHub.Run()

	exportSvc := services.NewExportService(exportJobRepo, resourceQuerySvc, rdb, dbMgr, cfg.ExportDir, wsHub, activityLogSvc)
	exportSvc.StartWorker()

	// AI Chat
	aiSchemaMapRepo := repository.NewAiSchemaMapRepository(gormDB)
	aiChatHistoryRepo := repository.NewAiChatHistoryRepository(gormDB)

	aiSchemaMapSvc := services.NewAiSchemaMapService(aiSchemaMapRepo)
	aiChatHistorySvc := services.NewAiChatHistoryService(aiChatHistoryRepo)

	ollamaClient := ai.NewOllamaClient(cfg.OllamaURL, cfg.OllamaModel)
	intentDetector := ai.NewIntentDetector(ollamaClient)
	sqlGenerator := ai.NewSQLGenerator(ollamaClient)
	sqlFirewall := ai.NewSQLFirewall()
	sqlExecutor := ai.NewSQLExecutor(func(driver, host string, port int, username, password, dbName string) (*sql.DB, error) {
		edsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			username, password, host, port, dbName)
		return sql.Open("mysql", edsn)
	})

	chatPipeline := ai.NewChatPipeline(intentDetector, sqlGenerator, sqlFirewall, sqlExecutor, ollamaClient, aiSchemaMapRepo, aiChatHistoryRepo)

	aiSchemaMapHandler := handlers.NewAiSchemaMapHandler(aiSchemaMapSvc)
	aiChatHandler := handlers.NewAiChatHandler(chatPipeline, aiChatHistorySvc)

	authHandler := handlers.NewAuthHandler(authSvc)
	userHandler := handlers.NewUserHandler(userRepo, roleRepo)
	roleHandler := handlers.NewRoleHandler(roleRepo)
	dbConnHandler := handlers.NewDBConnectionHandler(dbConnSvc)
	resourceDBHandler := handlers.NewResourceDBConfigHandler(resourceDBConfigSvc)
	resourceQueryHandler := handlers.NewResourceQueryHandler(resourceQuerySvc)
	exportHandler := handlers.NewExportHandler(exportSvc)
	activityLogHandler := handlers.NewActivityLogHandler(activityLogSvc)

	ginMode := gin.DebugMode
	if cfg.AppEnv == "production" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Setup(r, authHandler, jwtSvc, dbConnHandler, resourceDBHandler, resourceQueryHandler,
		userHandler, roleHandler, exportHandler, aiSchemaMapHandler, aiChatHandler,
		wsHub, activityLogHandler, cfg.JWTSecret)

	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	r.Run(addr)
}
