package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/kdiffin/grug-platform/internal/deploy"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: grug <command>")
	}

	switch args[0] {
	case "status":
		fmt.Println("everything is okay")
	// todo: break this up into it's own function
	case "deploy":
		deploy.DeployCommand(args[1:])
		deploy.Deploy()
	default:
		return errors.New("this command does not exist in the grugverse: " + args[0])
	}

	return nil
}
