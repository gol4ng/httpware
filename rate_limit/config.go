package rate_limit

import (
	"net/http"
)

type Config struct {
	ErrorCallback ErrorCallback
}

func (c *Config) apply(options ...Option) *Config {
	for _, option := range options {
		option(c)
	}
	return c
}

func GetConfig(options ...Option) *Config {
	config := &Config{
		ErrorCallback: nopErrorCallback(),
	}
	return config.apply(options...)
}

func nopErrorCallback() ErrorCallback {
	return func(_ *http.Request, err error) (*http.Response, error) {
		return nil, err
	}
}

type Option func(*Config)

func WithErrorCallback(callback ErrorCallback) Option {
	return func(config *Config) {
		config.ErrorCallback = callback
	}
}
