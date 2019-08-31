package request_id

import (
	"net/http"
)

const HeaderName = "Request-Id"

type Config struct {
	HeaderName string
	GenerateId func(*http.Request) string
}

func NewConfig() *Config {
	return &Config{
		HeaderName: HeaderName,
		GenerateId: func(request *http.Request) string {
			return RandomString(10) // eg: XPF0G5kqEG
		},
	}
}
