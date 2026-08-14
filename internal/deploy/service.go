package deploy

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func HandleDeploy() {
	// reads from cli
	// gets the path of the directory to deploy
	// handles flags
	// sends port name, app name, and health endpoint to deploy
	// outputs "deploy succeeded, healthy" in cli if done with work.
}

func Deploy() error {
	// gets port, app, name (appconfig)
	// deploys it to desired runtime (currently only docker)
	// returns with either or error or success. this is the part of the codebase which actually touches the commands.
	// the handler is the plumbing around it.
}

// the deep module
func Deploy(appConfig AppConfig) error {
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

	return nil
}

type AppConfig struct {
	Name   string `yaml:"name"`
	Port   int    `yaml:"port"`
	Health string `yaml:"health"`
}
