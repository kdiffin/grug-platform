package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: grug <command>")
		os.Exit(1)
		return
	}

	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
		return
	}
}

func run(args []string) error {
	switch args[0] {
	case "status":
		fmt.Println("everything is okay")
	case "deploy":
		fmt.Println("deploying this")
	case "help":
		fmt.Println("usage: grug <command>")
	default:
		return errors.New("this command does not exist in the grugverse: " + args[0])
	}

	return nil
}
