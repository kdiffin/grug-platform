package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("starting server on :8090")
	http.ListenAndServe(":8090", nil)
}
