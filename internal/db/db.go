package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	_ "github.com/lib/pq" // PostgreSQL driver for CockroachDB
	"github.com/spf13/viper"
)

var DB *sql.DB

// InitDB initializes the CockroachDB connection with connection pooling
func InitDB() error {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.dbname"),
		viper.GetString("database.sslmode"),
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pooling settings
	DB.SetMaxOpenConns(viper.GetInt("database.max_open_conns"))
	DB.SetMaxIdleConns(viper.GetInt("database.max_idle_conns"))
	DB.SetConnMaxLifetime(time.Duration(viper.GetInt("database.conn_max_lifetime")) * time.Second)

	// Verify the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to CockroachDB successfully!")

	// Run migrations
	if err = RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed.")
	}
}

// RunMigrations executes schema migrations from the migrations directory
func RunMigrations() error {
	migrationsDir := "internal/db/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Ensure the migrations table exists
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		executed_at TIMESTAMP DEFAULT current_timestamp
	);`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrationName := file.Name()

			// Check if migration has already been applied
			var count int
			err := DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", migrationName).Scan(&count)
			if err != nil {
				return fmt.Errorf("failed to check migration status: %w", err)
			}

			if count > 0 {
				log.Printf("Skipping already applied migration: %s", migrationName)
				continue
			}

			// Read and execute migration
			data, err := os.ReadFile(filepath.Join(migrationsDir, migrationName))
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", migrationName, err)
			}

			_, err = DB.Exec(string(data))
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", migrationName, err)
			}

			// Record migration
			_, err = DB.Exec("INSERT INTO migrations (name) VALUES ($1)", migrationName)
			if err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migrationName, err)
			}

			log.Printf("Applied migration: %s", migrationName)
		}
	}

	log.Println("Database migrations applied successfully.")
	return nil
}
