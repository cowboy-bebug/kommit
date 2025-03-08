package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

func execCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute: %v", err)
	}

	return out.String(), nil
}

func ExecGit(args ...string) (string, error) {
	return execCmd("git", args...)
}
