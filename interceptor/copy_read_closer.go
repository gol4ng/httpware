package interceptor

import (
	"bytes"
	"io"
	"io/ioutil"
)

// io.Reader with Read method reset offset when EOF
type bufReader struct {
	buf []byte
	off int
}

func (r *bufReader) Read(p []byte) (n int, err error) {
	if r.off == len(r.buf) {
		if len(p) == 0 {
			return 0, nil
		}
		r.off = 0
		return 0, io.EOF
	}

	n = copy(p, r.buf[r.off:])
	r.off += n

	return n, nil
}

type copyReadCloser struct {
	io.ReadCloser
	// write in bytes.Buffer
	copyTemp *bytes.Buffer
	// read in copy
	copy *bufReader
}

// First read with io.TeeReader
//      -> copyBuffered
//    /
// src --> output
// Second read after EOF
// copyBuffered --> copy BufReader simple buffer with fix size
// when BufReader is EOF offset is reset to read again
func NewCopyReadCloser(src io.ReadCloser) *copyReadCloser {
	buf := &bytes.Buffer{}
	tr := &copyReadCloser{
		copyTemp: buf,
	}

	tr.ReadCloser = &struct {
		io.Reader
		io.Closer
	}{io.TeeReader(src, buf), src}

	return tr
}

func (tr *copyReadCloser)Read(p []byte) (n int, err error) {
	n, err = tr.ReadCloser.Read(p)
	if err == io.EOF {
		if tr.copy == nil {
			tr.ReadCloser.Close()
			tr.copy = &bufReader{buf: tr.copyTemp.Bytes()}
			tr.copyTemp.Reset()
			tr.ReadCloser = ioutil.NopCloser(tr.copy)
		}
	}

	return n, err
}
