package loader

import (
	"fmt"

	"github.com/spf13/viper"
)

// FileLoader loads configuration from files.
// Supported formats: JSON, YAML, TOML, Properties, HCL
type FileLoader struct {
	filePath string
	fileType string
}

// NewFileLoader creates a new FileLoader.
//
// Parameters:
//   - path: path to config file
//   - fileType: file type (json, yaml, toml, properties, hcl)
//
// Example:
//
//	loader := loader.NewFileLoader("config.yaml", "yaml")
//	loader := loader.NewFileLoader("config.json", "json")
func NewFileLoader(path, fileType string) *FileLoader {
	return &FileLoader{
		filePath: path,
		fileType: fileType,
	}
}

// Load reads config file and unmarshals it into dst.
func (f *FileLoader) Load(dst interface{}) error {
	v := viper.New()
	v.SetConfigFile(f.filePath)
	v.SetConfigType(f.fileType)

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file %s: %w", f.filePath, err)
	}

	if err := v.Unmarshal(dst); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
