package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/config"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the database connection using the provided configuration
func InitDB(cfg *config.Config) error {
    // Create connection string
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.User, 
        cfg.DBConfig.Password, cfg.DBConfig.DBName,
    )

    // Connect to the database
    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("error connecting to database: %w", err)
    }

    // Test the connection
    if err = DB.Ping(); err != nil {
        return fmt.Errorf("error pinging database: %w", err)
    }

    log.Println("Database connection established successfully")

    // Initialize database tables and data
    if err = initTables(); err != nil {
        return fmt.Errorf("error initializing tables: %w", err)
    }

    // Create admin user if needed
    if err = createAdminUser(cfg.Admin.Username, cfg.Admin.Password); err != nil {
        return fmt.Errorf("error creating admin user: %w", err)
    }

    // Insert sample data if needed
    if err = insertSampleData(); err != nil {
        return fmt.Errorf("error inserting sample data: %w", err)
    }

    return nil
}

// initTables creates tables if they don't exist
func initTables() error {
    // Create users table
    _, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(100) UNIQUE NOT NULL,
            password VARCHAR(100) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        return err
    }

    // Create products table
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS products (
            id VARCHAR(100) PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            unit VARCHAR(10) NOT NULL CHECK(unit IN ('L', 'kg')),
            quantity NUMERIC NOT NULL DEFAULT 0
        )
    `)
    if err != nil {
        return err
    }

    // Create history table
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS history (
            id VARCHAR(100) PRIMARY KEY,
            date VARCHAR(100) NOT NULL,
            changes JSONB NOT NULL
        )
    `)
    if err != nil {
        return err
    }

    log.Println("Database tables created/verified successfully")
    return nil
}

// createAdminUser creates an admin user if it doesn't exist
func createAdminUser(username, password string) error {
    // Check if admin user already exists
    var count int
    err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
    if err != nil {
        return err
    }

    if count == 0 {
        // Hash the password
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }

        // Insert admin user
        _, err = DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hashedPassword))
        if err != nil {
            return err
        }
        log.Printf("Admin user created successfully: %s\n", username)
    }

    return nil
}

// insertSampleData inserts initial sample data if the products table is empty
func insertSampleData() error {
    var count int
    err := DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
    if err != nil {
        return err
    }

    if count == 0 {
        // Default products from Node.js backend
        defaultProducts := []struct {
            ID       string
            Name     string
            Unit     string
            Quantity float64
        }{
            {"1", "Alade", "L", 210},
            {"2", "Curbix", "L", 71},
            {"3", "Magnum", "kg", 110},
            {"4", "Instivo", "L", 3},
            {"5", "Kasumin", "L", 50},
            {"6", "Priori", "L", 33},
        }

        for _, product := range defaultProducts {
            _, err := DB.Exec(
                "INSERT INTO products (id, name, unit, quantity) VALUES ($1, $2, $3, $4)",
                product.ID, product.Name, product.Unit, product.Quantity,
            )
            if err != nil {
                return err
            }
        }

        log.Println("Sample data inserted into database")
    }

    return nil
}
