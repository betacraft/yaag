// file: main_test.go
package main

import (
	"testing"
	"time"

	"github.com/kataras/iris/httptest"
)

func TestIrisYaag(t *testing.T) {
	app := newApp()
	e := httptest.New(t, app)

	e.GET("/json").Expect().Status(httptest.StatusOK).
		JSON().Equal(map[string]interface{}{"result": "Hello World!"})

	e.POST("/hello").WithFormField("username", "kataras").Expect().Status(httptest.StatusOK).
		Body().Equal("Hello kataras")

	e.POST("/reqbody").WithJSON(myModel{Username: "kataras"}).Expect().Status(httptest.StatusOK).
		Body().Equal("kataras")

	// give time to "yaag" to generate the doc, 5 seconds are more than enough
	time.Sleep(5 * time.Second)
}
