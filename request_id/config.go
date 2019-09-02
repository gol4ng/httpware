package request_id

import (
	"net/http"
)

const HeaderName = "Request-Id"

type Config struct {
	HeaderName  string
	IdGenerator func(*http.Request) string
}

func NewConfig() *Config {
	return &Config{
		HeaderName:  HeaderName,
		IdGenerator: RandomIdGenerator,
	}
}
