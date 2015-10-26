/*
 * This is the main core of the yaag package
 */
package yaag

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
)

const TEMPLATE = `<!DOCTYPE html>
<html>
<head lang="en">
    <title> API Documentation </title>
    <!-- Latest compiled and minified CSS -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css"
          integrity="sha512-dTfge/zgoMYpP7QbHy4gWMEGsbsdZeCXz7irItjcC3sPUFtf0kuFbDz/ixG7ArTxmDjLXDmezHubeNikyKGVyQ=="
          crossorigin="anonymous">

    <!-- Optional theme -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap-theme.min.css"
          integrity="sha384-aUGj/X2zp5rLCbBxumKTCw2Z50WgIr1vs/PFN4praOTvYXWlVyh2UtNUU0KAUhAX" crossorigin="anonymous">

    <!-- Latest compiled and minified JavaScript -->
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/js/bootstrap.min.js"
            integrity="sha512-K1qjQ+NcF2TYO/eI3M6v8EiNYZfA95pQumfvcVrTHtwQVDG+aHRqLi/ETn2uB+1JqwYqVG3LIvdm9lj6imS/pQ=="
            crossorigin="anonymous"></script>

    <script src="http://google-code-prettify.googlecode.com/svn/loader/run_prettify.js"></script>
    <link href='http://fonts.googleapis.com/css?family=Roboto' rel='stylesheet' type='text/css'>
    <link rel="stylesheet" href="http://google-code-prettify.googlecode.com/svn/trunk/src/prettify.css">
    <link rel="stylesheet" href="style.css">
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
    <!-- Latest compiled and minified JavaScript -->
    <style type="text/css">
        body {
            font-family: 'Roboto', sans-serif;
        }
    </style>
    <style type="text/css">
        pre.prettyprint {
            border: 1px solid #ccc;
            margin-bottom: 0;
            padding: 9.5px;
        }
    </style>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.4/highlight.min.js"></script>
    <script>hljs.initHighlightingOnLoad();</script>
</head>
<body>
<nav class="navbar navbar-default navbar-static-top">
    <div class="container-fluid">
        <!-- Brand and toggle get grouped for better mobile display -->
        <div class="navbar-header">
            <button type="button" class="navbar-toggle collapsed" data-toggle="collapse"
                    data-target="#bs-example-navbar-collapse-1">
                <span class="sr-only">Toggle navigation</span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
                <span class="icon-bar"></span>
            </button>
            <a class="navbar-brand" href="#">{{.Title}}</a>
        </div>
        <!-- /.navbar-collapse -->
    </div>
    <!-- /.container-fluid -->
</nav>
<div class="container-fluid">
    <div class="alert alert-default">
        <p>Base URL : <strong>{{.BaseLink}}</strong></p>
    </div>
    <div class="wrapper">
        <section class="col--cwcl--l">
            <div class="nav--cwcl">
                <h2 class="title--column js-nav-title">Paths</h2>
                {{ range $key, $value := .array }}
                <ul class="list--reset list--column js-nav-list">
                    <li><a class="nav--cwcl__item anchor">{{$value.Path}}</a> </li>
                </ul>

                {{ end }}
            </div>
        </section>
        <section class="col-cwcl-r float--clear">
            <div class="wrapper--content">

            </div>
        </section>

    </div>
</div>
</body>
</html>`

var CommonHeaders = []string{
	"Accept",
	"Accept-Encoding",
	"Accept-Language",
	"Cache-Control",
	"Content-Length",
	"Content-Type",
	"Origin",
	"User-Agent",
	"X-Forwarded-For",
}

var count int

type APICall struct {
	Id int

	CurrentPath string
	MethodType  string

	PostForm map[string]string

	RequestHeader        map[string]string
	CommonRequestHeaders map[string]string
	ResponseHeader       map[string]string
	RequestUrlParams     map[string]string

	RequestBody  string
	ResponseBody string
	ResponseCode int
}

type PathSpec struct {
	HttpVerb   string
	Path       string
	HtmlValues []APICall
}

type ApiCallValue struct {
	BaseLink string
	Path     []PathSpec
}

type Config struct {
	On       bool
	DocTitle string
	DocPath  string
}

var ApiCallValueInstance = &ApiCallValue{}
var config *Config = &Config{On: false, DocPath: "apidoc.html", DocTitle: "YAAG"}

func IsOn() bool {
	return config.On
}

func Init(conf *Config) {
	filePath, err := filepath.Abs(conf.DocPath + ".json")
	dataFile, err := os.Open(filePath)
	defer dataFile.Close()

	if err == nil {
		json.NewDecoder(io.Reader(dataFile)).Decode(ApiCallValueInstance)
	}
	config = conf
}

func GenerateHtml(htmlValue *APICall) {
	shouldAddPathSpec := true
	log.Printf("PathSpec : %v", ApiCallValueInstance.Path)
	for k, pathSpec := range ApiCallValueInstance.Path {
		if pathSpec.Path == htmlValue.CurrentPath && pathSpec.HttpVerb == htmlValue.MethodType {
			shouldAddPathSpec = false
			shouldAdd := true
			if shouldAdd {
				htmlValue.Id = count
				count += 1
				deleteCommonHeaders(htmlValue)
				ApiCallValueInstance.Path[k].HtmlValues = append(pathSpec.HtmlValues, *htmlValue)
			}
		}
	}

	if shouldAddPathSpec {
		pathSpec := PathSpec{
			HttpVerb: htmlValue.MethodType,
			Path:     htmlValue.CurrentPath,
		}
		htmlValue.Id = count
		count += 1
		deleteCommonHeaders(htmlValue)
		pathSpec.HtmlValues = append(pathSpec.HtmlValues, *htmlValue)
		ApiCallValueInstance.Path = append(ApiCallValueInstance.Path, pathSpec)
	}

	t := template.New("API Documentation")
	filePath, err := filepath.Abs(config.DocPath)
	htmlString := TEMPLATE

	t, err = t.Parse(htmlString)
	if err != nil {
		log.Println(err)
		return
	}
	homeHtmlFile, err := os.Create(filePath)
	defer homeHtmlFile.Close()

	dataFile, err := os.Create(filePath + ".json")
	if err != nil {
		log.Println(err)
		return
	}
	defer dataFile.Close()

	data, err := json.Marshal(ApiCallValueInstance)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = dataFile.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
	homeWriter := io.Writer(homeHtmlFile)
	t.Execute(homeWriter, map[string]interface{}{"array": ApiCallValueInstance.Path,
		"BaseLink": ApiCallValueInstance.BaseLink, "Title": config.DocTitle})
}

func deleteCommonHeaders(call *APICall) {
	delete(call.RequestHeader, "Accept")
	delete(call.RequestHeader, "Accept-Encoding")
	delete(call.RequestHeader, "Accept-Language")
	delete(call.RequestHeader, "Cache-Control")
	delete(call.RequestHeader, "Connection")
	delete(call.RequestHeader, "Cookie")
	delete(call.RequestHeader, "Origin")
	delete(call.RequestHeader, "User-Agent")
}
