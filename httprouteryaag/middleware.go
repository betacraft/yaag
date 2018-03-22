package httprouteryaag

import (
	"net/http"

	"github.com/betacraft/yaag/middleware"
	"github.com/betacraft/yaag/yaag"
	"github.com/betacraft/yaag/yaag/models"

	"github.com/julienschmidt/httprouter"
)

func HandleFunc(next func(http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if !yaag.IsOn() {
			next(w, r, ps)
			return
		}
		apiCall := models.ApiCall{}
		writer := middleware.NewResponseRecorder(w)
		middleware.Before(&apiCall, r)
		next(writer, r, ps)
		middleware.After(&apiCall, writer, r)
	})
}
