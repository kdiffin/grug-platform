package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: grug <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "status":
		fmt.Println("everything is okay")
	case "deploy":
		fmt.Println("deploying this")
	case "help":
		fmt.Println("usage: grug <command>")
	}

	fmt.Print(os.Args[0])
	fmt.Print(os.Args[1])

	log.Printf("http server started on port %v", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("failed to start http server")
		return
	}
}
