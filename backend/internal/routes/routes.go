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
	genericHandler *handlers.ResourceQueryHandler,
	moduleHandler *handlers.ModuleHandler,
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

			// =========================================
			// Fase 2: DB Connection Management
			// =========================================
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

			// =========================================
			// Fase 3: Resource Table APIs (generic)
			// =========================================
			resources := auth.Group("/resources")
			{
				resources.GET("/:resource", genericHandler.List)
				resources.GET("/:resource/:id", genericHandler.GetByID)
				resources.POST("/:resource", genericHandler.Create)
				resources.PUT("/:resource/:id", genericHandler.Update)
				resources.DELETE("/:resource/:id", genericHandler.Delete)
				resources.GET("/:resource/columns", genericHandler.GetColumns)
			}

			// =========================================
			// Fase 3: Building Module (RTBA/RTBC/RTBE)
			// =========================================
			building := auth.Group("/building")
			{
				for _, t := range []string{"rtba1", "rtba2", "rtba3"} {
					tbl := t
					building.GET("/"+tbl, moduleHandler.ResourceList(tbl))
					building.GET("/"+tbl+"/:id", moduleHandler.ResourceGetByID(tbl, "recid"))
					building.POST("/"+tbl, moduleHandler.ResourceCreate(tbl))
					building.PUT("/"+tbl+"/:id", moduleHandler.ResourceUpdate(tbl, "recid"))
					building.DELETE("/"+tbl+"/:id", moduleHandler.ResourceDelete(tbl, "recid"))
					building.GET("/"+tbl+"/barcode/:barcode", moduleHandler.ResourceBarcodeLookup(tbl))
				}
				for _, t := range []string{"rtbc1", "rtbc2", "rtbc3", "rtbc4"} {
					tbl := t
					building.GET("/"+tbl, moduleHandler.ResourceList(tbl))
					building.GET("/"+tbl+"/:id", moduleHandler.ResourceGetByID(tbl, "recid"))
					building.POST("/"+tbl, moduleHandler.ResourceCreate(tbl))
					building.PUT("/"+tbl+"/:id", moduleHandler.ResourceUpdate(tbl, "recid"))
					building.DELETE("/"+tbl+"/:id", moduleHandler.ResourceDelete(tbl, "recid"))
					building.GET("/"+tbl+"/barcode/:barcode", moduleHandler.ResourceBarcodeLookup(tbl))
				}
				for _, t := range []string{"rtbe1", "rtbe2"} {
					tbl := t
					building.GET("/"+tbl, moduleHandler.ResourceList(tbl))
					building.GET("/"+tbl+"/:id", moduleHandler.ResourceGetByID(tbl, "recid"))
					building.POST("/"+tbl, moduleHandler.ResourceCreate(tbl))
					building.PUT("/"+tbl+"/:id", moduleHandler.ResourceUpdate(tbl, "recid"))
					building.DELETE("/"+tbl+"/:id", moduleHandler.ResourceDelete(tbl, "recid"))
					building.GET("/"+tbl+"/barcode/:barcode", moduleHandler.ResourceBarcodeLookup(tbl))
				}
			}

			// =========================================
			// Fase 3: Extruder Module
			// =========================================
			extruder := auth.Group("/extruder")
			{
				extruder.GET("/rteex1", moduleHandler.ResourceList("rteex1"))
				extruder.GET("/rteex1/:id", moduleHandler.ResourceGetByID("rteex1", "id"))
				extruder.POST("/rteex1", moduleHandler.ResourceCreate("rteex1"))
				extruder.PUT("/rteex1/:id", moduleHandler.ResourceUpdate("rteex1", "id"))
				extruder.DELETE("/rteex1/:id", moduleHandler.ResourceDelete("rteex1", "id"))

				extruder.GET("/rteex2", moduleHandler.ResourceList("rteex2"))
				extruder.GET("/rteex2/:id", moduleHandler.ResourceGetByID("rteex2", "recid"))
				extruder.POST("/rteex2", moduleHandler.ResourceCreate("rteex2"))
				extruder.PUT("/rteex2/:id", moduleHandler.ResourceUpdate("rteex2", "recid"))
				extruder.DELETE("/rteex2/:id", moduleHandler.ResourceDelete("rteex2", "recid"))

				extruder.GET("/rteex3head", moduleHandler.ResourceList("rteex3head"))
				extruder.GET("/rteex3head/:id", moduleHandler.ResourceGetByID("rteex3head", "recid"))
				extruder.POST("/rteex3head", moduleHandler.ResourceCreate("rteex3head"))
				extruder.PUT("/rteex3head/:id", moduleHandler.ResourceUpdate("rteex3head", "recid"))
				extruder.DELETE("/rteex3head/:id", moduleHandler.ResourceDelete("rteex3head", "recid"))

				extruder.GET("/recorddatacyclic", moduleHandler.ResourceList("recorddatacyclic"))
				extruder.GET("/recorddatacyclic/:id", moduleHandler.ResourceGetByID("recorddatacyclic", "id"))
				extruder.POST("/recorddatacyclic", moduleHandler.ResourceCreate("recorddatacyclic"))
				extruder.PUT("/recorddatacyclic/:id", moduleHandler.ResourceUpdate("recorddatacyclic", "id"))
				extruder.DELETE("/recorddatacyclic/:id", moduleHandler.ResourceDelete("recorddatacyclic", "id"))

				extruder.GET("/recorddatapcs", moduleHandler.ResourceList("recorddatapcs"))
				extruder.GET("/recorddatapcs/:id", moduleHandler.ResourceGetByID("recorddatapcs", "id"))
				extruder.POST("/recorddatapcs", moduleHandler.ResourceCreate("recorddatapcs"))
				extruder.PUT("/recorddatapcs/:id", moduleHandler.ResourceUpdate("recorddatapcs", "id"))
				extruder.DELETE("/recorddatapcs/:id", moduleHandler.ResourceDelete("recorddatapcs", "id"))

				extruder.GET("/datalog", moduleHandler.ResourceList("datalog"))
				extruder.GET("/datalog/:id", moduleHandler.ResourceGetByID("datalog", "id"))
				extruder.POST("/datalog", moduleHandler.ResourceCreate("datalog"))
				extruder.PUT("/datalog/:id", moduleHandler.ResourceUpdate("datalog", "id"))
				extruder.DELETE("/datalog/:id", moduleHandler.ResourceDelete("datalog", "id"))
			}

			// =========================================
			// Fase 3: Curing Module
			// =========================================
			curing := auth.Group("/curing")
			{
				curing.GET("/curtire", moduleHandler.ResourceList("curtire"))
				curing.GET("/curtire/:id", moduleHandler.ResourceGetByID("curtire", "recid"))
				curing.POST("/curtire", moduleHandler.ResourceCreate("curtire"))
				curing.PUT("/curtire/:id", moduleHandler.ResourceUpdate("curtire", "recid"))
				curing.DELETE("/curtire/:id", moduleHandler.ResourceDelete("curtire", "recid"))
				curing.GET("/curtire/barcode/:barcode", moduleHandler.ResourceBarcodeLookup("curtire"))

				curing.GET("/item_measurement", moduleHandler.ResourceList("item_measurement"))
				curing.GET("/item_measurement/:id", moduleHandler.ResourceGetByID("item_measurement", "recid"))
				curing.POST("/item_measurement", moduleHandler.ResourceCreate("item_measurement"))
				curing.PUT("/item_measurement/:id", moduleHandler.ResourceUpdate("item_measurement", "recid"))
				curing.DELETE("/item_measurement/:id", moduleHandler.ResourceDelete("item_measurement", "recid"))

				curing.GET("/gtentire", moduleHandler.ResourceList("gtentire"))
				curing.GET("/gtentire/:id", moduleHandler.ResourceGetByID("gtentire", "recid"))
				curing.POST("/gtentire", moduleHandler.ResourceCreate("gtentire"))
				curing.PUT("/gtentire/:id", moduleHandler.ResourceUpdate("gtentire", "recid"))
				curing.DELETE("/gtentire/:id", moduleHandler.ResourceDelete("gtentire", "recid"))
				curing.GET("/gtentire/barcode/:barcode", moduleHandler.ResourceBarcodeLookup("gtentire"))
			}

			// =========================================
			// Fase 3: Trimming Module
			// =========================================
			trimming := auth.Group("/trimming")
			{
				trimming.GET("/trimming", moduleHandler.ResourceList("trimming"))
				trimming.GET("/trimming/:id", moduleHandler.ResourceGetByID("trimming", "id"))
				trimming.POST("/trimming", moduleHandler.ResourceCreate("trimming"))
				trimming.PUT("/trimming/:id", moduleHandler.ResourceUpdate("trimming", "id"))
				trimming.DELETE("/trimming/:id", moduleHandler.ResourceDelete("trimming", "id"))

				trimming.GET("/rtc-tr1", moduleHandler.ResourceList("rtc-tr1"))
				trimming.GET("/rtc-tr1/:id", moduleHandler.ResourceGetByID("rtc-tr1", "recid"))
				trimming.POST("/rtc-tr1", moduleHandler.ResourceCreate("rtc-tr1"))
				trimming.PUT("/rtc-tr1/:id", moduleHandler.ResourceUpdate("rtc-tr1", "recid"))
				trimming.DELETE("/rtc-tr1/:id", moduleHandler.ResourceDelete("rtc-tr1", "recid"))
			}

			// =========================================
			// Fase 3: Monitoring Module
			// =========================================
			monitoring := auth.Group("/monitoring")
			{
				monitoring.GET("/monitoringtl1", moduleHandler.ResourceList("monitoringtl1"))
				monitoring.GET("/monitoringtl1/:id", moduleHandler.ResourceGetByID("monitoringtl1", "id"))
				monitoring.POST("/monitoringtl1", moduleHandler.ResourceCreate("monitoringtl1"))
				monitoring.PUT("/monitoringtl1/:id", moduleHandler.ResourceUpdate("monitoringtl1", "id"))
				monitoring.DELETE("/monitoringtl1/:id", moduleHandler.ResourceDelete("monitoringtl1", "id"))

				monitoring.GET("/rtl-tl1", moduleHandler.ResourceList("rtl-tl1"))
				monitoring.GET("/rtl-tl1/:id", moduleHandler.ResourceGetByID("rtl-tl1", "id"))
				monitoring.POST("/rtl-tl1", moduleHandler.ResourceCreate("rtl-tl1"))
				monitoring.PUT("/rtl-tl1/:id", moduleHandler.ResourceUpdate("rtl-tl1", "id"))
				monitoring.DELETE("/rtl-tl1/:id", moduleHandler.ResourceDelete("rtl-tl1", "id"))

				monitoring.GET("/rtltl1", moduleHandler.ResourceList("rtltl1"))
				monitoring.GET("/rtltl1/:id", moduleHandler.ResourceGetByID("rtltl1", "id"))
				monitoring.POST("/rtltl1", moduleHandler.ResourceCreate("rtltl1"))
				monitoring.PUT("/rtltl1/:id", moduleHandler.ResourceUpdate("rtltl1", "id"))
				monitoring.DELETE("/rtltl1/:id", moduleHandler.ResourceDelete("rtltl1", "id"))

				monitoring.GET("/alarm_history", moduleHandler.ResourceList("alarm_history"))
				monitoring.GET("/alarm_history/:id", moduleHandler.ResourceGetByID("alarm_history", "id"))
				monitoring.POST("/alarm_history", moduleHandler.ResourceCreate("alarm_history"))
				monitoring.PUT("/alarm_history/:id", moduleHandler.ResourceUpdate("alarm_history", "id"))
				monitoring.DELETE("/alarm_history/:id", moduleHandler.ResourceDelete("alarm_history", "id"))
			}

			// =========================================
			// Fase 3: Recipe & Order Module
			// =========================================
			recipe := auth.Group("/recipe")
			{
				recipe.GET("/recipe1", moduleHandler.ResourceList("recipe1"))
				recipe.GET("/recipe1/:id", moduleHandler.ResourceGetByID("recipe1", "id"))
				recipe.POST("/recipe1", moduleHandler.ResourceCreate("recipe1"))
				recipe.PUT("/recipe1/:id", moduleHandler.ResourceUpdate("recipe1", "id"))
				recipe.DELETE("/recipe1/:id", moduleHandler.ResourceDelete("recipe1", "id"))

				recipe.GET("/recipe1queue", moduleHandler.ResourceList("recipe1queue"))
				recipe.GET("/recipe1queue/:id", moduleHandler.ResourceGetByID("recipe1queue", "id"))
				recipe.POST("/recipe1queue", moduleHandler.ResourceCreate("recipe1queue"))
				recipe.PUT("/recipe1queue/:id", moduleHandler.ResourceUpdate("recipe1queue", "id"))
				recipe.DELETE("/recipe1queue/:id", moduleHandler.ResourceDelete("recipe1queue", "id"))

				recipe.GET("/recipe_history", moduleHandler.ResourceList("recipe_history"))
				recipe.GET("/recipe_history/:id", moduleHandler.ResourceGetByID("recipe_history", "id"))
				recipe.POST("/recipe_history", moduleHandler.ResourceCreate("recipe_history"))
				recipe.PUT("/recipe_history/:id", moduleHandler.ResourceUpdate("recipe_history", "id"))
				recipe.DELETE("/recipe_history/:id", moduleHandler.ResourceDelete("recipe_history", "id"))

				recipe.GET("/order_report", moduleHandler.ResourceList("order_report"))
				recipe.GET("/order_report/:id", moduleHandler.ResourceGetByID("order_report", "id"))
				recipe.POST("/order_report", moduleHandler.ResourceCreate("order_report"))
				recipe.PUT("/order_report/:id", moduleHandler.ResourceUpdate("order_report", "id"))
				recipe.DELETE("/order_report/:id", moduleHandler.ResourceDelete("order_report", "id"))

				recipe.GET("/batch_report", moduleHandler.ResourceList("batch_report"))
				recipe.GET("/batch_report/:id", moduleHandler.ResourceGetByID("batch_report", "id"))
				recipe.POST("/batch_report", moduleHandler.ResourceCreate("batch_report"))
				recipe.PUT("/batch_report/:id", moduleHandler.ResourceUpdate("batch_report", "id"))
				recipe.DELETE("/batch_report/:id", moduleHandler.ResourceDelete("batch_report", "id"))
			}

			// =========================================
			// Fase 3: Supporting / Master Module
			// =========================================
			master := auth.Group("/master")
			{
				master.GET("/mastermcn", moduleHandler.ResourceList("mastermcn"))
				master.GET("/mastermcn/:id", moduleHandler.ResourceGetByID("mastermcn", "recid"))
				master.POST("/mastermcn", moduleHandler.ResourceCreate("mastermcn"))
				master.PUT("/mastermcn/:id", moduleHandler.ResourceUpdate("mastermcn", "recid"))
				master.DELETE("/mastermcn/:id", moduleHandler.ResourceDelete("mastermcn", "recid"))

				master.GET("/bpbl", moduleHandler.ResourceList("bpbl"))
				master.GET("/bpbl/:id", moduleHandler.ResourceGetByID("bpbl", "recid"))
				master.POST("/bpbl", moduleHandler.ResourceCreate("bpbl"))
				master.PUT("/bpbl/:id", moduleHandler.ResourceUpdate("bpbl", "recid"))
				master.DELETE("/bpbl/:id", moduleHandler.ResourceDelete("bpbl", "recid"))

				master.GET("/rsc_pc1", moduleHandler.ResourceList("rsc_pc1"))
				master.GET("/rsc_pc1/:id", moduleHandler.ResourceGetByID("rsc_pc1", "id"))
				master.POST("/rsc_pc1", moduleHandler.ResourceCreate("rsc_pc1"))
				master.PUT("/rsc_pc1/:id", moduleHandler.ResourceUpdate("rsc_pc1", "id"))
				master.DELETE("/rsc_pc1/:id", moduleHandler.ResourceDelete("rsc_pc1", "id"))

				master.GET("/material", moduleHandler.ResourceList("material"))
				master.GET("/material/:id", moduleHandler.ResourceGetByID("material", "recid"))
				master.POST("/material", moduleHandler.ResourceCreate("material"))
				master.PUT("/material/:id", moduleHandler.ResourceUpdate("material", "recid"))
				master.DELETE("/material/:id", moduleHandler.ResourceDelete("material", "recid"))
			}
		}
	}
}
