package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello dud")
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/healthz", health)

	http.ListenAndServe(":8081", nil)
}
