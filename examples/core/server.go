package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"yaag/middleware"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.EqualFold(r.Method, "POST") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Illegal request"))
		return
	}
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, "Yaag is awesome", string(body))
}

func main() {
	http.HandleFunc("/", middleware.HandleFunc(handler))
	http.HandleFunc("/say_it", middleware.HandleFunc(postHandler))
	http.ListenAndServe(":8080", nil)
}
