package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

// Test default values
func TestInitConfigDefaults(t *testing.T) {
	initConfig()

	if viper.GetString("port") != "5000" {
		t.Errorf("Expected port 5000, got %s", viper.GetString("port"))
	}
	if viper.GetString("logDir") != "/var/log/autoinstall-webhook" {
		t.Errorf("Expected default logDir, got %s", viper.GetString("logDir"))
	}
}
