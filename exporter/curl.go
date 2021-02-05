package exporter

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func CurlFmtPrinter(cmd *Cmd, err error) {
	if err == nil {
		fmt.Println(cmd)
		return
	}
	fmt.Println("cannot print curl command", err)
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
