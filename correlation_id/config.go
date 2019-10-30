package correlation_id

import (
	"net/http"
)

const HeaderName = "Correlation-Id"

type Config struct {
	HeaderName  string
	IdGenerator func(*http.Request) string
}

func (c *Config) apply(options ...Option) *Config {
	for _, option := range options {
		option(c)
	}
	return c
}

// GetConfig return a new correlation configuration with all options applied
func GetConfig(options ...Option) *Config {
	config := &Config{
		HeaderName: HeaderName,
		IdGenerator: func(_ *http.Request) string {
			return DefaultIdGenerator.Generate(10)
		},
	}
	return config.apply(options...)
}

// Option was correlation middleware/tripperware configurable options
type Option func(*Config)

// WithHeaderName will configure HeaderName correlation options
func WithHeaderName(headerName string) Option {
	return func(config *Config) {
		config.HeaderName = headerName
	}
}

// WithIdGenerator will configure IdGenerator correlation options
func WithIdGenerator(idGenerator func(*http.Request) string) Option {
	return func(config *Config) {
		config.IdGenerator = idGenerator
	}
}
