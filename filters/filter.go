package filters

import (
	"encoding/xml"
	"github.com/revel/revel"
	"log"
	"net/http/httptest"
)

type 

func FilterForApiDoc(c *revel.Controller, fc []revel.Filter) {
	w := httptest.NewRecorder()
	c.Response = revel.NewResponse(w)	
	fc[0](c, fc[1:])
	httpVerb := c.Request.Method
	params := c.Params
	var customParams map[string]interface{}
	hasJson := false
	hasXml := false
	if c.Request.ContentType == "application/json" {
		if httpVerb == "POST" || httpVerb == "PUT" || httpVerb == "PATCH" {			
				err := json.Unmarshal(c.Request.Body, customParams)
				if err != nil {
					hasJson = true
				}			
		} else {
			//TODO check if query params are json encoded
		}

	} else if c.Request.ContentType == "application/xml" {
		if httpVerb == "POST" || httpVerb == "PUT" || httpVerb == "PATCH" {
			err := xml.NewDecoder(c.Request.Body).Decode(customParams)
			if err != nil {
				hasXml = true
			}
		} else {
			//TODO check if query params are json encoded
		}
	}

	// call yaag generator
	c.Result.Apply(c.Request, c.Response)
	
}
