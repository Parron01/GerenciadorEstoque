package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/config"
)

// BackupManager handles database backups
type BackupManager struct {
    Config *config.Config
}

// NewBackupManager creates a new backup manager
func NewBackupManager(cfg *config.Config) *BackupManager {
    return &BackupManager{
        Config: cfg,
    }
}

// CreateBackup creates a database backup
func (bm *BackupManager) CreateBackup() (string, error) {
    // Create backup directory if it doesn't exist
    backupDir := "backups"
    if err := os.MkdirAll(backupDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create backup directory: %w", err)
    }

    // Generate backup filename with timestamp
    timestamp := time.Now().Format("2006-01-02")
    backupFile := filepath.Join(backupDir, fmt.Sprintf("inventory-backup-%s.sql", timestamp))

    // Execute pg_dump command
    cmd := exec.Command(
        "pg_dump",
        "-h", bm.Config.DBConfig.Host,
        "-p", bm.Config.DBConfig.Port,
        "-U", bm.Config.DBConfig.User,
        "-F", "c", // Custom format
        "-b", // Binary format
        "-v", // Verbose
        "-f", backupFile,
        bm.Config.DBConfig.DBName,
    )

    // Set PGPASSWORD environment variable
    cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", bm.Config.DBConfig.Password))

    // Execute command
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("backup failed: %w, output: %s", err, string(output))
    }

    log.Printf("Backup created successfully: %s", backupFile)
    
    // Clean up old backups
    err = bm.cleanupOldBackups(backupDir)
    if err != nil {
        log.Printf("Failed to clean up old backups: %v", err)
    }

    return backupFile, nil
}

// cleanupOldBackups removes backups older than 30 days
func (bm *BackupManager) cleanupOldBackups(backupDir string) error {
    // Keep backups for 30 days
    cutoffTime := time.Now().AddDate(0, 0, -30)

    entries, err := os.ReadDir(backupDir)
    if err != nil {
        return fmt.Errorf("failed to read backup directory: %w", err)
    }

    for _, entry := range entries {
        if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
            filePath := filepath.Join(backupDir, entry.Name())
            info, err := entry.Info()
            if err != nil {
                log.Printf("Failed to get file info for %s: %v", filePath, err)
                continue
            }

            if info.ModTime().Before(cutoffTime) {
                if err := os.Remove(filePath); err != nil {
                    log.Printf("Failed to remove old backup %s: %v", filePath, err)
                } else {
                    log.Printf("Removed old backup: %s", filePath)
                }
            }
        }
    }

    return nil
}

// Helper function to get environment variable with default
func getEnvWithDefault(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}