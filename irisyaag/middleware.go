package irisyaag

import (
	"github.com/betacraft/yaag/middleware"
	"github.com/betacraft/yaag/yaag"
	"github.com/betacraft/yaag/yaag/models"

	"github.com/kataras/iris/context" // after go 1.9, users can use iris package directly.
	"net/http/httptest"
	"bytes"
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

		//iris recorder is not httptest.ResponseRecorder! So need to map it.
		r := httptest.NewRecorder()
		r.Body = bytes.NewBuffer(ctx.Recorder().Body())
		r.HeaderMap = ctx.Recorder().Header()
		r.Code = ctx.Recorder().StatusCode()

		//iris recorder writes the recorded data to its original response recorder. So pass the testrecorder
		// as responsewriter to after call.
		middleware.After(apiCall, r, r , ctx.Request())
	}
}
