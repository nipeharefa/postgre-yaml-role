package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
}

// DatabaseConfig represents the database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// Load loads the configuration from file or environment variables
func Load() (*Config, error) {
	// Set default values
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "postgres")
	viper.SetDefault("database.sslmode", "disable")

	// Read from environment variables
	viper.SetEnvPrefix("pgctl")
	viper.AutomaticEnv()

	// Try to read from config file
	configPath := viper.GetString("config")
	if configPath == "" {
		// Look for config file in home directory
		home, err := os.UserHomeDir()
		if err == nil {
			configPath = filepath.Join(home, ".pgctl", "config.yaml")
		}
	}

	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			// If file doesn't exist, continue with defaults and environment variables
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("error reading config file: %w", err)
			}
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// ConnectionString returns a PostgreSQL connection string
func (c *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode)
}