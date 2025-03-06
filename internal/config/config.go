package config

import (
	"time"

	"github.com/greysquirr3l/lashes/internal/storage"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Proxy    ProxyConfig    `yaml:"proxy"`
	Captcha  CaptchaConfig  `yaml:"captcha"`
}

type DatabaseConfig struct {
	Type           storage.DatabaseType `yaml:"type"`
	DSN            string               `yaml:"dsn"`
	MaxConnections int                  `yaml:"max_connections"`
	RetentionDays  int                  `yaml:"retention_days"`
	MetricsEnabled bool                 `yaml:"metrics_enabled"`
	ConnTimeout    time.Duration        `yaml:"conn_timeout"`
	QueryTimeout   time.Duration        `yaml:"query_timeout"`
}

type ProxyConfig struct {
	RotationStrategy string        `yaml:"rotation_strategy"`
	ValidateOnStart  bool          `yaml:"validate_on_start"`
	RefreshInterval  time.Duration `yaml:"refresh_interval"`
	MaxRetries       int           `yaml:"max_retries"`
	Timeout          time.Duration `yaml:"timeout"`
}

type CaptchaConfig struct {
	TwoCaptchaAPIKey string        `yaml:"2captcha_api_key" env:"TWOCAPTCHA_API_KEY"`
	Timeout          time.Duration `yaml:"timeout"`
	Debug            bool          `yaml:"debug"`
}

func GetDefaultConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Type:           storage.SQLite,
			MaxConnections: 10,
			RetentionDays:  30,
			MetricsEnabled: true,
			ConnTimeout:    time.Minute,
			QueryTimeout:   time.Second * 30,
		},
		Proxy: ProxyConfig{
			RotationStrategy: "round-robin",
			ValidateOnStart:  true,
			RefreshInterval:  time.Hour,
			MaxRetries:       3,
			Timeout:          time.Second * 10,
		},
		Captcha: CaptchaConfig{
			Timeout: time.Minute * 2,
			Debug:   false,
		},
	}
}
