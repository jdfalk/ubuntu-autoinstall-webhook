package db

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/logger"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

// DB is the global database connection.
var DB *sql.DB

// migrationFiles contains all embedded SQL migration files.
//
//go:embed migrations/*.sql
var migrationFiles embed.FS

// InitDB initializes the database connection and runs migrations.
func InitDB() error {
	logger.Debugf("Starting DB initialization")

	// Retrieve DB settings from Viper.
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	dbname := viper.GetString("database.dbname")
	sslmode := viper.GetString("database.sslmode")

	// Debug: Log the DB settings (sensitive details should be omitted in production)
	logger.Debugf("DB Config - host: %s, port: %d, user: %s, dbname: %s, sslmode: %s", host, port, user, dbname, sslmode)

	// Construct the DSN. (Ensure that your driver accepts the scheme "postgresql://")
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
	logger.Debugf("Constructed DSN: %s", dsn)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		logger.Errorf("Failed to open database: %v", err)
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool parameters.
	DB.SetMaxOpenConns(viper.GetInt("database.max_open_conns"))
	DB.SetMaxIdleConns(viper.GetInt("database.max_idle_conns"))
	DB.SetConnMaxLifetime(time.Duration(viper.GetInt("database.conn_max_lifetime")) * time.Second)

	// Ensure the DB connection is valid.
	if err = DB.Ping(); err != nil {
		logger.Errorf("Failed to ping database: %v", err)
		return fmt.Errorf("failed to ping database: %w", err)
	}
	logger.Debugf("DB ping successful")

	// Inject the DB connection into the logger package.
	logger.SetDBExecutor(DB)
	logger.Info("Connected to CockroachDB successfully!")

	// Run migrations.
	if err := runMigrations(); err != nil {
		logger.Errorf("Failed to run migrations: %v", err)
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

// SaveCloudInitVersion stores a new cloud-init configuration in the history table
// and maintains only the last five versions.
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

// SaveIPXEConfigVersion stores a new version of the iPXE configuration for a client,
// and ensures that only the last five versions are kept.
// (Legacy function; use SaveIPXEConfiguration for phase support.)
func SaveIPXEConfigVersion(clientID, config string) error {
	query := `
		INSERT INTO ipxe_history (client_id, config)
		VALUES ($1, $2);
	`
	_, err := DB.Exec(query, clientID, config)
	if err != nil {
		return fmt.Errorf("error inserting IPXE history: %w", err)
	}

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

// SaveIPXEConfiguration inserts a new iPXE configuration record into the database,
// including the phase field (e.g. "install" or "post-install").
func SaveIPXEConfiguration(macAddress, config, phase string) error {
	query := `
        INSERT INTO ipxe_configuration (client_id, config, phase, created_at)
        VALUES ($1, $2, $3, $4)
    `
	// For simplicity, assume client_id is the same as macAddress.
	clientID := macAddress
	_, err := DB.Exec(query, clientID, config, phase, time.Now())
	if err != nil {
		logger.Errorf("Error inserting iPXE configuration: %v", err)
		return fmt.Errorf("error inserting iPXE configuration: %w", err)
	}
	logger.Infof("Saved new iPXE configuration for MAC %s with phase %s", macAddress, phase)
	return nil
}

// GetLatestIPXEConfig retrieves the latest iPXE configuration for the given MAC address and phase.
func GetLatestIPXEConfig(macAddress, phase string) (IpxeConfig, error) {
	var cfg IpxeConfig
	// Assuming client_id is derived from macAddress.
	clientID := macAddress
	query := `
        SELECT id, client_id, config, created_at
        FROM ipxe_configuration
        WHERE client_id = $1 AND phase = $2
        ORDER BY created_at DESC
        LIMIT 1
    `
	err := DB.QueryRow(query, clientID, phase).Scan(&cfg.ID, &cfg.ClientID, &cfg.Config, &cfg.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return cfg, fmt.Errorf("no iPXE configuration found for MAC %s with phase %s", macAddress, phase)
		}
		return cfg, fmt.Errorf("error querying latest iPXE configuration: %w", err)
	}
	return cfg, nil
}

// SaveClientLog saves a client log to the database.
func SaveClientLog(event interface{}) {
	e, ok := event.(map[string]interface{})
	if !ok {
		logger.Errorf("SaveClientLog: invalid event format")
		return
	}
	// Extract required fields.
	sourceIP, _ := e["source_ip"].(string)
	timestampFloat, _ := e["timestamp"].(float64)
	origin, _ := e["origin"].(string)
	description, _ := e["description"].(string)
	name, _ := e["name"].(string)
	result, _ := e["result"].(string)
	eventType, _ := e["event_type"].(string)
	filesBytes, _ := json.Marshal(e["files"])

	query := `INSERT INTO client_logs (client_id, timestamp, origin, description, name, result, event_type, files, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`
	_, err := DB.Exec(query, sourceIP, time.Unix(int64(timestampFloat), 0), origin, description, name, result, eventType, string(filesBytes))
	if err != nil {
		logger.ClientErrorf("Error saving client log: %v", err)
	}
}

// SaveClientStatus saves a client status update to the database.
func SaveClientStatus(event interface{}) {
	e, ok := event.(map[string]interface{})
	if !ok {
		logger.Errorf("SaveClientStatus: invalid event format")
		return
	}
	sourceIP, _ := e["source_ip"].(string)
	status, _ := e["status"].(string)
	progress, _ := e["progress"].(float64)
	message, _ := e["message"].(string)

	query := `INSERT INTO client_status (client_id, status, progress, message, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (client_id) DO UPDATE
		SET status = $2, progress = $3, message = $4, updated_at = NOW();`
	_, err := DB.Exec(query, sourceIP, status, int(progress), message)
	if err != nil {
		logger.ClientErrorf("Error saving client status: %v", err)
	}
}

// CloseDB gracefully closes the database connection.
func CloseDB() error {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			logger.Errorf("Failed to close database: %v", err)
			return fmt.Errorf("failed to close database: %w", err)
		}
		logger.Info("Database connection closed successfully.")
	} else {
		logger.Infof("No database connection to close.")
	}
	return nil
}

