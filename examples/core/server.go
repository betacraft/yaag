package main

import (
	"fmt"
	"github.com/betacraft/yaag/middleware"
	"github.com/betacraft/yaag/yaag"
	"io/ioutil"
	"net/http"
	"strings"
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
	w.WriteHeader(http.StatusOK)
	w.Header().Add("test", "tesasasdasd")
	fmt.Fprintf(w, string(body))
}

func main() {
	yaag.Init(&yaag.Config{On: true, DocTitle: "Core", DocPath: "apidoc.html"})
	http.HandleFunc("/", middleware.HandleFunc(handler))
	http.HandleFunc("/say_it", middleware.HandleFunc(postHandler))
	http.ListenAndServe(":8080", nil)
}
