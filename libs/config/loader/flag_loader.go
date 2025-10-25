package loader

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// FlagLoader loads configuration from command-line flags.
// Uses pflag library (POSIX-compliant, compatible with standard flag package).
type FlagLoader struct {
	flagSet *pflag.FlagSet
}

// NewFlagLoader creates a new FlagLoader.
// If flagSet is nil, uses pflag.CommandLine (global default).
//
// Example:
//
//	// Using global flag set
//	pflag.String("server.host", "localhost", "Server host")
//	pflag.Int("server.port", 8080, "Server port")
//	pflag.Parse()
//	loader := loader.NewFlagLoader(nil)
//
//	// Using custom flag set
//	flags := pflag.NewFlagSet("app", pflag.ExitOnError)
//	flags.String("server.host", "localhost", "Server host")
//	flags.Parse(os.Args[1:])
//	loader := loader.NewFlagLoader(flags)
func NewFlagLoader(flagSet *pflag.FlagSet) *FlagLoader {
	if flagSet == nil {
		flagSet = pflag.CommandLine
	}
	return &FlagLoader{
		flagSet: flagSet,
	}
}

// Load binds flags and unmarshals them into dst.
//
// Note:
//   - Flags must be parsed (call flagSet.Parse()) before calling Load()
//   - Flag names with dots (.) create nested structures
//     Example: --server.port=8080 -> struct{Server: {Port: 8080}}
func (f *FlagLoader) Load(dst interface{}) error {
	v := viper.New()

	if err := v.BindPFlags(f.flagSet); err != nil {
		return fmt.Errorf("failed to bind flags: %w", err)
	}

	if err := v.Unmarshal(dst); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
