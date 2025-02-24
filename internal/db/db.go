package db

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger" // Import the logger package
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var DB *sql.DB

//go:embed migrations/*.sql
var migrationFiles embed.FS

// InitDB initializes the database connection and runs migrations.
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

	DB.SetMaxOpenConns(viper.GetInt("database.max_open_conns"))
	DB.SetMaxIdleConns(viper.GetInt("database.max_idle_conns"))
	DB.SetConnMaxLifetime(time.Duration(viper.GetInt("database.conn_max_lifetime")) * time.Second)

	// Ensure the DB connection is valid.
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Inject the DB connection into the logger package.
	logger.SetDBExecutor(DB)

	// Log the successful connection using the new logger.
	logger.Info("Connected to CockroachDB successfully!")

	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// runMigrations executes all embedded SQL migration files.
func runMigrations() error {
	files, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read embedded migrations: %w", err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		executed_at TIMESTAMP DEFAULT current_timestamp
	);`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationName := file.Name()

			var count int
			err := DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = $1", migrationName).Scan(&count)
			if err != nil {
				return fmt.Errorf("failed to check migration status: %w", err)
			}

			if count > 0 {
				logger.Infof("Skipping already applied migration: %s", migrationName)
				continue
			}

			data, err := migrationFiles.ReadFile("migrations/" + migrationName)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", migrationName, err)
			}

			_, err = DB.Exec(string(data))
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", migrationName, err)
			}

			_, err = DB.Exec("INSERT INTO migrations (name) VALUES ($1)", migrationName)
			if err != nil {
				return fmt.Errorf("failed to record migration %s: %w", migrationName, err)
			}

			logger.Infof("Applied migration: %s", migrationName)
		}
	}

	logger.Info("Database migrations applied successfully.")
	return nil
}

// SaveNetworkInterface inserts or updates a network interface in the database.
func SaveNetworkInterface(clientID, macAddress, interfaceName, chipset, driver string) error {
	query := `
		INSERT INTO network_interfaces (client_id, mac_address, interface_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (mac_address) DO UPDATE SET interface_name = EXCLUDED.interface_name
		RETURNING id;
	`
	var networkID string
	err := DB.QueryRow(query, clientID, macAddress, interfaceName).Scan(&networkID)
	if err != nil {
		return fmt.Errorf("error inserting network interface: %w", err)
	}

	query = `
		INSERT INTO network_chipsets (network_interface_id, chipset)
		VALUES ($1, $2)
		ON CONFLICT (network_interface_id) DO UPDATE SET chipset = EXCLUDED.chipset;
	`
	_, err = DB.Exec(query, networkID, chipset)
	if err != nil {
		return fmt.Errorf("error inserting chipset: %w", err)
	}

	query = `
		UPDATE network_interfaces SET driver = $1 WHERE id = $2;
	`
	_, err = DB.Exec(query, driver, networkID)
	if err != nil {
		return fmt.Errorf("error updating network driver: %w", err)
	}

	return nil
}

// SaveCloudInitVersion stores a new cloud-init configuration and maintains only the last five versions.
func SaveCloudInitVersion(clientID, macAddress, userData string) error {
	query := `
		INSERT INTO cloud_init_history (client_id, mac_address, user_data)
		VALUES ($1, $2, $3);
	`
	_, err := DB.Exec(query, clientID, macAddress, userData)
	if err != nil {
		return fmt.Errorf("error inserting cloud-init history: %w", err)
	}

	query = `
		DELETE FROM cloud_init_history
		WHERE client_id = $1 AND mac_address = $2
		AND id NOT IN (
			SELECT id FROM cloud_init_history
			WHERE client_id = $1 AND mac_address = $2
			ORDER BY created_at DESC
			LIMIT 5
		);
	`
	_, err = DB.Exec(query, clientID, macAddress)
	if err != nil {
		return fmt.Errorf("error pruning old cloud-init versions: %w", err)
	}

	return nil
}

// SaveIPXEConfigVersion stores a new version of the IPXE configuration for a client,
// and ensures that only the last five versions are kept.
func SaveIPXEConfigVersion(clientID, config string) error {
	query := `
		INSERT INTO ipxe_history (client_id, config)
		VALUES ($1, $2);
	`
	_, err := DB.Exec(query, clientID, config)
	if err != nil {
		return fmt.Errorf("error inserting IPXE history: %w", err)
	}

	// Prune old versions, keeping only the most recent five.
	query = `
		DELETE FROM ipxe_history
		WHERE client_id = $1
		AND id NOT IN (
			SELECT id FROM ipxe_history
			WHERE client_id = $1
			ORDER BY created_at DESC
			LIMIT 5
		);
	`
	_, err = DB.Exec(query, clientID)
	if err != nil {
		return fmt.Errorf("error pruning old IPXE history versions: %w", err)
	}

	return nil
}

// --- Type definitions ---

// ClientLogDetail provides detailed information for a client log.
type ClientLogDetail struct {
	ID        int
	ClientID  string
	Detail    string
	CreatedAt time.Time
}

// IpxeConfig represents an iPXE configuration.
type IpxeConfig struct {
	ID        int
	ClientID  string
	Config    string
	CreatedAt time.Time
}

// CloudInitConfig represents a cloud-init configuration.
type CloudInitConfig struct {
	ID         int
	ClientID   string
	MacAddress string
	UserData   string
	CreatedAt  time.Time
}

// --- DB Functions ---

// GetClientLogs returns a list of client logs.
func GetClientLogs() ([]ClientLog, error) {
	query := `SELECT id, client_id, timestamp, origin, description, name, result, event_type, files, created_at FROM client_logs ORDER BY created_at DESC`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying client logs: %w", err)
	}
	defer rows.Close()

	var logs []ClientLog
	for rows.Next() {
		var logEntry ClientLog
		if err := rows.Scan(&logEntry.ID, &logEntry.ClientID, &logEntry.Timestamp, &logEntry.Origin, &logEntry.Description, &logEntry.Name, &logEntry.Result, &logEntry.EventType, &logEntry.Files, &logEntry.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning client log row: %w", err)
		}
		logs = append(logs, logEntry)
	}
	return logs, rows.Err()
}

// GetClientLogDetail retrieves detailed information for a specific client log by id.
func GetClientLogDetail(id string) (ClientLog, error) {
	query := `SELECT id, client_id, timestamp, origin, description, name, result, event_type, files, created_at FROM client_logs WHERE id = $1`
	var logDetail ClientLog
	err := DB.QueryRow(query, id).Scan(&logDetail.ID, &logDetail.ClientID, &logDetail.Timestamp, &logDetail.Origin, &logDetail.Description, &logDetail.Name, &logDetail.Result, &logDetail.EventType, &logDetail.Files, &logDetail.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return logDetail, fmt.Errorf("client log %s not found", id)
		}
		return logDetail, fmt.Errorf("error querying client log detail: %w", err)
	}
	return logDetail, nil
}

// GetServerLogs returns a list of server logs.
func GetServerLogs() ([]ServerLog, error) {
	query := `SELECT id, message, created_at FROM server_logs ORDER BY created_at DESC`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying server logs: %w", err)
	}
	defer rows.Close()

	var logs []ServerLog
	for rows.Next() {
		var sl ServerLog
		if err := rows.Scan(&sl.ID, &sl.Message, &sl.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning server log row: %w", err)
		}
		logs = append(logs, sl)
	}
	return logs, rows.Err()
}

// GetIpxeConfigs returns the latest iPXE configuration per client.
// This query uses PostgreSQL's DISTINCT ON to select the most recent config per client.
func GetIpxeConfigs() ([]IpxeConfig, error) {
	query := `
        SELECT DISTINCT ON (client_id) id, client_id, config, created_at
        FROM ipxe_history
        ORDER BY client_id, created_at DESC
    `
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying current iPXE configs: %w", err)
	}
	defer rows.Close()

	var configs []IpxeConfig
	for rows.Next() {
		var cfg IpxeConfig
		if err := rows.Scan(&cfg.ID, &cfg.ClientID, &cfg.Config, &cfg.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning iPXE config row: %w", err)
		}
		configs = append(configs, cfg)
	}
	return configs, rows.Err()
}

// GetHistoricalIpxeConfigs returns all historical iPXE configurations.
func GetHistoricalIpxeConfigs() ([]IpxeConfig, error) {
	query := `SELECT id, client_id, config, created_at FROM ipxe_history ORDER BY created_at DESC`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying historical iPXE configs: %w", err)
	}
	defer rows.Close()

	var configs []IpxeConfig
	for rows.Next() {
		var cfg IpxeConfig
		if err := rows.Scan(&cfg.ID, &cfg.ClientID, &cfg.Config, &cfg.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning historical iPXE config row: %w", err)
		}
		configs = append(configs, cfg)
	}
	return configs, rows.Err()
}

// GetCloudInitConfigs returns the latest cloud-init configuration per client and MAC address.
func GetCloudInitConfigs() ([]CloudInitConfig, error) {
	query := `
        SELECT DISTINCT ON (client_id, mac_address) id, client_id, mac_address, user_data, created_at
        FROM cloud_init_history
        ORDER BY client_id, mac_address, created_at DESC
    `
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying current cloud-init configs: %w", err)
	}
	defer rows.Close()

	var configs []CloudInitConfig
	for rows.Next() {
		var cfg CloudInitConfig
		if err := rows.Scan(&cfg.ID, &cfg.ClientID, &cfg.MacAddress, &cfg.UserData, &cfg.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning cloud-init config row: %w", err)
		}
		configs = append(configs, cfg)
	}
	return configs, rows.Err()
}

// GetHistoricalCloudInitConfigs returns all historical cloud-init configurations.
func GetHistoricalCloudInitConfigs() ([]CloudInitConfig, error) {
	query := `SELECT id, client_id, mac_address, user_data, created_at FROM cloud_init_history ORDER BY created_at DESC`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying historical cloud-init configs: %w", err)
	}
	defer rows.Close()

	var configs []CloudInitConfig
	for rows.Next() {
		var cfg CloudInitConfig
		if err := rows.Scan(&cfg.ID, &cfg.ClientID, &cfg.MacAddress, &cfg.UserData, &cfg.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning historical cloud-init config row: %w", err)
		}
		configs = append(configs, cfg)
	}
	return configs, rows.Err()
}
