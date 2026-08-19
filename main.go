package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
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
	// todo: break this up into it's own function
	case "deploy":
		deployCmd := flag.NewFlagSet("deploy", flag.ContinueOnError)

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
		path := ""
		paths := deployCmd.Args()
		if len(paths) >= 1 {
			path = paths[0]
		}

		appConfig, err := loadConfig("./test-api-deployment/")
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}
		portStr := strconv.Itoa(appConfig.Port)

		buildCmd := exec.Command("docker", "build", "-t", "grug/"+appConfig.Name, path)
		if err := buildCmd.Run(); err != nil {
			return fmt.Errorf("build dockerfile: %w", err)
		}

		cmd := exec.Command("docker", "run", "--name", appConfig.Name, "-p", portStr+":"+"8081", "-d", "grug/"+appConfig.Name)

		log.Print(cmd.String())

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf(
				"docker failed: %w: %s",
				err,
				strings.TrimSpace(string(output)),
			)
		}
	default:
		return errors.New("this command does not exist in the grugverse: " + args[0])
	}

	return nil
}

// TODO: add validation to the grug.yaml file
func loadConfig(appDir string) (*AppConfig, error) {
	fmt.Println("printing appDir", appDir)
	if appDir[len(appDir)-1] == '/' {
		appDir = removeRuneAtIndex(len(appDir)-1, appDir)
	}

	fmt.Println("printing appDir", appDir)

	var appConfig AppConfig
	bytes, err := os.ReadFile(appDir + "/grug.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading %s/grug.yaml: %w", appDir, err)
	}

	if err := yaml.Unmarshal(bytes, &appConfig); err != nil {
		fmt.Println(err)
	}
	return &appConfig, nil
}

// TODO: remove this functiona nd replace its usage with filepath lib
// https://pkg.go.dev/path/filepath
func removeRuneAtIndex(i int, str string) string {
	runes := []rune(str)

	// Correct order: Keep everything BEFORE index + everything AFTER index
	s := string(runes[:i]) + string(runes[i+1:])
	return s
}
