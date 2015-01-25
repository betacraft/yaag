package main

import (
	"github.com/go-martini/martini"
	"yaag/middleware"
	"yaag/yaag"
)

func main() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Run()
}
