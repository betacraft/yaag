# YAAG : Yet Another API doc Generator

Golang is awesome for developing web apps. And people have created a bunch of awesome Web-Frameworks, Web helper libraries. If we consider the entire ecosystem of web apps in Golang everything except API documentation seems to be in place. So we have created the first API doc generator for Golang based web apps and calling it Yet another.

## Why ?

Most of the web services expose their APIs to the mobile or third party developers. Documenting them is somewhat pain in the ass. We are trying to reduce the pain, atleast for in house projects where you don't have to expose your documentation to the world. YAAG generates simple bootstrap based API documentation without writing a single line of comments.

## How it works ?

YAAG is a middleware. You have to add YAAG handler in your routes and you are done. Just go on calling your APIs using POSTMAN, Curl or from any client, and YAAG will keep on updating the API Doc html. 


## How to use with basic net.http package

1. Import github.com/gophergala/yaag/yaag
2. Import github.com/gophergala/yaag/middleware
3. Initialize yaag ```yaag.Init(&yaag.Config{On: true, DocTitle: "Core", DocPath: "apidoc.html"})```
4. Use it in your handlers as ```http.HandleFunc("/", middleware.HandleFunc(handler))```

####Sample code

```go
func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
  yaag.Init(&yaag.Config{On: true, DocTitle: "Core", DocPath: "apidoc.html"})
  http.HandleFunc("/", middleware.HandleFunc(handler))
  http.ListenAndServe(":8080", nil)
}
```

## How to use with Gorilla Mux
1. Import github.com/gophergala/yaag/yaag
2. Import github.com/gophergala/yaag/middleware
3. Initialize yaag ```yaag.Init(&yaag.Config{On: true, DocTitle: "Core", DocPath: "apidoc.html"})```
4. Use it in your handlers as ```r.HandleFunc("/", middleware.HandleFunc(handler))```

####Sample code

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

1. Import github.com/gophergala/yaag/yaag
2. Import github.com/gophergala/yaag/martiniyaag
3. Initialize yaag ```yaag.Init(&yaag.Config{On: true, DocTitle: "Core", DocPath: "apidoc.html"})```
4. Add Yaag middleware like ```m.Use(martiniyaag.Document)```

####Sample Code

```go
func main() {
  yaag.Init(&yaag.Config{On: true, DocTitle: "Gorilla Mux", DocPath: "apidoc.html"})
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
2. import github.com/rainingclouds/yaag/filters in app/init.go
3. add 'filters.FilterForApiDoc' in revel.Filters
4. Start recording Api calls


## Screenshots

#### API doc is generated based on the paths
![alt first](https://raw.github.com/gophergala/yaag/master/1.png)
#### Click on any call to see the details of the API
![alt second](https://raw.github.com/gophergala/yaag/master/2.png)

## Screencast

[YAAG ScreenCast](https://www.youtube.com/watch?v=dQWXxJn6_iE&feature=youtu.be)

## Adding Support for

1. Gin framework

## Team

* Aniket Awati (aniket@rainingclouds.com)
* Akshay Deo (akshay@rainingclouds.com)
* Kaustubh Deshmukh (kaustubh@rainingclouds.com)

This project is initiated by RainingClouds Inc during GopherGala 2015.
