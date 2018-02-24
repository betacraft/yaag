[![Build Status](https://travis-ci.org/CoryARamirez/yaag.svg)](https://travis-ci.org/CoryARamirez/yaag)

[Trello Board](https://trello.com/b/jCZlTsNj/yaag)

# YAAG : Yet Another API doc Generator

Golang is awesome for developing web apps. And people have created a bunch of awesome Web-Frameworks, Web helper libraries. If we consider the entire ecosystem of web apps in Golang everything except API documentation seems to be in place. So we have created the first API doc generator for Golang based web apps and calling it Yet another.

## Why ?

Most of the web services expose their APIs to the mobile or third party developers. Documenting them is somewhat pain in the ass. We are trying to reduce the pain, at least for in house projects where you don't have to expose your documentation to the world. YAAG generates simple bootstrap based API documentation without writing a single line of comments.

## How it works ?

YAAG is a middleware. You have to add YAAG handler in your routes and you are done. Just go on calling your APIs using POSTMAN, Curl or from any client, and YAAG will keep on updating the API Doc html. (Note: We are also generating a json file containing data af all API calls)

## How to use with basic net.http package

1. Import github.com/CoryARamirez/yaag/yaag
2. Import github.com/CoryARamirez/yaag/middleware
3. Initialize yaag ```yaag.Init(&yaag.Config{On: true, DocTitle: "Core", DocPath: "apidoc.html"})```
4. Use it in your handlers as ```http.HandleFunc("/", middleware.HandleFunc(handler))```

#### Sample code

```go
func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
  yaag.Init(&yaag.Config{On: true, DocTitle: "Core", DocPath: "apidoc.html", BaseUrls : map[string]string{"Production":"","Staging":""} })
  http.HandleFunc("/", middleware.HandleFunc(handler))
  http.ListenAndServe(":8080", nil)
}
```

## How to use with Gorilla Mux
1. Import github.com/CoryARamirez/yaag/yaag
2. Import github.com/CoryARamirez/yaag/middleware
3. Initialize yaag ```yaag.Init(&yaag.Config{On: true, DocTitle: "Gorilla Mux", DocPath: "apidoc.html"})```
4. Use it in your handlers as ```r.HandleFunc("/", middleware.HandleFunc(handler))```

#### Sample code

```go
func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, time.Now().String())
}

func main() {
  yaag.Init(&yaag.Config{On: true, DocTitle: "Gorilla Mux", DocPath: "apidoc.html"})
  r := mux.NewRouter()
  r.HandleFunc("/", middleware.HandleFunc(handler)) 
  http.ListenAndServe(":8080", r)
}
```

## How to use with Martini

1. Import github.com/CoryARamirez/yaag/yaag
2. Import github.com/CoryARamirez/yaag/martiniyaag
3. Initialize yaag ```yaag.Init(&yaag.Config{On: true, DocTitle: "Martini", DocPath: "apidoc.html"})```
4. Add Yaag middleware like ```m.Use(martiniyaag.Document)```

#### Sample Code

```go
func main() {
  yaag.Init(&yaag.Config{On: true, DocTitle: "Martini", DocPath: "apidoc.html"})
  m := martini.Classic()
  m.Use(martiniyaag.Document)
  m.Get("/", func() string {
    return "Hello world!"
  })
  m.Run()
}
```

## How to use with Revel

1. Add yaag.record = true in conf/app.conf (before starting to record the api calls)
2. import github.com/CoryARamirez/yaag/filters in app/init.go
3. add 'filters.FilterForApiDoc' in revel.Filters
4. Start recording Api calls

### Sample Code

```go
func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		filters.FilterForApiDoc,       // This enables yaag to record apicalls
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.OnAppStart(func() {
		yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
			On:       true,
			DocTitle: "Revel",
			DocPath:  "examples/revel/apidoc.html",
			BaseUrls: map[string]string{"Production": "", "Staging": ""},
		})
	})
}
```


## How to use with Gin

1. Import github.com/CoryARamirez/yaag/yaag
2. Import github.com/CoryARamirez/yaag/gin
3. Initialize yaag ```yaag.Init(&yaag.Config(On: true, DocTile: "Gin", DocPath: "apidpc.html"))```
4. Add yaag middleware like ```r.User(yaag_gin.Document())```

#### Sample Code

```go
import (
    "net/http"
    yaag_gin "github.com/CoryARamirez/yaag/gin/v1"
    "github.com/CoryARamirez/yaag/yaag"
    "gopkg.in/gin-gonic/gin.v1"
    )
func main() {
    r := gin.Default()
    yaag.Init(&yaag.Config{On: true, DocTitle: "Gin", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Staging": ""}})
    r.Use(yaag_gin.Document())
    // use other middlewares ...
    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello World!")
    })
    r.Run(":8080")
}
```

### Using github for gin dependency 

```go
import (
    "net/http"
    yaag_gin "github.com/CoryARamirez/yaag/gin"
    "github.com/CoryARamirez/yaag/yaag"
    "github.com/gin-gonic/gin"
    )
func main() {
    r := gin.Default()
    yaag.Init(&yaag.Config{On: true, DocTitle: "Gin", DocPath: "apidoc.html", BaseUrls: map[string]string{"Production": "", "Staging": ""}})
    r.Use(yaag_gin.Document())
    // use other middlewares ...
    r.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello World!")
    })
    r.Run(":8080")
}
```


## How to use with Iris

1. Import `github.com/CoryARamirez/yaag/yaag`
2. Import `github.com/CoryARamirez/yaag/irisyaag`
3. Initialize yaag `yaag.Init(&yaag.Config(On: true, DocTile: "Iris", DocPath: "apidoc.html"))`
4. Register yaag middleware like `app.Use(irisyaag.New())`

> `irisyaag` records the response body and provides all the necessary information to the apidoc.

### Sample Code

```go
package main

import (
  "github.com/kataras/iris"
  "github.com/kataras/iris/context"

  "github.com/CoryARamirez/yaag/irisyaag"
  "github.com/CoryARamirez/yaag/yaag"
)

type myXML struct {
  Result string `xml:"result"`
}

func main() {
  app := iris.New()

  yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
    On:       true,
    DocTitle: "Iris",
    DocPath:  "apidoc.html",
    BaseUrls: map[string]string{"Production": "", "Staging": ""},
  })
  app.Use(irisyaag.New()) // <- IMPORTANT, register the middleware.

  app.Get("/json", func(ctx context.Context) {
    ctx.JSON(context.Map{"result": "Hello World!"})
  })

  app.Get("/plain", func(ctx context.Context) {
    ctx.Text("Hello World!")
  })

  app.Get("/xml", func(ctx context.Context) {
    ctx.XML(myXML{Result: "Hello World!"})
  })

  app.Get("/complex", func(ctx context.Context) {
    value := ctx.URLParam("key")
    ctx.JSON(context.Map{"value": value})
  })

  // Run our HTTP Server.
  //
  // Note that on each incoming request the yaag will generate and update the "apidoc.html".
  // Recommentation:
  // Write tests that calls those handlers, save the generated "apidoc.html".
  // Turn off the yaag middleware when in production.
  //
  // Example usage:
  // Visit all paths and open the generated "apidoc.html" file to see the API's automated docs.
  app.Run(iris.Addr(":8080"))
}
```

## Screenshots

#### API doc is generated based on the paths
![alt first](https://raw.github.com/CoryARamirez/yaag/master/1.png)
#### Click on any call to see the details of the API
![alt second](https://raw.github.com/CoryARamirez/yaag/master/2.png)

## Screencast

[YAAG ScreenCast](https://www.youtube.com/watch?v=dQWXxJn6_iE&feature=youtu.be)

## Contributors 

* Aniket Awati (aniket@CoryARamirez.co)
* Akshay Deo (akshay@CoryARamirez.co)
* Brian Newsom (Brian.Newsom@Colorado.edu)

This project is initiated by Betacraft during GopherGala 2015.
