package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/config"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/database"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/routes"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()
    
    // Set default environment
    if os.Getenv("GO_ENV") == "" {
        os.Setenv("GO_ENV", "development")
    }
    fmt.Printf("Servidor rodando em modo: %s\n", os.Getenv("GO_ENV"))

    // Initialize database connection
    if err := database.InitDB(cfg); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    // Setup Gin router
    r := gin.Default()

    // Setup CORS
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Setup routes
    routes.SetupRoutes(r, cfg)

    // Serve static files in production
    if os.Getenv("GO_ENV") == "production" {
        r.Static("/", "./dist")
        r.NoRoute(func(c *gin.Context) {
            c.File("./dist/index.html")
        })
    }

    // Setup backup manager
    backupManager := utils.NewBackupManager(cfg)

    // Set up cron job for weekly backups (Sunday at 3:00 AM)
    c := cron.New()
    _, err := c.AddFunc("0 3 * * 0", func() {
        log.Println("Executando backup semanal...")
        backupPath, err := backupManager.CreateBackup()
        if err != nil {
            log.Printf("Erro ao criar backup: %v", err)
        } else {
            log.Printf("Backup criado com sucesso: %s", backupPath)
        }
    })
    if err != nil {
        log.Printf("Erro ao configurar agendamento de backup: %v", err)
    } else {
        c.Start()
        log.Println("Agendamento de backup configurado")
    }

    // Start the server
    port := cfg.Port
    log.Printf("Servidor rodando na porta %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}