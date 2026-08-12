package main

import (
	"fmt"
	"log"
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

	log.Println("starting server at 8081")
	http.ListenAndServe(":8081", nil)
}
