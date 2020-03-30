package rate_limit

import (
	"net/http"
)

type Config struct {
	ErrorCallback func(err error, request *http.Request) error
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

func nopErrorCallback() func(error, *http.Request) error {
	return func(err error, _ *http.Request) error {
		return err
	}
}

type Option func(*Config)

func WithErrorCallback(callback func(error, *http.Request) error) Option {
	return func(config *Config) {
		config.ErrorCallback = callback
	}
}
