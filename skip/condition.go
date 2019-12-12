package skip

import (
	"net/http"
)

type Condition func(*http.Request) bool
