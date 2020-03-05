package auth_test

import (
	"net/http"
	"testing"

	"github.com/gol4ng/httpware/v2/auth"
	"github.com/stretchr/testify/assert"
)

func TestFromHeader(t *testing.T) {
	req := &http.Request{
		Header: http.Header{
			"Authorization": []string{"foo"},
		},
	}

	credProvider := auth.FromHeader(req)
	cred := string(credProvider())
	assert.Equal(t, "foo", cred)
}


func TestAddHeader(t *testing.T) {
	req := &http.Request{
		Header: make(http.Header),
	}

	credSetter := auth.AddHeader(req)
	credSetter("foo")
	assert.Equal(t, "foo", req.Header.Get("Authorization"))
}
