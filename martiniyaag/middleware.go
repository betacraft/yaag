package martiniyaag

import (
	"github.com/CoryARamirez/yaag/middleware"
	"github.com/CoryARamirez/yaag/yaag"
	"github.com/CoryARamirez/yaag/yaag/models"
	"github.com/go-martini/martini"
	"net/http"
)

func Document(c martini.Context, w http.ResponseWriter, r *http.Request) {
	if !yaag.IsOn() {
		c.Next()
		return
	}
	apiCall := models.ApiCall{}
	writer := middleware.NewResponseRecorder(w)
	c.MapTo(writer, (*http.ResponseWriter)(nil))
	middleware.Before(&apiCall, r)
	c.Next()
	middleware.After(&apiCall, writer,  r)
}
