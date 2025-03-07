// internal/configuration/service.go
package configuration

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds the top-level configuration.
type Config struct {
	Environment string           `mapstructure:"environment"`
	Database    DatabaseConfig   `mapstructure:"database"`
	FileEditor  FileEditorConfig `mapstructure:"file_editor"`
	// You can add additional top-level configuration sections as needed.
}

// DatabaseConfig holds the database-related configuration.
type DatabaseConfig struct {
	Type      string          `mapstructure:"type"`
	SQLite    SQLiteConfig    `mapstructure:"sqlite"`
	Cockroach CockroachConfig `mapstructure:"cockroachdb"`
}

// SQLiteConfig holds configuration for the SQLite database.
type SQLiteConfig struct {
	Path string `mapstructure:"path"`
}

// CockroachConfig holds configuration for the CockroachDB database.
type CockroachConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// FileEditorConfig holds file editor-related configuration.
type FileEditorConfig struct {
	BasePath          string `mapstructure:"base_path"`
	CloudInitBasePath string `mapstructure:"cloud_init_base_path"`
	HostnameSymlink   bool   `mapstructure:"hostname_symlink"`
}

// ConfigService defines operations to manage configurations.
type ConfigService interface {
	LoadConfig() (Config, error)
	ValidateConfig(cfg Config) error
	WatchConfigUpdates() (<-chan Config, error)
	GenerateTemplates(data interface{}) (map[string][]byte, error)
}

// Service implements the ConfigService interface.
type Service struct{}

// NewService creates a new configuration service instance.
func NewService() ConfigService {
	return &Service{}
}

// LoadConfig reads the configuration using Viper.
func (s *Service) LoadConfig() (Config, error) {
	fmt.Println("Loading configuration from file...")
	var cfg Config

	// Set the config name and add search paths.
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	// Additional paths can be added externally (e.g., via tests).
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}
	fmt.Printf("Loaded configuration: %+v\n", cfg)
	return cfg, nil
}

// ValidateConfig performs basic validation of the configuration.
func (s *Service) ValidateConfig(cfg Config) error {
	fmt.Println("Validating configuration...")
	// Check required fields.
	if cfg.Database.Type == "" {
		return fmt.Errorf("database.type cannot be empty")
	}
	if cfg.FileEditor.BasePath == "" {
		return fmt.Errorf("file_editor.base_path cannot be empty")
	}
	return nil
}

// WatchConfigUpdates sets up a simple watcher that polls for config changes.
func (s *Service) WatchConfigUpdates() (<-chan Config, error) {
	cfgCh := make(chan Config)
	go func() {
		for {
			cfg, err := s.LoadConfig()
			if err == nil {
				cfgCh <- cfg
			}
			time.Sleep(30 * time.Second) // Poll every 30 seconds.
		}
	}()
	return cfgCh, nil
}

// GenerateTemplates creates configuration file templates based on provided data.
func (s *Service) GenerateTemplates(data interface{}) (map[string][]byte, error) {
	fmt.Println("Generating configuration templates...")
	// TODO: Implement template generation logic.
	templates := make(map[string][]byte)
	templates["meta-data"] = []byte("instance-id: {{.InstanceID}}\nlocal-hostname: {{.Hostname}}\n")
	templates["network-config"] = []byte("")
	templates["user-data"] = []byte("#cloud-config\n...")
	return templates, nil
}
