package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func livez(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "ok")
	if err != nil {
		fmt.Print("cannot write to file")
	}
}

func main() {
	// load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("couldn't load environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/livez", livez)
	// todo
	http.HandleFunc("/readyz", livez)

	log.Printf("http server started on port %v", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("failed to start http server")
		return
	}
}
