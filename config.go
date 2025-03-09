package lashes

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/greysquirr3l/lashes/internal/rotation"
	"github.com/greysquirr3l/lashes/internal/storage"
)

// Strategy name constants
const (
	StrategyRoundRobin = "round-robin"
	StrategyRandom     = "random"
	StrategyWeighted   = "weighted"
	StrategyLeastUsed  = "least-used"

	// Alternative names for backward compatibility
	StrategyRoundRobinAlt = "roundrobin"
	StrategyLeastUsedAlt  = "leastused"
)

// Config represents the configuration for the proxy rotator
type Config struct {
	Storage struct {
		Type             string `json:"type"`
		ConnectionString string `json:"connection_string,omitempty"`
		FilePath         string `json:"file_path,omitempty"`
		QueryTimeout     string `json:"query_timeout,omitempty"`
	} `json:"storage"`

	Strategy        string `json:"strategy"`
	TestURL         string `json:"test_url"`
	ValidateOnStart bool   `json:"validate_on_start"`

	Timeouts struct {
		Request    string `json:"request"`
		Validation string `json:"validation"`
		Retry      string `json:"retry"`
	} `json:"timeouts"`

	MaxRetries int `json:"max_retries"`

	CircuitBreaker struct {
		Enabled             bool   `json:"enabled"`
		MaxFailures         int    `json:"max_failures"`
		ResetTimeout        string `json:"reset_timeout"`
		EnableGlobalBreaker bool   `json:"enable_global_breaker"`
	} `json:"circuit_breaker"`
}

// LoadConfig loads configuration from a file
func LoadConfig(filePath string) (Options, error) {
	// Split this complex function into smaller parts
	options := DefaultOptions()

	config, err := readConfigFile(filePath)
	if err != nil {
		return options, err
	}

	return applyConfig(config, options)
}

