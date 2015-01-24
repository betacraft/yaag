/*
 * This is yaag middleware for the web apps using the middlewares that supports http handleFunc
 */
package middleware

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"strings"
)

type YaagHandler struct {
	next func(http.ResponseWriter, *http.Request)
}

func Handle(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return &YaagHandler{next: next}
}

func (y *YaagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writer := httptest.NewRecorder()
	before(r)
	y.next(writer, r)
	after(writer, w, r)
}

func HandleFunc(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		writer := httptest.NewRecorder()
		before(r)
		next(writer, r)
		after(writer, w, r)
	}
}

func before(req *http.Request) {
	log.Println(*readBody(req))
}

func readBody(req *http.Request) *string {
	save := req.Body
	var err error
	if req.Body == nil {
		req.Body = nil
	} else {
		save, req.Body, err = drainBody(req.Body)
		if err != nil {
			return nil
		}
	}
	b := bytes.NewBuffer([]byte(""))
	chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"
	if req.Body == nil {
		return nil
	}

	var dest io.Writer = b
	if chunked {
		dest = httputil.NewChunkedWriter(dest)
	}
	_, err = io.Copy(dest, req.Body)
	if chunked {
		dest.(io.Closer).Close()
	}
	req.Body = save
	body := b.String()
	return &body
}

func after(writer *httptest.ResponseRecorder, w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.RequestURI, ".ico") {
		fmt.Fprintf(w, writer.Body.String())
		return
	}
	log.Println(r.RequestURI)
	log.Println(writer.Body.String())
	log.Println(writer.Code)
	for header := range writer.Header() {
		log.Println(header)
	}
	w.WriteHeader(writer.Code)
	w.Write(writer.Body.Bytes())
}

// One of the copies, say from b to r2, could be avoided by using a more
// elaborate trick where the other copy is made during Request/Response.Write.
// This would complicate things too much, given that these functions are for
// debugging only.
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
