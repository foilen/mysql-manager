package main

import (
	"os/exec"
)

func getCommandOutput(name string, arg ...string) (string, error) {

	cmd := exec.Command(name, arg...)

	outputBytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(outputBytes), nil

}
