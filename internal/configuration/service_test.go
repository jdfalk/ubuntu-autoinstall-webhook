// internal/configuration/service_test.go
package configuration_test

import (
	"testing"

	"github.com/jdfalk/ubuntu-autoinstall-webhook/internal/configuration"
	"github.com/spf13/viper"
)

func TestLoadConfig(t *testing.T) {
	// Add an additional config path to locate the config file.
	// Adjust the path if necessary based on your project structure.
	viper.AddConfigPath("../..")

	configService := configuration.NewService()
	cfg, err := configService.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() returned error: %v", err)
	}
	// Basic checks: ensure the expected fields are set.
	if cfg.Database.Type == "" {
		t.Errorf("Expected Database.Type to be set, got empty")
	}
	if cfg.FileEditor.BasePath == "" {
		t.Errorf("Expected FileEditor.BasePath to be set, got empty")
	}
}

func TestValidateConfig(t *testing.T) {
	configService := configuration.NewService()
	// Create a valid config.
	validCfg := configuration.Config{
		Database: configuration.DatabaseConfig{
			Type: "sqlite",
		},
		FileEditor: configuration.FileEditorConfig{
			BasePath: "/var/www/html/ipxe/boot",
		},
	}
	if err := configService.ValidateConfig(validCfg); err != nil {
		t.Errorf("ValidateConfig() returned error on valid config: %v", err)
	}

	// Create an invalid config.
	invalidCfg := configuration.Config{
		Database: configuration.DatabaseConfig{
			Type: "",
		},
		FileEditor: configuration.FileEditorConfig{
			BasePath: "",
		},
	}
	if err := configService.ValidateConfig(invalidCfg); err == nil {
		t.Errorf("ValidateConfig() did not return an error for invalid config")
	}
}

func TestGenerateTemplates(t *testing.T) {
	configService := configuration.NewService()
	data := map[string]string{
		"InstanceID": "jf-123456",
		"Hostname":   "test-host",
	}
	templates, err := configService.GenerateTemplates(data)
	if err != nil {
		t.Fatalf("GenerateTemplates() returned error: %v", err)
	}
	if len(templates) == 0 {
		t.Errorf("Expected non-empty templates map")
	}
}