// readConfigFile reads and parses a config file
func readConfigFile(filePath string) (Config, error) {
	var config Config

	// #nosec G304 - filePath comes from trusted source in this context
	file, err := os.Open(filePath)
	if err != nil {
		return config, fmt.Errorf("error opening config file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log or handle close error if needed
			// This approach avoids overwriting any existing error
			err = closeErr
		}
	}()

	// Read the file content
	bytes, err := io.ReadAll(file)
	if err != nil {
		return config, fmt.Errorf("error reading config file: %w", err)
	}

	// Parse the configuration
	if err := json.Unmarshal(bytes, &config); err != nil {
		return config, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}

// applyConfig applies config values to options
func applyConfig(config Config, options Options) (Options, error) {
	// Apply storage configuration
	options = applyStorageConfig(config, options)

	// Apply strategy
	options = applyStrategyConfig(config, options)

	// Apply timeouts
	options = applyTimeoutConfig(config, options)

	// Apply other settings
	if config.TestURL != "" {
		options.TestURL = config.TestURL
	}

	options.ValidateOnStart = config.ValidateOnStart

	if config.MaxRetries > 0 {
		options.MaxRetries = config.MaxRetries
	}

	return options, nil
}

// applyStorageConfig configures storage options
func applyStorageConfig(config Config, options Options) Options {
	switch config.Storage.Type {
	case "sqlite", "mysql", "postgres":
		storageType := storage.DatabaseType(config.Storage.Type)
		options.Storage = &storage.Options{
			Type:             storageType,
			FilePath:         config.Storage.FilePath,
			ConnectionString: config.Storage.ConnectionString,
		}

		if config.Storage.QueryTimeout != "" {
			if timeout, err := time.ParseDuration(config.Storage.QueryTimeout); err == nil {
				options.Storage.QueryTimeout = timeout
			}
		}
	}
	return options
}

// applyStrategyConfig configures rotation strategy
func applyStrategyConfig(config Config, options Options) Options {
	switch config.Strategy {
	case StrategyRoundRobin, StrategyRoundRobinAlt:
		options.Strategy = rotation.RoundRobinStrategy
	case StrategyRandom:
		options.Strategy = rotation.RandomStrategy
	case StrategyWeighted:
		options.Strategy = rotation.WeightedStrategy
	case StrategyLeastUsed, StrategyLeastUsedAlt:
		options.Strategy = rotation.LeastUsedStrategy
	}
	return options
}

// applyTimeoutConfig configures timeout settings
func applyTimeoutConfig(config Config, options Options) Options {
	if config.Timeouts.Request != "" {
		if timeout, err := time.ParseDuration(config.Timeouts.Request); err == nil {
			options.RequestTimeout = timeout
		}
	}

	if config.Timeouts.Validation != "" {
		if timeout, err := time.ParseDuration(config.Timeouts.Validation); err == nil {
			options.ValidationTimeout = timeout
		}
	}

	if config.Timeouts.Retry != "" {
		if timeout, err := time.ParseDuration(config.Timeouts.Retry); err == nil {
			options.RetryDelay = timeout
		}
	}
	return options
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() Options {
	// Split into smaller functions
	options := DefaultOptions()

	// Apply different config types
	options = applyStorageEnvConfig(options)
	options = applyStrategyEnvConfig(options)
	options = applyTimeoutEnvConfig(options)

	return options
}

// applyStorageEnvConfig applies storage configuration from environment
func applyStorageEnvConfig(options Options) Options {
	storageType := os.Getenv("LASHES_STORAGE_TYPE")
	switch storageType {
	case "sqlite":
		options.Storage = &storage.Options{
			Type:     storage.SQLite,
			FilePath: os.Getenv("LASHES_SQLITE_PATH"),
		}
	case "mysql":
		options.Storage = &storage.Options{
			Type:             storage.MySQL,
			ConnectionString: os.Getenv("LASHES_MYSQL_DSN"),
		}
	case "postgres":
		options.Storage = &storage.Options{
			Type:             storage.Postgres,
			ConnectionString: os.Getenv("LASHES_POSTGRES_DSN"),
		}
	}
	return options
}

// applyStrategyEnvConfig applies rotation strategy configuration from environment
func applyStrategyEnvConfig(options Options) Options {
	strategyEnv := os.Getenv("LASHES_STRATEGY")
	switch strategyEnv {
	case StrategyRoundRobin, StrategyRoundRobinAlt:
		options.Strategy = rotation.RoundRobinStrategy
	case StrategyRandom:
		options.Strategy = rotation.RandomStrategy
	case StrategyWeighted:
		options.Strategy = rotation.WeightedStrategy
	case StrategyLeastUsed, StrategyLeastUsedAlt:
		options.Strategy = rotation.LeastUsedStrategy
	}

	// Test URL
	if testURL := os.Getenv("LASHES_TEST_URL"); testURL != "" {
		options.TestURL = testURL
	}

	return options
}

// applyTimeoutEnvConfig applies timeout configuration from environment
func applyTimeoutEnvConfig(options Options) Options {
	// Timeouts
	if timeout := os.Getenv("LASHES_REQUEST_TIMEOUT"); timeout != "" {
		if parsed, err := time.ParseDuration(timeout); err == nil {
			options.RequestTimeout = parsed
		}
	}

	if timeout := os.Getenv("LASHES_VALIDATION_TIMEOUT"); timeout != "" {
		if parsed, err := time.ParseDuration(timeout); err == nil {
			options.ValidationTimeout = parsed
		}
	}

	return options
}

// SaveConfig saves the current options to a configuration file
func SaveConfig(options Options, configFilePath string) error {
	// Create a config structure from the options
	config := buildConfigFromOptions(options)

	// Convert to JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("error creating JSON config: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(configFilePath)
	// #nosec G301 - This directory needs to be accessible by the application
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("error creating config directory: %w", err)
	}

	// Write to file
	// #nosec G306 - This file needs to be readable by the application
	if err := os.WriteFile(configFilePath, data, 0o600); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

// buildConfigFromOptions constructs a Config from Options
func buildConfigFromOptions(options Options) Config {
	config := Config{}

	// Storage settings
	if options.Storage != nil {
		config.Storage.Type = string(options.Storage.Type)
		config.Storage.FilePath = options.Storage.FilePath
		config.Storage.ConnectionString = options.Storage.ConnectionString
		config.Storage.QueryTimeout = options.Storage.QueryTimeout.String()
	}

	// Strategy
	switch options.Strategy {
	case rotation.RoundRobinStrategy:
		config.Strategy = StrategyRoundRobin
	case rotation.RandomStrategy:
		config.Strategy = StrategyRandom
	case rotation.WeightedStrategy:
		config.Strategy = StrategyWeighted
	case rotation.LeastUsedStrategy:
		config.Strategy = StrategyLeastUsed
	}

	// Other settings
	config.TestURL = options.TestURL
	config.ValidateOnStart = options.ValidateOnStart
	config.MaxRetries = options.MaxRetries

	// Timeouts
	config.Timeouts.Request = options.RequestTimeout.String()
	config.Timeouts.Validation = options.ValidationTimeout.String()
	config.Timeouts.Retry = options.RetryDelay.String()

	return config
}
