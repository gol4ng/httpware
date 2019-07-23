package metrics

import "net/http"

type Config struct {
	Recorder               Recorder
	SplitStatus            bool
	DisableMeasureSize     bool
	DisableMeasureInflight bool
	HandlerNameExtractor   func(req *http.Request) string
}

func NewConfig(recorder Recorder) *Config {
	return &Config{
		Recorder:               recorder,
		SplitStatus:            false,
		DisableMeasureSize:     false,
		DisableMeasureInflight: false,
		HandlerNameExtractor: func(req *http.Request) string {
			return req.URL.String()
		},
	}
}
