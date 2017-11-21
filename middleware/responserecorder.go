package middleware

import (
	"net/http"
	"bytes"
	"net"
	"bufio"
	"errors"
)

//go:generate easytags $GOFILE

type responseRecorder struct {
	writer http.ResponseWriter
	Status int
	Body   *bytes.Buffer
}

func NewResponseRecorder(w http.ResponseWriter) *responseRecorder {
	r := &responseRecorder{
		writer:w,
		Status:http.StatusOK,
		Body:bytes.NewBuffer(nil),
	}
	return r
}

func (r *responseRecorder) Header() http.Header {
	return r.writer.Header()
}

func (r *responseRecorder) WriteHeader(status int) {
	r.Status = status
	r.writer.WriteHeader(status)
}

func (r *responseRecorder) Write(buf []byte) (int, error) {
	n, err := r.writer.Write(buf)
	if err == nil {
		r.Body.Write(buf)
	}
	return n, err
}

func (r *responseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := r.writer.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("Error in hijacker")
}