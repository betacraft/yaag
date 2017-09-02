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
		// fire the "main" handler.
		ctx.Next()

		if statusCode := ctx.GetStatusCode(); statusCode != 404 {
			apiCall.MethodType = ctx.Method()
			apiCall.CurrentPath = strings.Split(ctx.Request().RequestURI, "?")[0]
			// we could use the response recorder built'n inside Iris but better keep that empty,
			// most users don't need that, it's aligned with other middlewares as well.
			apiCall.ResponseBody = ""
			apiCall.ResponseCode = statusCode
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
