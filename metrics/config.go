package metrics

import "net/http"

type Config struct {
	Recorder                Recorder
	// func that allows you to provide a strategy to identify/group metrics
	// you can group metrics by request host/url/... or app name ...
	// by default, we group metrics by request url
	IdentifierProvider      func(req *http.Request) string
	// if set to true, each response status will be represented by a metrics
	// if set to false, response status codes will be grouped by first digit (204/201/200/... -> 2xx; 404/403/... -> 4xx)
	SplitStatus             bool
	// if set to true, recorder will add a responseSize metric
	ObserveResponseSize     bool
	// if set to true, recorder will add a metric representing the number of inflight requests
	MeasureInflightRequests bool
}

func NewConfig(recorder Recorder) *Config {
	return &Config{
		Recorder:                recorder,
		SplitStatus:             false,
		ObserveResponseSize:     true,
		MeasureInflightRequests: true,
		IdentifierProvider: func(req *http.Request) string {
			return req.URL.String()
		},
	}
}
