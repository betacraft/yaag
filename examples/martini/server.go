package main

import (
	"github.com/go-martini/martini"
	"yaag/martiniyaag"
	"yaag/yaag"
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
