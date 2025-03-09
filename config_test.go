package lashes

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/greysquirr3l/lashes/internal/rotation"
	"github.com/greysquirr3l/lashes/internal/storage"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	configJSON := `{
		"storage": {
			"type": "sqlite",
			"file_path": "test.db",
			"query_timeout": "5s"
		},
		"strategy": "weighted",
		"test_url": "https://test.example.com",
		"validate_on_start": true,
		"timeouts": {
			"request": "10s",
			"validation": "3s",
			"retry": "500ms"
		},
		"max_retries": 5,
		"circuit_breaker": {
			"enabled": true,
			"max_failures": 3,
			"reset_timeout": "20s",
			"enable_global_breaker": true
		}
	}`

	if err := os.WriteFile(configPath, []byte(configJSON), 0o644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Load the config
	opts, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Check values were loaded correctly
	if opts.Storage == nil {
		t.Fatal("Storage options should not be nil")
	}

	if opts.Storage.Type != storage.SQLite {
		t.Errorf("Storage.Type = %s, want %s", opts.Storage.Type, storage.SQLite)
	}

	if opts.Storage.FilePath != "test.db" {
		t.Errorf("Storage.FilePath = %s, want %s", opts.Storage.FilePath, "test.db")
	}

	if opts.Storage.QueryTimeout != time.Second*5 {
		t.Errorf("Storage.QueryTimeout = %v, want %v", opts.Storage.QueryTimeout, time.Second*5)
	}

	if opts.Strategy != rotation.WeightedStrategy {
		t.Errorf("Strategy = %s, want %s", opts.Strategy, rotation.WeightedStrategy)
	}

	if opts.TestURL != "https://test.example.com" {
		t.Errorf("TestURL = %s, want %s", opts.TestURL, "https://test.example.com")
	}

	if !opts.ValidateOnStart {
		t.Error("ValidateOnStart = false, want true")
	}

	if opts.RequestTimeout != time.Second*10 {
		t.Errorf("RequestTimeout = %v, want %v", opts.RequestTimeout, time.Second*10)
	}

	if opts.ValidationTimeout != time.Second*3 {
		t.Errorf("ValidationTimeout = %v, want %v", opts.ValidationTimeout, time.Second*3)
	}

	if opts.RetryDelay != time.Millisecond*500 {
		t.Errorf("RetryDelay = %v, want %v", opts.RetryDelay, time.Millisecond*500)
	}

	if opts.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want %d", opts.MaxRetries, 5)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Save original env values to restore later
	origStorageType := os.Getenv("LASHES_STORAGE_TYPE")
	origSQLitePath := os.Getenv("LASHES_SQLITE_PATH")
	origStrategy := os.Getenv("LASHES_STRATEGY")
	origTestURL := os.Getenv("LASHES_TEST_URL")
	origTimeout := os.Getenv("LASHES_REQUEST_TIMEOUT")

	// Restore env after test
	defer func() {
		os.Setenv("LASHES_STORAGE_TYPE", origStorageType)
		os.Setenv("LASHES_SQLITE_PATH", origSQLitePath)
		os.Setenv("LASHES_STRATEGY", origStrategy)
		os.Setenv("LASHES_TEST_URL", origTestURL)
		os.Setenv("LASHES_REQUEST_TIMEOUT", origTimeout)
	}()

	// Set test env variables
	os.Setenv("LASHES_STORAGE_TYPE", "sqlite")
	os.Setenv("LASHES_SQLITE_PATH", "env_test.db")
	os.Setenv("LASHES_STRATEGY", "random")
	os.Setenv("LASHES_TEST_URL", "https://env-test.example.com")
	os.Setenv("LASHES_REQUEST_TIMEOUT", "15s")

	// Load config from env
	opts := LoadConfigFromEnv()

	// Check values were loaded from env
	if opts.Storage == nil {
		t.Fatal("Storage options should not be nil")
	}

	if opts.Storage.Type != storage.SQLite {
		t.Errorf("Storage.Type = %s, want %s", opts.Storage.Type, storage.SQLite)
	}

	if opts.Storage.FilePath != "env_test.db" {
		t.Errorf("Storage.FilePath = %s, want %s", opts.Storage.FilePath, "env_test.db")
	}

	if opts.Strategy != rotation.RandomStrategy {
		t.Errorf("Strategy = %s, want %s", opts.Strategy, rotation.RandomStrategy)
	}

	if opts.TestURL != "https://env-test.example.com" {
		t.Errorf("TestURL = %s, want %s", opts.TestURL, "https://env-test.example.com")
	}

	if opts.RequestTimeout != time.Second*15 {
		t.Errorf("RequestTimeout = %v, want %v", opts.RequestTimeout, time.Second*15)
	}
}

func TestSaveConfig(t *testing.T) {
	// Create test options
	opts := Options{
		Storage: &storage.Options{
			Type:             storage.MySQL,
			ConnectionString: "user:pass@tcp(localhost:3306)/dbname",
			QueryTimeout:     10 * time.Second,
		},
		Strategy:          rotation.LeastUsedStrategy,
		TestURL:           "https://save-test.example.com",
		ValidateOnStart:   true,
		MaxRetries:        4,
		RequestTimeout:    8 * time.Second,
		ValidationTimeout: 2 * time.Second,
		RetryDelay:        750 * time.Millisecond,
	}

	// Create a temporary file for saving
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "saved_config.json")

	// Save the config
	if err := SaveConfig(opts, configPath); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Load the config to verify it saved correctly
	loadedOpts, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed after save: %v", err)
	}

	// Check values were preserved
	if loadedOpts.Storage == nil {
		t.Fatal("Loaded Storage options should not be nil")
	}

	if loadedOpts.Storage.Type != storage.MySQL {
		t.Errorf("Loaded Storage.Type = %s, want %s", loadedOpts.Storage.Type, storage.MySQL)
	}

	if loadedOpts.Storage.ConnectionString != "user:pass@tcp(localhost:3306)/dbname" {
		t.Errorf("Loaded Storage.ConnectionString doesn't match original")
	}

	if loadedOpts.Strategy != rotation.LeastUsedStrategy {
		t.Errorf("Loaded Strategy = %s, want %s", loadedOpts.Strategy, rotation.LeastUsedStrategy)
	}

	if loadedOpts.TestURL != "https://save-test.example.com" {
		t.Errorf("Loaded TestURL = %s, want %s", loadedOpts.TestURL, "https://save-test.example.com")
	}
}
