package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/handlers"
	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/services"
)

func Setup(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	jwtSvc *services.JWTService,
	dbConnHandler *handlers.DBConnectionHandler,
	resourceDBHandler *handlers.ResourceDBConfigHandler,
	resourceQueryHandler *handlers.ResourceQueryHandler,
) {
	r.Use(middleware.CORS())

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(jwtSvc))
		{
			auth.GET("/auth/me", authHandler.Me)

			dbConns := auth.Group("/db-connections")
			{
				dbConns.GET("", dbConnHandler.List)
				dbConns.GET("/:id", dbConnHandler.GetByID)
				dbConns.POST("", dbConnHandler.Create)
				dbConns.PUT("/:id", dbConnHandler.Update)
				dbConns.DELETE("/:id", dbConnHandler.Delete)
				dbConns.POST("/test", dbConnHandler.TestConnection)
				dbConns.POST("/:id/test", dbConnHandler.TestConnection)
			}

			resConns := auth.Group("/resource-db-configs")
			{
				resConns.GET("", resourceDBHandler.List)
				resConns.GET("/:id", resourceDBHandler.GetByID)
				resConns.POST("", resourceDBHandler.Create)
				resConns.PUT("/:id", resourceDBHandler.Update)
				resConns.DELETE("/:id", resourceDBHandler.Delete)
				resConns.POST("/:id/test", resourceDBHandler.TestConnection)
			}

			resources := auth.Group("/resources")
			{
				resources.GET("/:resource", resourceQueryHandler.Query)
				resources.GET("/:resource/:id", resourceQueryHandler.GetByID)
				resources.GET("/:resource/columns", resourceQueryHandler.GetColumns)
			}
		}
	}
}
