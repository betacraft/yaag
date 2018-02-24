package filters

import (
	"encoding/json"
	"encoding/xml"
	"github.com/CoryARamirez/yaag/middleware"
	"github.com/CoryARamirez/yaag/yaag"
	"github.com/CoryARamirez/yaag/yaag/models"
	"github.com/revel/revel"
	"log"
	"net/http/httptest"
	"strings"
	"net/url"
	"net/http"
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
	r := Route(c.Request.Request)
	path := c.Request.URL.Path
	if r != nil {
		path = r.Path
	}
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
	log.Println(hasJson, hasXml)
	// call remaiing filters
	fc[0](c, fc[1:])

	c.Result.Apply(c.Request, c.Response)
	if !yaag.IsStatusCodeValid(c.Response.Status) {
		return
	}
	htmlValues := models.ApiCall{}
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
	htmlValues.CurrentPath = path
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
	go yaag.GenerateHtml(&htmlValues)
}

func Route(req *http.Request) (route *revel.Route) {
	// Override method if set in header
	if method := req.Header.Get("X-HTTP-Method-Override"); method != "" && req.Method == "POST" {
		req.Method = method
	}
	treePath := func(method, path string) string {
		if method == "*" {
			method = ":METHOD"
		}
		return "/" + method + path
	}
	leaf, expansions := revel.MainRouter.Tree.Find(treePath(req.Method, req.URL.Path))
	if leaf == nil {
		return nil
	}

	// Create a map of the route parameters.
	var params url.Values
	if len(expansions) > 0 {
		params = make(url.Values)
		for i, v := range expansions {
			params[leaf.Wildcards[i]] = []string{v}
		}
	}
	var controllerName, methodName string

	// The leaf value is now a list of possible routes to match, only a controller
	routeList := leaf.Value.([]*revel.Route)

	//INFO.Printf("Found route for path %s %#v", req.URL.Path, len(routeList))
	for index := range routeList {
		route = routeList[index]
		methodName = route.MethodName

		// Special handling for explicit 404's.
		if route.Action == "404" {
			route = nil
			break
		}

		// If wildcard match on method name use the method name from the params
		if methodName[0] == ':' {
			methodName = strings.ToLower(params[methodName[1:]][0])
		}

		// If the action is variablized, replace into it with the captured args.
		controllerName = route.ControllerName
		if controllerName[0] == ':' {
			controllerName = strings.ToLower(params[controllerName[1:]][0])
			if route.ModuleSource.ControllerByName(controllerName, methodName) != nil {
				break
			}
		} else {
			break
		}
		route = nil
	}
	return
}
