package core

// Loader defines a generic interface for loading configuration from various sources.
type Loader[T any] interface {
	// Load reads config from source and fills the provided data structure.
	Load(T) error
}
