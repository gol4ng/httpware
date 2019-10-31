package metrics

import "net/http"

type Config struct {
	Recorder Recorder
	// func that allows you to provide a strategy to identify/group metrics
	// you can group metrics by request host/url/... or app name ...
	// by default, we group metrics by request url
	IdentifierProvider func(req *http.Request) string
	// if set to true, each response status will be represented by a metrics
	// if set to false, response status codes will be grouped by first digit (204/201/200/... -> 2xx; 404/403/... -> 4xx)
	SplitStatus bool
	// if set to true, recorder will add a responseSize metric
	ObserveResponseSize bool
	// if set to true, recorder will add a metric representing the number of inflight requests
	MeasureInflightRequests bool
}

func (c *Config) apply(options ...Option) *Config {
	for _, option := range options {
		option(c)
	}
	return c
}

// NewConfig returns a new metrics configuration with all options applied
func NewConfig(recorder Recorder, options ...Option) *Config {
	config := &Config{
		Recorder:                recorder,
		SplitStatus:             false,
		ObserveResponseSize:     true,
		MeasureInflightRequests: true,
		IdentifierProvider: func(req *http.Request) string {
			return req.URL.String()
		},
	}
	return config.apply(options...)
}

// Option defines a metrics middleware/tripperware configuration option
type Option func(*Config)

// WithSplitStatus will configure SplitStatus metrics option
func WithSplitStatus(splitStatus bool) Option {
	return func(config *Config) {
		config.SplitStatus = splitStatus
	}
}

// WithObserveResponseSize will configure ObserveResponseSize metrics option
func WithObserveResponseSize(observeResponseSize bool) Option {
	return func(config *Config) {
		config.ObserveResponseSize = observeResponseSize
	}
}

// WithMeasureInflightRequests will configure MeasureInflightRequests metrics option
func WithMeasureInflightRequests(measureInflightRequests bool) Option {
	return func(config *Config) {
		config.MeasureInflightRequests = measureInflightRequests
	}
}

// WithIdentifierProvider will configure IdentifierProvider metrics option
func WithIdentifierProvider(identifierProvider func(req *http.Request) string) Option {
	return func(config *Config) {
		config.IdentifierProvider = identifierProvider
	}
}
