package irisyaag

import (
	"strings"

	"github.com/betacraft/yaag/middleware"
	"github.com/betacraft/yaag/yaag"
	"github.com/betacraft/yaag/yaag/models"

	"github.com/kataras/iris/context" // after go 1.9, users can use iris package directly.
)

// New returns a new yaag iris-compatible handler which is responsible to generate the rest API.
func New() context.Handler {
	return func(ctx context.Context) {
		if !yaag.IsOn() {
			// execute the main handler and exit if yaag is off.
			ctx.Next()
			return
		}

		// prepare the middleware.
		apiCall := &models.ApiCall{}
		middleware.Before(apiCall, ctx.Request())

		// start the recorder instead of raw response writer,
		// response writer is changed for that handler now.
		ctx.Record()
		// and then fire the "main" handler.
		ctx.Next()

		if statusCode := ctx.GetStatusCode(); statusCode != 404 {
			apiCall.MethodType = ctx.Method()
			apiCall.CurrentPath = strings.Split(ctx.Request().RequestURI, "?")[0]
			apiCall.ResponseCode = statusCode
			apiCall.RequestUrlParams = ctx.URLParams()
			apiCall.ResponseBody = string(ctx.Recorder().Body())
			// copy resp headers.
			headers := map[string]string{}
			for k, v := range ctx.ResponseWriter().Header() {
				headers[k] = strings.Join(v, " ")
			}
			apiCall.ResponseHeader = headers
			// different goroutine for generating html document.
			go yaag.GenerateHtml(apiCall)
		}
	}
}
