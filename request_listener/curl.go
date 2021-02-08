package request_listener

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// CurlLogDumper will log the given http.Request
// it gonna log the request using cli curl command format
// Caution request header and Body can be heavy!!
// Pay attention to the dump quantity/weight
func CurlLogDumper(request *http.Request) {
	if request == nil {
		return
	}
	cmd, err := GetCurlCommand(request)
	if err != nil {
		log.Println("cannot print curl command", err)
		return
	}
	log.Println(cmd)
}

func GetCurlCommand(req *http.Request) (*Cmd, error) {
	cmd := &Cmd{
		"curl",
		"-X", escape(req.Method),
		escape(req.URL.String()),
	}

	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = nopCloser{bytes.NewBuffer(body)}
		if len(string(body)) > 0 {
			cmd.append("-d", escape(string(body)))
		}
	}

	for h := range req.Header {
		cmd.append("-H", escape(h+": "+strings.Join(req.Header[h], " ")))
	}

	return cmd, nil
}

type Cmd []string

func (c *Cmd) append(newSlice ...string) {
	*c = append(*c, newSlice...)
}

func (c *Cmd) String() string {
	return strings.Join(*c, " ")
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func escape(str string) string {
	return `'` + strings.Replace(str, `'`, `'\''`, -1) + `'`
}
