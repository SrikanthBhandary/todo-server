package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := GetDefaultConfig()
	if config.NumOfWorkers != 3 {
		t.Errorf("default config number of workers is not matching, expected %d and actual is %d", 3, config.NumOfWorkers)
	}
	if config.SmtpHost != "localhost" {
		t.Errorf("default config SmtpHost is not matching, expected %s and actual is %s", "localhost", config.SmtpHost)
	}
}

func TestNewConfig(t *testing.T) {
	configContent := `
num_of_workers: 5
smtp_host: smtp.example.com
smtp_port: 587
smtp_user_name: user@example.com
smtp_password: password123
`

	tmpFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up
	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temporary file: %v", err)
	}

	// Test the NewConfig function
	config, err := NewConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if config.NumOfWorkers != 5 {
		t.Errorf("expected NumOfWorkers to be 5, got %d", config.NumOfWorkers)
	}
	if config.SmtpHost != "smtp.example.com" {
		t.Errorf("expected SmtpHost to be 'smtp.example.com', got '%s'", config.SmtpHost)
	}
	if config.SmtpPort != 587 {
		t.Errorf("expected SmtpPort to be 587, got %d", config.SmtpPort)
	}
	if config.SmtpUserName != "user@example.com" {
		t.Errorf("expected SmtpUserName to be 'user@example.com', got '%s'", config.SmtpUserName)
	}
	if config.SmtpPassword != "password123" {
		t.Errorf("expected SmtpPassword to be 'password123', got '%s'", config.SmtpPassword)
	}

}
