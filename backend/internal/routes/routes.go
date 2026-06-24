package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/anraaa/visual-mesin/internal/handlers"
	"github.com/anraaa/visual-mesin/internal/middleware"
	"github.com/anraaa/visual-mesin/internal/services"
	"github.com/anraaa/visual-mesin/internal/ws"
)

func Setup(
	r *gin.Engine,
	authHandler *handlers.AuthHandler,
	jwtSvc *services.JWTService,
	dbConnHandler *handlers.DBConnectionHandler,
	resourceDBHandler *handlers.ResourceDBConfigHandler,
	resourceQueryHandler *handlers.ResourceQueryHandler,
	userHandler *handlers.UserHandler,
	roleHandler *handlers.RoleHandler,
	permHandler *handlers.PermissionHandler,
	exportHandler *handlers.ExportHandler,
	aiSchemaMapHandler *handlers.AiSchemaMapHandler,
	aiChatHandler *handlers.AiChatHandler,
	wsHub *ws.Hub,
	activityLogHandler *handlers.ActivityLogHandler,
	dataProduksiConfigHandler *handlers.DataProduksiConfigHandler,
	dashboardHandler *handlers.DashboardHandler,
	jwtSecret string,
) {
	r.Use(middleware.CORS())

	r.Static("/exports", "./exports")

	r.GET("/ws", func(c *gin.Context) {
		ws.HandleWebSocket(wsHub, jwtSecret)(c.Writer, c.Request)
	})

	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", authHandler.Login)
		api.POST("/auth/register", authHandler.Register)

		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(jwtSvc))
		{
			auth.GET("/auth/me", authHandler.Me)

			users := auth.Group("/users")
			{
				users.GET("", userHandler.List)
				users.GET("/:id", userHandler.GetByID)
				users.POST("", userHandler.Create)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", userHandler.Delete)
				users.POST("/:id/assign-role", userHandler.AssignRole)
				users.PUT("/:id/sync-roles", userHandler.SyncRoles)
			}

			roles := auth.Group("/roles")
			{
				roles.GET("", roleHandler.List)
				roles.GET("/:id", roleHandler.GetByID)
				roles.POST("", roleHandler.Create)
				roles.DELETE("/:id", roleHandler.Delete)
				roles.POST("/:id/assign-permission", roleHandler.AssignPermission)
				roles.DELETE("/:id/revoke-permission/:permId", roleHandler.RevokePermission)
				roles.POST("/:id/sync-permissions", roleHandler.SyncPermissions)
			}

			permissions := auth.Group("/permissions")
			{
				permissions.GET("", permHandler.List)
				permissions.POST("", permHandler.Create)
				permissions.DELETE("/:id", permHandler.Delete)
			}

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
				resources.GET("/:resource", resourceQueryHandler.List)
				resources.GET("/:resource/stats", resourceQueryHandler.GetStats)
				resources.GET("/:resource/trend", resourceQueryHandler.GetTrend)
				resources.GET("/:resource/spc", resourceQueryHandler.GetSPC)
				resources.GET("/:resource/distribution", resourceQueryHandler.GetDistribution)
				resources.GET("/:resource/quality-trend", resourceQueryHandler.GetQualityTrend)
				resources.GET("/:resource/columns", resourceQueryHandler.GetColumns)
				resources.GET("/:resource/:id", resourceQueryHandler.GetByID)
				resources.POST("/:resource", resourceQueryHandler.Create)
				resources.PUT("/:resource/:id", resourceQueryHandler.Update)
				resources.DELETE("/:resource/:id", resourceQueryHandler.Delete)
			}

			auth.GET("/resources-judgment/:resource", resourceQueryHandler.GetJudgmentSummary)

			dpConfig := auth.Group("/data-produksi-config")
			{
				dpConfig.GET("/groups", dataProduksiConfigHandler.ListGroups)
				dpConfig.POST("/groups", dataProduksiConfigHandler.CreateGroup)
				dpConfig.PUT("/groups/:id", dataProduksiConfigHandler.UpdateGroup)
				dpConfig.DELETE("/groups/:id", dataProduksiConfigHandler.DeleteGroup)

				dpConfig.POST("/items", dataProduksiConfigHandler.CreateItem)
				dpConfig.PUT("/items/:itemId", dataProduksiConfigHandler.UpdateItem)
				dpConfig.DELETE("/items/:itemId", dataProduksiConfigHandler.DeleteItem)

				dpConfig.POST("/create-table", dataProduksiConfigHandler.CreateResourceWithTable)
				dpConfig.GET("/columns/:resource", dataProduksiConfigHandler.GetColumnDefs)
			}

			exports := auth.Group("/exports")
			{
				exports.GET("", exportHandler.List)
				exports.POST("", exportHandler.Submit)
				exports.GET("/:id", exportHandler.GetStatus)
				exports.GET("/:id/download", exportHandler.Download)
			}

			dashboard := auth.Group("/dashboard")
			{
				dashboard.GET("/summary", dashboardHandler.GetSummary)
			}

			activityLog := auth.Group("/activity-logs")
			{
				activityLog.GET("", activityLogHandler.List)
				activityLog.GET("/me", activityLogHandler.ListMy)
			}
		}
	}
}
