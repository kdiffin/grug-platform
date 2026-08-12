package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type AppConfig struct {
	Name   string `yaml:"name"`
	Port   int    `yaml:"port"`
	Health string `yaml:"health"`
}

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
	case "deploy":
		deployCmd := flag.NewFlagSet("deploy", flag.ContinueOnError)

		// help isnt an error usually its just regular output
		deployCmd.SetOutput(os.Stdout)
		showHelp := deployCmd.Bool("help", false, "you can point grug at a directory using the --dir flag and if it has a grug.yaml the project will be deployed")

		if err := deployCmd.Parse(args[1:]); err != nil {
			return err
		}

		if *showHelp {
			deployCmd.PrintDefaults()
			return nil
		}

		// the path of the repo we want to deploy
		paths := deployCmd.Args()
		if len(paths) >= 1 {
			path := paths[0]
			fmt.Println(path)

		}

	default:
		return errors.New("this command does not exist in the grugverse: " + args[0])
	}

	return nil
}
