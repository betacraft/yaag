/*
 * This is the main core of the yaag package
 */
package yaag

import (
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
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap.min.css">
    <link href='http://fonts.googleapis.com/css?family=Roboto' rel='stylesheet' type='text/css'>
    <!-- Optional theme -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/css/bootstrap-theme.min.css">
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
    <!-- Latest compiled and minified JavaScript -->
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.2/js/bootstrap.min.js"></script>
    <style type="text/css">
        body {
            font-family: 'Roboto', sans-serif;
        }
    </style>
</head>
<body>
<nav class="navbar navbar-default navbar-fixed-top">
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

        <!-- Collect the nav links, forms, and other content for toggling -->
        <div class="collapse navbar-collapse pull-right" id="bs-example-navbar-collapse-1">
            <form class="navbar-form navbar-left" role="search">
                <div class="form-group">
                    <input type="text" class="form-control" placeholder="Search">
                </div>
                <button type="submit" class="btn btn-default">Find</button>
            </form>
        </div>
        <!-- /.navbar-collapse -->
    </div>
    <!-- /.container-fluid -->
</nav>
<div class="container" style="margin-top: 70px;">
    <div class="alert alert-info">
    <p>Base URL => <strong>{{.BaseLink}}</strong></p></div>
    <hr>
    {{ range $wrapperKey, $wrapperValue := .array }}

    <h4 style="cursor:pointer;" type="button" data-toggle="collapse" data-target="#{{$wrapperValue.Id}}"
        aria-expanded="false" aria-controls="collapseExample"><span class="glyphicon glyphicon-link"
                                                                    aria-hidden="true"></span> <code>{{$wrapperValue.MethodType}}
        {{$wrapperValue.CurrentPath}}</code>
    </h4>

    <div class="collapse" id="{{$wrapperValue.Id}}">
        {{ if $wrapperValue.RequestHeader }}
        <p> <H4> Request Headers </H4> </p>
        <table class="table table-bordered table-striped">
            <tr>
                <th>Key</th>
                <th>Value</th>
            </tr>
            {{ range $key, $value := $wrapperValue.RequestHeader }}
            <tr>
                <td>{{ $key }}</td>
                <td> {{ $value }}</td>
            </tr>
            {{ end }}
        </table>
        {{ end }}

        {{ if $wrapperValue.PostForm }}
        <p> <H4> Post Form </H4> </p>
        {{ range $key, $value := $wrapperValue.PostForm }}
        <li><strong>{{ $key }}</strong>: {{ $value }}</li>
        {{ end }}
        {{ end }}

        {{ if $wrapperValue.RequestUrlParams }}
        <p> <H4> URL Params </H4> </p>
        {{ range $key, $value := $wrapperValue.RequestUrlParams }}
        <li><strong>{{ $key }}</strong>: {{ $value }}</li>
        {{ end }}
        {{ end }}

        {{ if $wrapperValue.RequestBody }}
        <p> <H4> Request Body </H4> </p>
        <pre> {{ $wrapperValue.RequestBody }} </pre>
        {{ end }}

        <p><h4> Response Code</h4></p>
        <pre>{{ $wrapperValue.ResponseCode }}</pre>

        {{ if $wrapperValue.ResponseHeader }}
        <p> <H4> Response Headers </H4> </p>
        {{ range $key, $value := $wrapperValue.ResponseHeader }}
        <li><strong>{{ $key }}</strong>: {{ $value }}</li>
        {{ end }}
        {{ end }}

        {{ if $wrapperValue.ResponseBody }}
        <p> <H4> Response Body </H4> </p>
        <pre> {{ $wrapperValue.ResponseBody }} </pre>
        {{ end }}
    </div>
    <hr>
    {{ end}}
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
	Init     bool
	DocTitle string
	DocPath  string
}

var ApiCallValueInstance = &ApiCallValue{}

func main() {
	//	firstApi := APICall{Id: 1, MethodType: "GET", CurrentPath: "/login/:id", RequestHeader: map[string]string{"Content-Type": "application/json", "Accept": "application/json"},

	//		RequestBody: "{ 'main' : { 'id' : 2, 'name' : 'Gopher' }}"}

	secondApi := APICall{Id: 2, MethodType: "POST", CurrentPath: "/singup", RequestHeader: map[string]string{"Content-Type": "application/json", "Accept": "application/json"},
		ResponseBody: "{ 'main' : { 'Key' : 'ABC-123-XYZ', 'name' : 'Gopher' }}"}

	config := Config{Init: false, DocPath: "html/home.html", DocTitle: "YAAG"}

	//	valueArray := []APICall{secondApi, firstApi}
	//	allApis := ApiCallValue{BaseLink: "www.google.com", HtmlValues: valueArray}
	GenerateHtml(&secondApi, &config)
}

func GenerateHtml(htmlValue *APICall, config *Config) {
	shouldAddPathSpec := true
	log.Printf("PathSpec : %v", ApiCallValueInstance.Path)
	for _, pathSpec := range ApiCallValueInstance.Path {
		if pathSpec.Path == htmlValue.CurrentPath && pathSpec.HttpVerb == htmlValue.MethodType {
			shouldAddPathSpec = false
			shouldAdd := true
			for _, value := range pathSpec.HtmlValues {
				if value.ResponseBody == htmlValue.ResponseBody {
					shouldAdd = false
				}
			}
			if shouldAdd {
				htmlValue.Id = len(pathSpec.HtmlValues) + 1
				pathSpec.HtmlValues = append(pathSpec.HtmlValues, *htmlValue)
			}
		}
	}

	if shouldAddPathSpec {
		pathSpec := PathSpec{
			HttpVerb: htmlValue.MethodType,
			Path:     htmlValue.CurrentPath,
		}
		pathSpec.HtmlValues = append(pathSpec.HtmlValues, *htmlValue)
		ApiCallValueInstance.Path = append(ApiCallValueInstance.Path, pathSpec)
	}

	t := template.New("API Documentation")
	filePath, err := filepath.Abs(config.DocPath)
	// file, err := ioutil.ReadFile("templates/main.html")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	htmlString := TEMPLATE
	t, err = t.Parse(htmlString)
	if err != nil {
		log.Println(err)
		return
	}
	homeHtmlFile, err := os.Create(filePath)
	defer homeHtmlFile.Close()

	if err != nil {
		log.Println(err)
		return
	}
	homeWriter := io.Writer(homeHtmlFile)
	t.Execute(homeWriter, map[string]interface{}{"array": ApiCallValueInstance.Path,
		"BaseLink": ApiCallValueInstance.BaseLink, "Title": config.DocTitle})
}
