package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	http.HandleFunc("/", home(templates))

	fmt.Println("starting server on :8090")
	http.ListenAndServe(":8090", nil)
}
