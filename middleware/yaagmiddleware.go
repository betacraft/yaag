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
	"net/url"
	"strings"
	"yaag/yaag"
)

var reqWriteExcludeHeaderDump = map[string]bool{
	"Host":              true, // not in Header map anyway
	"Content-Length":    true,
	"Transfer-Encoding": true,
	"Trailer":           true,
}

type YaagHandler struct {
	next func(http.ResponseWriter, *http.Request)
}

func Handle(next func(http.ResponseWriter, *http.Request)) http.Handler {
	return &YaagHandler{next: next}
}

func (y *YaagHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writer := httptest.NewRecorder()
	apiCall := yaag.APICall{}
	before(&apiCall, r)
	y.next(writer, r)
	after(&apiCall, writer, w, r)
}

func HandleFunc(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		apiCall := yaag.APICall{}
		writer := httptest.NewRecorder()
		before(&apiCall, r)
		next(writer, r)
		after(&apiCall, writer, w, r)
	}
}

func before(apiCall *yaag.APICall, req *http.Request) {
	headers := readHeaders(req)
	val, ok := headers["Content-Type"]
	log.Println(val)
	if ok {
		switch strings.TrimSpace(headers["Content-Type"]) {
		case "application/x-www-form-urlencoded":
			fallthrough
		case "application/json, application/x-www-form-urlencoded":
			log.Println("Reading form")
			readPostForm(req)
		case "application/json":
			log.Println("Reading body")
			readBody(req)
		}
	}
}

func readQueryParams(req *http.Request) map[string]string {
	params := map[string]string{}
	u, err := url.Parse(req.RequestURI)
	if err != nil {
		return params
	}
	for _, param := range strings.Split(u.Query().Encode(), "&") {
		value := strings.Split(param, "=")
		params[value[0]] = value[1]
	}
	return params
}

func printMap(m map[string]string) {
	for key, value := range m {
		log.Println(key, "=>", value)
	}
}

func readPostForm(req *http.Request) map[string]string {
	postForm := map[string]string{}
	log.Println("", *readBody(req))
	for _, param := range strings.Split(*readBody(req), "&") {
		value := strings.Split(param, "=")
		postForm[value[0]] = value[1]
	}
	return postForm
}

func readHeaders(req *http.Request) map[string]string {
	b := bytes.NewBuffer([]byte(""))
	err := req.Header.WriteSubset(b, reqWriteExcludeHeaderDump)
	if err != nil {
		return map[string]string{}
	}
	headers := map[string]string{}
	for _, header := range strings.Split(b.String(), "\n") {
		values := strings.Split(header, ":")
		if strings.EqualFold(values[0], "") {
			continue
		}
		headers[values[0]] = values[1]
	}
	//printMap(headers)
	return headers
}

func ReadBody(req *http.Request) *string {
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

func after(apiCall *yaag.APICall, writer *httptest.ResponseRecorder, w http.ResponseWriter, r *http.Request) {
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
