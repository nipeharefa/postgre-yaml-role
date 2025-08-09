package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	// This is a simple test to ensure the Config struct is correctly defined
	cfg := &Config{}
	
	// Test that the struct was created successfully
	_ = cfg
	
	// Test that the struct has the expected fields
	if cfg.Database.Host != "" {
		t.Error("Expected empty Host field")
	}
}

func TestConnectionString(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "testdb",
			SSLMode:  "disable",
		},
	}

	expected := "host=localhost port=5432 user=postgres password=postgres dbname=testdb sslmode=disable"
	actual := cfg.ConnectionString()

	if actual != expected {
		t.Errorf("Expected connection string '%s', got '%s'", expected, actual)
	}
}