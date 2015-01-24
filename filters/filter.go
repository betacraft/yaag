package filters

import (
	"encoding/json"
	"encoding/xml"
	"github.com/revel/revel"
	"log"
	"net/http/httptest"
	"strings"
	"yaag/middleware"
	"yaag/yaag"
)

func FilterForApiDoc(c *revel.Controller, fc []revel.Filter) {

	if record, _ := revel.Config.Bool("yaag.record"); !record {
		log.Printf("record %v ", record)
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
	log.Println(*body)

	if c.Request.ContentType == "application/json" {
		if httpVerb == "POST" || httpVerb == "PUT" || httpVerb == "PATCH" {
			err := json.Unmarshal([]byte(*body), &customParams)
			if err != nil {
				log.Println("Json Error ! ", err)
			} else {
				hasJson = true
			}
		} else {
			//TODO check if query params are json encoded
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
			//TODO check if query params are json encoded
		}
	}
	// call remaiing filters
	fc[0](c, fc[1:])

	c.Result.Apply(c.Request, c.Response)

	// get headers
	for k, v := range c.Request.Header {
		headers[k] = strings.Join(v, " ")
	}

	log.Println("Params:")
	if hasJson {
		log.Printf("%#v", customParams)
	} else if hasXml {
		log.Printf("%#v", customParams)
	}
	log.Printf("Standard Params %#v", c.Params)

	log.Printf("\nurl path %s", c.Request.URL.Path)

	log.Println("Headers")
	log.Printf("%#v", headers)

	log.Printf("\n Status %v Response %s", w.Code, w.Body.String())

	htmlValues := yaag.HtmlValueContainer{}

	htmlValues.BaseLink = c.Request.URL.Host
	htmlValues.MethodType = httpVerb
	htmlValues.CurrentPath = c.Request.URL.Path
	for k, v := range c.Params.Form {
		htmlValues.PostForm[k] = v
	}
	htmlValues.RequestBody = *body
	htmlValues.RequestHeader = headers
	for k, v := range c.Request.URL.Query() {
		htmlValues.RequestUrlParams[k] = strings.Join(v, " ")
	}
	htmlValues.ResponseBody = w.Body.String()
	for k, v := range w.Header() {
		htmlValues.ResponseHeader[k] = strings.Join(v, " ")
	}
	htmlValues.ResponseCode = w.Code
	yaag.GenerateHtml(&htmlValues)
}
