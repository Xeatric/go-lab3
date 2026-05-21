package main

import (
	"log"
	"os"
	"paving-tiles-api/internal/config"
	"paving-tiles-api/internal/controller"
	"paving-tiles-api/internal/database"
	"paving-tiles-api/internal/middleware"
	"paving-tiles-api/internal/repository"
	"paving-tiles-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Логируем все переменные окружения для отладки
	log.Println("=== ПЕРЕМЕННЫЕ ОКРУЖЕНИЯ ===")
	log.Printf("DB_HOST=%s", os.Getenv("DB_HOST"))
	log.Printf("DB_PORT=%s", os.Getenv("DB_PORT"))
	log.Printf("DB_USER=%s", os.Getenv("DB_USER"))
	log.Printf("DB_PASSWORD=%s", os.Getenv("DB_PASSWORD"))
	log.Printf("DB_NAME=%s", os.Getenv("DB_NAME"))
	log.Println("============================")

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	log.Printf("Конфиг после загрузки: DBName=%s", cfg.DBName)

	// Подключение к БД
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// Миграции
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Инициализация слоев
	tileRepo := repository.NewTileRepository(db)
	tileService := service.NewTileService(tileRepo)
	tileController := controller.NewTileController(tileService)

	// Настройка роутера
	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.Logger())

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Добро пожаловать в API каталога тротуарной плитки!")
	})

	api := router.Group("/api/v1")
	{
		tiles := api.Group("/tiles")
		{
			tiles.GET("", tileController.GetTiles)
			tiles.GET("/:id", tileController.GetTileByID)
			tiles.POST("", tileController.CreateTile)
			tiles.PUT("/:id", tileController.UpdateTile)
			tiles.PATCH("/:id", tileController.PatchTile)
			tiles.DELETE("/:id", tileController.DeleteTile)
		}

		api.GET("/info", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "API каталога тротуарной плитки",
				"version": "2.0.0",
			})
		})
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
