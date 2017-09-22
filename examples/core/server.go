package main

import (
	"fmt"
	"github.com/betacraft/yaag/middleware"
	"github.com/betacraft/yaag/yaag"
	"io/ioutil"
	"net/http"
	"strings"
	"path/filepath"
)

/*
@Description This is a dummy handler
@Path /
 */
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

/*
@Description This is a dummy post handler
@Path /say_it
 */
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
	c := &yaag.Config{On: true,
		DocTitle:         "Core",
		DocPath:          "apidoc.html",
		BaseUrls: map[string]string{"Production": "http://testing.com",
		}}
	p , _ := filepath.Abs(c.DocPath)
	c.RootPath = filepath.Dir(p)
	fmt.Println(c.RootPath)
	yaag.Init(c)
	http.HandleFunc("/", middleware.HandleFunc(handler))
	http.HandleFunc("/say_it", middleware.HandleFunc(postHandler))
	http.ListenAndServe(":8080", nil)
}
