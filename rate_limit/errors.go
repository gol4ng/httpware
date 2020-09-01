package rate_limit

import (
	"net/http"
)

const (
	RequestLimitReachedErr = "request limit reached"
)

type ErrorCallback func(request *http.Request, err error) (*http.Response, error)
