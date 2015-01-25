package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"yaag/middleware"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("test", "tesasasdasd")
	fmt.Fprintf(w, string(body))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", middleware.HandleFunc(handler))
	r.HandleFunc("/testing", middleware.HandleFunc(postHandler)).Methods("POST")
	http.ListenAndServe(":8080", r)
}
