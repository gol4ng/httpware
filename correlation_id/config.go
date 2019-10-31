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

// GetConfig returns a new correlation configuration with all options applied
func GetConfig(options ...Option) *Config {
	config := &Config{
		HeaderName: HeaderName,
		IdGenerator: func(_ *http.Request) string {
			return DefaultIdGenerator.Generate(10)
		},
	}
	return config.apply(options...)
}

// Option defines a correlationId middleware/tripperware configuration option
type Option func(*Config)

// WithHeaderName will configure HeaderName correlation option
func WithHeaderName(headerName string) Option {
	return func(config *Config) {
		config.HeaderName = headerName
	}
}

// WithIdGenerator will configure IdGenerator correlation option
func WithIdGenerator(idGenerator func(*http.Request) string) Option {
	return func(config *Config) {
		config.IdGenerator = idGenerator
	}
}
