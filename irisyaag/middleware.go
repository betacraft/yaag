package irisyaag

import (
	"strings"

	"github.com/betacraft/yaag/middleware"
	"github.com/betacraft/yaag/yaag"
	"github.com/betacraft/yaag/yaag/models"
	"github.com/kataras/iris/v12"
)

// New returns a new yaag iris-compatible handler which is responsible to generate the rest API.
func New() iris.Handler {
	return func(ctx iris.Context) {
		if !yaag.IsOn() {
			// execute the main handler and exit if yaag is off.
			ctx.Next()
			return
		}

		// prepare the middleware.
		apiCall := &models.ApiCall{}
		middleware.Before(apiCall, ctx.Request())

		w := ctx.Recorder() // starts recorder, if not already started and returns the writer.
		ctx.Next()

		if code := ctx.GetStatusCode(); yaag.IsStatusCodeValid(code) {
			apiCall.MethodType = ctx.Method()
			apiCall.CurrentPath = strings.Split(ctx.Request().RequestURI, "?")[0]
			apiCall.ResponseBody = string(w.Body()[0:])
			apiCall.ResponseCode = code

			headers := make(map[string]string, len(w.Header()))
			for k, v := range w.Header() {
				headers[k] = strings.Join(v, " ")
			}
			apiCall.ResponseHeader = headers

			go yaag.GenerateHtml(apiCall)
		}
	}
}
