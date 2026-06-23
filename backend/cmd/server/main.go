package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/anraaa/visual-mesin/internal/config"
	"github.com/anraaa/visual-mesin/internal/db"
	"github.com/anraaa/visual-mesin/internal/handlers"
	"github.com/anraaa/visual-mesin/internal/repository"
	"github.com/anraaa/visual-mesin/internal/routes"
	"github.com/anraaa/visual-mesin/internal/services"

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

	dbMgr := db.NewManager()

	jwtExpiry, _ := time.ParseDuration(cfg.JWTExpiry)
	jwtSvc := services.NewJWTService(cfg.JWTSecret, jwtExpiry)

	// Repositories
	userRepo := repository.NewUserRepository(gormDB)
	roleRepo := repository.NewRoleRepository(gormDB)
	dbConnRepo := repository.NewDBConnectionRepository(gormDB)
	resourceDBConfigRepo := repository.NewResourceDBConfigRepository(gormDB)

	// Services
	authSvc := services.NewAuthService(userRepo, jwtSvc)
	dbConnSvc := services.NewDBConnectionService(dbConnRepo, dbMgr, gormDB)
	resourceDBConfigSvc := services.NewResourceDBConfigService(resourceDBConfigRepo, dbMgr, gormDB)
	resourceQuerySvc := services.NewResourceQueryService(resourceDBConfigRepo, dbMgr)

	// Handlers
	authHandler := handlers.NewAuthHandler(authSvc)
	dbConnHandler := handlers.NewDBConnectionHandler(dbConnSvc)
	resourceDBHandler := handlers.NewResourceDBConfigHandler(resourceDBConfigSvc)
	resourceQueryHandler := handlers.NewResourceQueryHandler(resourceQuerySvc)
	moduleHandler := handlers.NewModuleHandler(resourceQuerySvc)

	ginMode := gin.DebugMode
	if cfg.AppEnv == "production" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.Setup(r, authHandler, jwtSvc, dbConnHandler, resourceDBHandler, resourceQueryHandler, moduleHandler)
	_ = roleRepo

	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	r.Run(addr)
}
