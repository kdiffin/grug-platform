package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"go.yaml.in/yaml/v3"
)

// the thin interface
// deeper module implemented in its own service folder
// take it as a given that deploy was the *2nd arg* in the command sent by the user to the CLI.
func Deploy(args []string) error {
	deployCmd := flag.NewFlagSet("deploy", flag.ContinueOnError)

	path := ""
	paths := deployCmd.Args()
	if len(paths) >= 1 {
		path = paths[0]
	}

	bytes, err := os.ReadFile(path + "/grug.yaml")
	if err != nil {
		return fmt.Errorf("error reading %q/grug.yaml: %w", path, err)
	}

	appConfig := AppConfig{}

	if err := yaml.Unmarshal(bytes, &appConfig); err != nil {
		return fmt.Errorf("error unmarshalling yaml: %w", err)
	}

	// help isnt an error usually its just regular output
	deployCmd.SetOutput(os.Stdout)
	showHelp := deployCmd.Bool("help", false, "you can point grug at a directory using the --dir flag and if it has a grug.yaml the project will be deployed")

	if err := deployCmd.Parse(args[1:]); err != nil {
		return err
	}

	// todo: the defaults here fucking suck add a structured cli output here for each of your commands
	if *showHelp || (len(deployCmd.Args()) < 1) {
		deployCmd.PrintDefaults()
		return nil
	}

	// the path of the repo we want to deploy
	portStr := strconv.Itoa(appConfig.Port)

	return nil
}
