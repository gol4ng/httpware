package correlation_id

import (
	"net/http"
)

const HeaderName = "Correlation-Id"

type Config struct {
	HeaderName  string
	IdGenerator func(*http.Request) string
}

func NewConfig() *Config {
	return &Config{
		HeaderName: HeaderName,
		IdGenerator: func(_ *http.Request) string {
			return DefaultIdGenerator.Generate(10)
		},
	}
}
