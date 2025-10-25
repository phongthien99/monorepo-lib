package loader

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// EnvLoader loads configuration from environment variables.
// Example: APP_SERVER_HOST will be converted to server.host
type EnvLoader struct {
	prefix string
	keys   []string // Optional: specific keys to bind
}

// NewEnvLoader creates a new EnvLoader with the given prefix.
// If prefix is "APP", it will read env vars like APP_*.
// Pass empty string "" if no prefix is needed.
func NewEnvLoader(prefix string) *EnvLoader {
	return &EnvLoader{
		prefix: prefix,
	}
}

// WithKeys specifies which keys to bind from environment.
// By default, EnvLoader will bind all env vars.
// Use WithKeys to bind only specific keys.
//
// Example:
//
//	loader := loader.NewEnvLoader("APP").
//	    WithKeys("server.host", "server.port", "database.password")
func (e *EnvLoader) WithKeys(keys ...string) *EnvLoader {
	e.keys = keys
	return e
}

// WithAutoKeys automatically extracts all keys from a struct type using reflection.
// This is more convenient than manually listing all keys.
//
// Example:
//
//	type AppConfig struct {
//	    Server struct {
//	        Host string `mapstructure:"host"`
//	        Port int    `mapstructure:"port"`
//	    } `mapstructure:"server"`
//	}
//
//	loader := loader.NewEnvLoader("APP").WithAutoKeys(AppConfig{})
func (e *EnvLoader) WithAutoKeys(example interface{}) *EnvLoader {
	e.keys = ExtractKeysFromType(example)
	return e
}

// Load reads environment variables and unmarshals them into dst.
//
// Conversion rules (handled automatically by Viper):
//   - Prefix is automatically uppercased: "app" -> "APP_"
//   - Underscore (_) is converted to dot (.): APP_SERVER_HOST -> server.host
//
// Example: with prefix="app", env var APP_SERVER_HOST maps to field server.host
func (e *EnvLoader) Load(dst interface{}) error {
	v := viper.New()

	if e.prefix != "" {
		v.SetEnvPrefix(e.prefix)
	}

	// Convert "." and "-" to "_" for env vars
	// Example: key "server.host" will look for env var "SERVER_HOST"
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	v.AutomaticEnv()

	// Bind specific keys if provided
	// This is necessary because AutomaticEnv() doesn't populate AllSettings()
	// but only works when Get() is called
	if len(e.keys) > 0 {
		for _, key := range e.keys {
			v.BindEnv(key)
		}
	}

	if err := v.Unmarshal(dst); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
