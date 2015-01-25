package martiniyaag

import (
	"github.com/go-martini/martini"
	"github.com/gophergala/yaag/middleware"
	"github.com/gophergala/yaag/yaag"
	"net/http"
	"net/http/httptest"
)

func Document(c martini.Context, w http.ResponseWriter, r *http.Request) {
	if !yaag.IsOn() {
		c.Next()
		return
	}
	apiCall := yaag.APICall{}
	writer := httptest.NewRecorder()
	c.MapTo(writer, (*http.ResponseWriter)(nil))
	middleware.Before(&apiCall, r)
	c.Next()
	middleware.After(&apiCall, writer, w, r)
}
