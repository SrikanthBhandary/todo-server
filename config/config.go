package config

import (
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Config holds the configuration settings for the application.
type Config struct {
	// Port is the server port
	Port string `yaml:"port"`

	// DSN specifies the database address
	DSN string `yaml:"dsn"`

	//JwtSecretKey specifies the JwtSecretKey
	JwtSecretKey string `yaml:"jwt_secret_key"`

	// RedisAddress specifies the redis address
	RedisAddress string `yaml:"redis_address"`

	// HtmlAssetsPath to store the html files
	HtmlAssetsPath string `yaml:"html_assets_path"`

	// NumOfWorkers specifies the number of worker threads to be used.
	// It determines how many concurrent tasks can be processed.
	NumOfWorkers int `yaml:"num_of_workers"`

	// SmtpHost is the address of the SMTP server used for sending emails.
	// This should be set to the hostname or IP address of the SMTP server.
	SmtpHost string `yaml:"smtp_host"`

	// SmtpPort is the port number on which the SMTP server is listening.
	// Commonly used ports for SMTP are 25, 465, and 587.
	SmtpPort int `yaml:"smtp_port"`

	// SmtpUserName is the username for authenticating with the SMTP server.
	// This should be a valid email address or username as required by the SMTP service.
	SmtpUserName string `yaml:"smtp_user_name"`

	// SmtpPassword is the password associated with the SmtpUserName.
	// It is used for authenticating to the SMTP server and should be kept secure.
	SmtpPassword string `yaml:"smtp_password"`

	PDFOutputPath string `yaml:"pdf_output_path"`
}

// GetDefaultConfig returns a Config instance with default values.
func GetDefaultConfig() *Config {
	return &Config{

		NumOfWorkers: 3,
		SmtpHost:     "localhost",
		SmtpPort:     587,
		SmtpUserName: "test",
		SmtpPassword: "test",
	}
}

// NewConfig reads a YAML configuration file and returns the corresponding Config.
func NewConfig(filePath string) (*Config, error) {
	if filePath == "" {
		filePath = "./config.yaml"
	}

	configData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling config data: %w", err)
	}

	// Basic validation
	if config.SmtpHost == "" || config.SmtpUserName == "" {
		return GetDefaultConfig(), nil // fall back to defaults if necessary
	}

	return config, nil
}
