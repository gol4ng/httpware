package request_listener_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gol4ng/httpware/v4/request_listener"
	"github.com/stretchr/testify/assert"
)

func TestCurlLogDumper(t *testing.T) {
	b := bytes.NewBuffer([]byte{})
	log.SetFlags(0)
	log.SetOutput(b)

	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", ioutil.NopCloser(strings.NewReader(url.Values{
		"mykey": {"myvalue"},
	}.Encode())))

	request_listener.CurlLogDumper(req)
	assert.Equal(t, "curl -X 'GET' 'http://fake-addr' -d 'mykey=myvalue'\n", b.String())
}

func TestGetCurlCommand(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", ioutil.NopCloser(strings.NewReader(url.Values{
		"mykey": {"myvalue"},
	}.Encode())))

	cmd, err := request_listener.GetCurlCommand(req)
	assert.NoError(t, err)
	assert.Len(t, *cmd, 6)
	assert.Equal(t, "curl -X 'GET' 'http://fake-addr' -d 'mykey=myvalue'", cmd.String())
}
