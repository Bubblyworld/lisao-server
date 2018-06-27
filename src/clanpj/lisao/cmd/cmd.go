// Package cmd provides utilities for working with shell commands in go.
package cmd

import (
	"os/exec"
	"strings"
)

func NewCommand(cmd string) *exec.Cmd {
	cmdWords := strings.Split(cmd, " ")

	return exec.Command(cmdWords[0], cmdWords[1:]...)
}

func Env(name, value string) string {
	return name + "=" + value
}
