package routes

import (
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/config"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/controllers"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/database" // For DB instance
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/middleware"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/repository" // Added
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/service"    // Added
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
    // Initialize Repositories
	loteRepository := repository.NewLoteRepository(database.DB)
	productRepository := repository.NewProductRepository(database.DB, loteRepository) // LoteRepo is a dependency for ProductRepo
	historyRepository := repository.NewHistoryRepository(database.DB)

    // Initialize Services
	historyService := service.NewHistoryService(historyRepository, productRepository) // Pass productRepository
	// Pass database.DB to LoteService for transaction management
	loteService := service.NewLoteService(loteRepository, productRepository, historyService, database.DB)


    // Create controllers
	authController := controllers.NewAuthController(cfg)
	productController := controllers.NewProductController(productRepository, historyService) // Updated
	historyController := controllers.NewHistoryController(historyService)                   // Updated
	loteController := controllers.NewLoteController(loteService)                           // Added

    // API routes
	api := router.Group("/api")
	{
        // Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.GET("/verify", middleware.AuthMiddleware(cfg), authController.Verify)
			auth.GET("/health", authController.Health)
		}

        // Product routes
		products := api.Group("/products")
		{
			products.GET("", middleware.AuthMiddleware(cfg), productController.GetAll)
			products.GET("/:product_id", middleware.AuthMiddleware(cfg), productController.GetByID) // Changed :id to :product_id
			products.POST("", middleware.AuthMiddleware(cfg), productController.Create)
			products.PUT("/:product_id", middleware.AuthMiddleware(cfg), productController.Update) // Changed :id to :product_id
			products.DELETE("/:product_id", middleware.AuthMiddleware(cfg), productController.Delete) // Changed :id to :product_id

            // Lote routes (nested under products for creation and listing)
			products.POST("/:product_id/lotes", middleware.AuthMiddleware(cfg), loteController.CreateLote)
			products.GET("/:product_id/lotes", middleware.AuthMiddleware(cfg), loteController.GetLotesForProduct)
		}

        // Standalone Lote routes (for updating/deleting specific lotes by their own ID)
		lotes := api.Group("/lotes")
		{
			// GET /lotes/:lote_id could be added if needed, but GetLotesForProduct might be sufficient
			lotes.PUT("/:lote_id", middleware.AuthMiddleware(cfg), loteController.UpdateLote)
			lotes.DELETE("/:lote_id", middleware.AuthMiddleware(cfg), loteController.DeleteLote)
		}


        // History routes
		history := api.Group("/history")
		{
			history.GET("", middleware.AuthMiddleware(cfg), historyController.GetAll) // Now supports ?batch_id=
			history.POST("", middleware.AuthMiddleware(cfg), historyController.Create)
			history.GET("/:entity_type/:entity_id", middleware.AuthMiddleware(cfg), historyController.GetHistoryForEntity)
			
			// New batch endpoints
			history.POST("/batch", middleware.AuthMiddleware(cfg), historyController.CreateBatch)
			history.GET("/batch/:batch_id", middleware.AuthMiddleware(cfg), historyController.GetByBatch)
			history.GET("/grouped", middleware.AuthMiddleware(cfg), historyController.GetGrouped) 
			history.POST("/product-context", middleware.AuthMiddleware(cfg), historyController.CreateProductBatchContext) // New route
		}
	}
}