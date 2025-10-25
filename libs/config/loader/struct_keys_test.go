package loader

import (
	"reflect"
	"sort"
	"testing"
)

func TestExtractStructKeys_Simple(t *testing.T) {
	type SimpleConfig struct {
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
	}

	keys := ExtractKeysFromType(SimpleConfig{})
	sort.Strings(keys)

	expected := []string{"name", "port"}
	sort.Strings(expected)

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

func TestExtractStructKeys_Nested(t *testing.T) {
	type NestedConfig struct {
		Server struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"server"`
		Database struct {
			Host string `mapstructure:"host"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"database"`
	}

	keys := ExtractKeysFromType(NestedConfig{})
	sort.Strings(keys)

	expected := []string{
		"database.host",
		"database.port",
		"server.host",
		"server.port",
	}
	sort.Strings(expected)

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

func TestExtractStructKeys_DeepNested(t *testing.T) {
	type DeepConfig struct {
		App struct {
			Server struct {
				HTTP struct {
					Host string `mapstructure:"host"`
					Port int    `mapstructure:"port"`
				} `mapstructure:"http"`
			} `mapstructure:"server"`
		} `mapstructure:"app"`
	}

	keys := ExtractKeysFromType(DeepConfig{})
	sort.Strings(keys)

	expected := []string{
		"app.server.http.host",
		"app.server.http.port",
	}
	sort.Strings(expected)

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

func TestExtractStructKeys_NoTags(t *testing.T) {
	type NoTagsConfig struct {
		Name string
		Port int
	}

	keys := ExtractKeysFromType(NoTagsConfig{})
	sort.Strings(keys)

	// Should use lowercase field names
	expected := []string{"name", "port"}
	sort.Strings(expected)

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

func TestExtractStructKeys_SkipDash(t *testing.T) {
	type SkipConfig struct {
		Name     string `mapstructure:"name"`
		Internal string `mapstructure:"-"`
		Port     int    `mapstructure:"port"`
	}

	keys := ExtractKeysFromType(SkipConfig{})
	sort.Strings(keys)

	// Should skip field with "-" tag
	expected := []string{"name", "port"}
	sort.Strings(expected)

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

func TestExtractStructKeys_Pointer(t *testing.T) {
	type PointerConfig struct {
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
	}

	// Test with pointer type
	keys := ExtractKeysFromType(&PointerConfig{})
	sort.Strings(keys)

	expected := []string{"name", "port"}
	sort.Strings(expected)

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

func TestExtractStructKeys_MixedTypes(t *testing.T) {
	type MixedConfig struct {
		StringVal string  `mapstructure:"string_val"`
		IntVal    int     `mapstructure:"int_val"`
		BoolVal   bool    `mapstructure:"bool_val"`
		FloatVal  float64 `mapstructure:"float_val"`
	}

	keys := ExtractKeysFromType(MixedConfig{})
	sort.Strings(keys)

	expected := []string{
		"bool_val",
		"float_val",
		"int_val",
		"string_val",
	}
	sort.Strings(expected)

	if !reflect.DeepEqual(keys, expected) {
		t.Errorf("Expected %v, got %v", expected, keys)
	}
}

func TestExtractStructKeys_EmptyStruct(t *testing.T) {
	type EmptyConfig struct{}

	keys := ExtractKeysFromType(EmptyConfig{})

	if len(keys) != 0 {
		t.Errorf("Expected empty keys, got %v", keys)
	}
}
