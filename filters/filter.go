package filters

import (
	"encoding/json"
	"encoding/xml"
	"github.com/betacraft/yaag/middleware"
	"github.com/betacraft/yaag/yaag"
	"github.com/revel/revel"
	"log"
	"net/http/httptest"
	"strings"
)

func FilterForApiDoc(c *revel.Controller, fc []revel.Filter) {
	if record, _ := revel.Config.Bool("yaag.record"); !record {
		fc[0](c, fc[1:])
		return
	}

	w := httptest.NewRecorder()
	c.Response = revel.NewResponse(w)
	httpVerb := c.Request.Method
	customParams := make(map[string]interface{})
	headers := make(map[string]string)
	hasJson := false
	hasXml := false

	body := middleware.ReadBody(c.Request.Request)

	if c.Request.ContentType == "application/json" {
		if httpVerb == "POST" || httpVerb == "PUT" || httpVerb == "PATCH" {
			err := json.Unmarshal([]byte(*body), &customParams)
			if err != nil {
				log.Println("Json Error ! ", err)
			} else {
				hasJson = true
			}
		} else {
			err := json.Unmarshal([]byte(c.Request.URL.RawQuery), &customParams)
			if err != nil {
				log.Println("Json Error ! ", err)
			} else {
				hasJson = true
			}
		}

	} else if c.Request.ContentType == "application/xml" {
		if httpVerb == "POST" || httpVerb == "PUT" || httpVerb == "PATCH" {
			err := xml.Unmarshal([]byte(*body), &customParams)
			if err != nil {
				log.Println("Xml Error ! ", err)
			} else {
				hasXml = true
			}
		} else {
			err := xml.Unmarshal([]byte(c.Request.URL.RawQuery), &customParams)
			if err != nil {
				log.Println("Json Error ! ", err)
			} else {
				hasXml = true
			}
		}
	}
	log.Println(hasJson,hasXml)
	// call remaiing filters
	fc[0](c, fc[1:])

	c.Result.Apply(c.Request, c.Response)
	htmlValues := yaag.APICall{}
	htmlValues.CommonRequestHeaders = make(map[string]string)
	// get headers
	for k, v := range c.Request.Header {
		isCommon := false
		for _, hk := range yaag.CommonHeaders {
			if k == hk {
				isCommon = true
				htmlValues.CommonRequestHeaders[k] = strings.Join(v, " ")
				break
			}
		}
		if !isCommon {
			headers[k] = strings.Join(v, " ")
		}
	}

	htmlValues.MethodType = httpVerb
	htmlValues.CurrentPath = c.Request.URL.Path
	htmlValues.PostForm = make(map[string]string)
	for k, v := range c.Params.Form {
		htmlValues.PostForm[k] = strings.Join(v, " ")
	}
	htmlValues.RequestBody = *body
	htmlValues.RequestHeader = headers
	htmlValues.RequestUrlParams = make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		htmlValues.RequestUrlParams[k] = strings.Join(v, " ")
	}
	htmlValues.ResponseHeader = make(map[string]string)
	htmlValues.ResponseBody = w.Body.String()
	for k, v := range w.Header() {
		htmlValues.ResponseHeader[k] = strings.Join(v, " ")
	}
	htmlValues.ResponseCode = w.Code

	yaag.ApiCallValueInstance.BaseLink = c.Request.Host

	go yaag.GenerateHtml(&htmlValues)
}
