package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Server ServerConfig `yaml:"server"`
	FTP    FTPConfig    `yaml:"ftp"`
	Auth   AuthConfig   `yaml:"auth"`
	Log    LogConfig    `yaml:"log"`
}

// ServerConfig contains server-specific configuration
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// FTPConfig contains FTP-specific configuration
type FTPConfig struct {
	RootDir        string `yaml:"root_dir"`
	MaxConnections int    `yaml:"max_connections"`
	Timeout        int    `yaml:"timeout"`
}

// AuthConfig contains authentication configuration
type AuthConfig struct {
	Anonymous bool              `yaml:"anonymous"`
	Users     map[string]string `yaml:"users"`
}

// LogConfig contains logging configuration
type LogConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load reads configuration from a file
func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save writes configuration to a file
func (c *Config) Save(filename string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 2121,
		},
		FTP: FTPConfig{
			RootDir:        "./ftp_root",
			MaxConnections: 100,
			Timeout:        300,
		},
		Auth: AuthConfig{
			Anonymous: true,
			Users: map[string]string{
				"anonymous": "anonymous",
			},
		},
		Log: LogConfig{
			Level:  "info",
			Format: "text",
		},
	}
}
