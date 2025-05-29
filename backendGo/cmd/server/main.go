package main

import (
	"database/sql"
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
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // Driver for postgres
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Driver for file source
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
)

func waitForDB(cfg *config.Config) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.DBName,
	)

	log.Printf("Tentando conectar ao banco de dados em %s:%s", cfg.DBConfig.Host, cfg.DBConfig.Port)

	maxRetries := 30
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err := sql.Open("postgres", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Printf("Conexão com banco de dados estabelecida após %d tentativas", i+1)
				db.Close()
				return nil
			}
		}

		if i < maxRetries-1 {
			log.Printf("Tentativa %d falhou: %v. Tentando novamente em %v...", i+1, err, retryInterval)
			time.Sleep(retryInterval)
		} else {
			return fmt.Errorf("não foi possível conectar ao banco de dados após %d tentativas: %v", maxRetries, err)
		}
	}

	return nil
}

func runMigrations(cfg *config.Config) {
	migrationPath := "file://migrations" // Correct: resolves to ./migrations relative to workdir

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host, // This should be 'db' when running in Docker via docker-compose
		cfg.DBConfig.Port,
		cfg.DBConfig.DBName,
	)

	log.Printf("Attempting to run migrations from %s on database %s (host: %s)", migrationPath, cfg.DBConfig.DBName, cfg.DBConfig.Host)

	m, err := migrate.New(migrationPath, dsn)
	if err != nil {
		log.Fatalf("Failed to initialize migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	} else if err == migrate.ErrNoChange {
		log.Println("No new migrations to apply.")
	} else {
		log.Println("Migrations applied successfully.")
	}
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set default environment
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "development")
	}
	fmt.Printf("Servidor rodando em modo: %s\n", os.Getenv("GO_ENV"))

	// Wait for database to be ready
	if err := waitForDB(cfg); err != nil {
		log.Fatalf("Falha ao aguardar o banco de dados: %v", err)
	}

	// Run database migrations before initializing DB for the app
	// This is important if InitDB relies on tables being present
	// Note: The DB connection for migrations is separate from database.DB for now.
	// Ensure the DB exists before running migrations. Docker-compose handles DB creation.
	runMigrations(cfg)

	// Initialize database connection for the application
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database for application: %v", err)
	}

	// Setup Gin router
	r := gin.Default()

	// Setup CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // For development. In production, specify your frontend origin.
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Operation-Batch-ID"}, // Added X-Operation-Batch-ID
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup routes
	// Repository and service initialization is now handled within SetupRoutes or passed to it.
	// For this structure, SetupRoutes takes care of it using the global database.DB.
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