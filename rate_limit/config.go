package rate_limit

import (
	"net/http"
)

type Config struct {
	ErrorCallback func(err error, req *http.Request) error
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

func nopErrorCallback() func(err error, req *http.Request) error {
	return func(err error, req *http.Request) error {
		return err
	}
}

type Option func(*Config)

func WithErrorCallback(callback func(err error, req *http.Request) error) Option {
	return func(config *Config) {
		config.ErrorCallback = callback
	}
}
