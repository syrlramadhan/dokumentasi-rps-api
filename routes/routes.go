package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/syrlramadhan/dokumentasi-rps-api/controllers"
	"github.com/syrlramadhan/dokumentasi-rps-api/repositories"
	"github.com/syrlramadhan/dokumentasi-rps-api/services"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	programRepo := repositories.NewProgramRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	templateRepo := repositories.NewTemplateRepository(db)
	templateVersionRepo := repositories.NewTemplateVersionRepository(db)
	generatedRPSRepo := repositories.NewGeneratedRPSRepository(db)
	auditLogRepo := repositories.NewAuditLogRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	programService := services.NewProgramService(programRepo)
	courseService := services.NewCourseService(courseRepo)
	templateService := services.NewTemplateService(templateRepo)
	templateVersionService := services.NewTemplateVersionService(templateVersionRepo)
	generatedRPSService := services.NewGeneratedRPSService(generatedRPSRepo)
	auditLogService := services.NewAuditLogService(auditLogRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	programController := controllers.NewProgramController(programService)
	courseController := controllers.NewCourseController(courseService)
	templateController := controllers.NewTemplateController(templateService)
	templateVersionController := controllers.NewTemplateVersionController(templateVersionService)
	generatedRPSController := controllers.NewGeneratedRPSController(generatedRPSService)
	auditLogController := controllers.NewAuditLogController(auditLogService)

	// API v1 group
	v1 := r.Group("/api/v1")
	{
		// Auth routes (TODO: implement auth controller)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Login endpoint - TODO: implement JWT auth"})
			})
			auth.POST("/refresh", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Refresh token endpoint - TODO: implement"})
			})
		}

		// Users routes
		users := v1.Group("/users")
		{
			users.POST("", userController.Create)
			users.GET("", userController.FindAll)
			users.GET("/:id", userController.FindByID)
			users.PUT("/:id", userController.Update)
			users.DELETE("/:id", userController.Delete)
		}

		// Programs routes
		programs := v1.Group("/programs")
		{
			programs.POST("", programController.Create)
			programs.GET("", programController.FindAll)
			programs.GET("/:id", programController.FindByID)
			programs.PUT("/:id", programController.Update)
			programs.DELETE("/:id", programController.Delete)
		}

		// Courses routes
		courses := v1.Group("/courses")
		{
			courses.POST("", courseController.Create)
			courses.GET("", courseController.FindAll)
			courses.GET("/:id", courseController.FindByID)
			courses.GET("/program/:program_id", courseController.FindByProgramID)
			courses.PUT("/:id", courseController.Update)
			courses.DELETE("/:id", courseController.Delete)
		}

		// Templates routes
		templates := v1.Group("/templates")
		{
			templates.POST("", templateController.Create)
			templates.GET("", templateController.FindAll)
			templates.GET("/:id", templateController.FindByID)
			templates.GET("/program/:program_id", templateController.FindByProgramID)
			templates.GET("/program/:program_id/active", templateController.FindActiveByProgramID)
			templates.PUT("/:id", templateController.Update)
			templates.DELETE("/:id", templateController.Delete)

			// Template versions (nested under templates)
			templates.POST("/:id/versions", templateVersionController.Create)
			templates.GET("/:id/versions", templateVersionController.FindByTemplateID)
			templates.GET("/:id/versions/latest", templateVersionController.FindLatestByTemplateID)
		}

		// Template versions routes (standalone)
		templateVersions := v1.Group("/template-versions")
		{
			templateVersions.GET("", templateVersionController.FindAll)
			templateVersions.GET("/:id", templateVersionController.FindByID)
			templateVersions.PUT("/:id", templateVersionController.Update)
			templateVersions.DELETE("/:id", templateVersionController.Delete)
		}

		// Generate routes
		generate := v1.Group("/generate")
		{
			generate.POST("", generatedRPSController.Create)
			generate.GET("/:job_id/status", generatedRPSController.FindByID)
		}

		// Generated RPS routes
		generated := v1.Group("/generated")
		{
			generated.GET("", generatedRPSController.FindAll)
			generated.GET("/:id", generatedRPSController.FindByID)
			generated.GET("/:id/export", generatedRPSController.Export)
			generated.GET("/course/:course_id", generatedRPSController.FindByCourseID)
			generated.GET("/status/:status", generatedRPSController.FindByStatus)
			generated.PUT("/:id", generatedRPSController.Update)
			generated.PATCH("/:id/status", generatedRPSController.UpdateStatus)
			generated.DELETE("/:id", generatedRPSController.Delete)
		}

		// Admin routes
		admin := v1.Group("/admin")
		{
			// Audit logs
			audit := admin.Group("/audit")
			{
				audit.GET("", auditLogController.FindAll)
				audit.GET("/:id", auditLogController.FindByID)
				audit.GET("/user/:user_id", auditLogController.FindByUserID)
				audit.GET("/action/:action", auditLogController.FindByAction)
				audit.GET("/date-range", auditLogController.FindByDateRange)
				audit.DELETE("/:id", auditLogController.Delete)
			}

			// AI prompts (TODO: implement)
			ai := admin.Group("/ai")
			{
				ai.GET("/prompts/:id", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "AI prompt trace - TODO: implement"})
				})
			}
		}

		// Internal routes (for worker/microservices)
		internal := v1.Group("/internal")
		{
			internal.POST("/complete_generation", generatedRPSController.CompleteGeneration)
		}
	}
}
