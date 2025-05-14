package routes

import (
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/config"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/controllers"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
    // Create controllers
    authController := controllers.NewAuthController(cfg)
    productController := controllers.NewProductController()
    historyController := controllers.NewHistoryController()

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
            products.GET("", productController.GetAll)
            products.GET("/:id", productController.GetByID)
            products.POST("", middleware.AuthMiddleware(cfg), productController.Create)
            products.PUT("/:id", middleware.AuthMiddleware(cfg), productController.Update)
            products.DELETE("/:id", middleware.AuthMiddleware(cfg), productController.Delete)
        }

        // History routes
        history := api.Group("/history")
        {
            history.GET("", historyController.GetAll)
            history.POST("", middleware.AuthMiddleware(cfg), historyController.Create)
        }
    }
}