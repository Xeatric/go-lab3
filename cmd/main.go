package main

import (
	"log"

	authController "paving-tiles-api/internal/auth/controller"
	authMiddleware "paving-tiles-api/internal/auth/middleware"
	authRepo "paving-tiles-api/internal/auth/repository"
	authService "paving-tiles-api/internal/auth/service"
	"paving-tiles-api/internal/config"
	"paving-tiles-api/internal/controller"
	"paving-tiles-api/internal/database"
	"paving-tiles-api/internal/middleware"
	"paving-tiles-api/internal/repository"
	"paving-tiles-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Подключение к БД
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Миграции
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Auth слои
	authRepoInstance := authRepo.NewAuthRepository(db)
	authServiceInstance := authService.NewAuthService(authRepoInstance, cfg)
	authControllerInstance := authController.NewAuthController(authServiceInstance, cfg)

	// Auth middleware
	authMiddlewareInstance := authMiddleware.NewAuthMiddleware(authServiceInstance)

	// Business слои
	tileRepo := repository.NewTileRepository(db)
	tileService := service.NewTileService(tileRepo)
	tileController := controller.NewTileController(tileService)

	// Настройка роутера
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.Logger())

	// Публичные маршруты
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Добро пожаловать в API каталога тротуарной плитки!")
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth маршруты (публичные)
	auth := router.Group("/auth")
	{
		auth.POST("/register", authControllerInstance.Register)
		auth.POST("/login", authControllerInstance.Login)
		auth.POST("/refresh", authControllerInstance.Refresh)
		auth.GET("/oauth/:provider", authControllerInstance.OAuthLogin)
		auth.GET("/oauth/:provider/callback", authControllerInstance.OAuthCallback)
	}

	// Защищенные маршруты (требуют аутентификации)
	protected := router.Group("/api/v1")
	protected.Use(authMiddlewareInstance.Authenticate())
	{
		// Auth эндпоинты (требуют аутентификации)
		protected.GET("/auth/whoami", authControllerInstance.Whoami)
		protected.POST("/auth/logout", authControllerInstance.Logout)
		protected.POST("/auth/logout-all", authControllerInstance.LogoutAll)

		// Business эндпоинты
		tiles := protected.Group("/tiles")
		{
			tiles.GET("", tileController.GetTiles)
			tiles.GET("/:id", tileController.GetTileByID)
			tiles.POST("", tileController.CreateTile)
			tiles.PUT("/:id", tileController.UpdateTile)
			tiles.PATCH("/:id", tileController.PatchTile)
			tiles.DELETE("/:id", tileController.DeleteTile)
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
