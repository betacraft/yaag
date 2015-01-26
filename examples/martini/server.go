package main

import (
	"github.com/go-martini/martini"
	"github.com/gophergala/yaag/martiniyaag"
	"github.com/gophergala/yaag/yaag"
)

func main() {
	yaag.Init(&yaag.Config{On: true, DocTitle: "Gorilla Mux", DocPath: "apidoc.html"})
	m := martini.Classic()
	m.Use(martiniyaag.Document)
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Run()
}