// --- Type definitions for log queries ---

// ClientLogDetail provides detailed information for a client log.
type ClientLogDetail struct {
	ID          int
	ClientID    string
	Timestamp   time.Time
	Origin      string
	Description string
	Name        string
	Result      string
	EventType   string
	Files       string
	CreatedAt   time.Time
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

// --- DB Query Functions ---

// GetClientLogs returns a list of client logs.
func GetClientLogs() ([]ClientLogDetail, error) {
	query := `SELECT id, client_id, timestamp, origin, description, name, result, event_type, files, created_at FROM client_logs ORDER BY created_at DESC`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying client logs: %w", err)
	}
	defer rows.Close()

	var logs []ClientLogDetail
	for rows.Next() {
		var logEntry ClientLogDetail
		if err := rows.Scan(&logEntry.ID, &logEntry.ClientID, &logEntry.Timestamp, &logEntry.Origin, &logEntry.Description, &logEntry.Name, &logEntry.Result, &logEntry.EventType, &logEntry.Files, &logEntry.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning client log row: %w", err)
		}
		logs = append(logs, logEntry)
	}
	return logs, rows.Err()
}

// GetClientLogDetail retrieves detailed information for a specific client log by id.
func GetClientLogDetail(id string) (ClientLogDetail, error) {
	query := `SELECT id, client_id, timestamp, origin, description, name, result, event_type, files, created_at FROM client_logs WHERE id = $1`
	var logDetail ClientLogDetail
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

// GetIpxeConfigs returns the latest iPXE configuration per client from ipxe_history.
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
			return nil, fmt.Errorf("error scanning iPXe config row: %w", err)
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
		return nil, fmt.Errorf("error querying historical iPXe configs: %w", err)
	}
	defer rows.Close()

	var configs []IpxeConfig
	for rows.Next() {
		var cfg IpxeConfig
		if err := rows.Scan(&cfg.ID, &cfg.ClientID, &cfg.Config, &cfg.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning historical iPXe config row: %w", err)
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
