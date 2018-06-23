// Package cmd provides utilities for working with shell commands in go.
package cmd

import (
	"io"
	"os/exec"
	"strings"
)

// Command holds data about a shell command that should be run. Commands can be
// run multiple times safely, and should be built up with the fluent API.
type Command struct {
	cmd  string
	dir  string
	args []string

	env   map[string]string
	logTo io.Writer
}

func NewCommand(cmd string) Command {
	cmdWords := strings.Split(cmd, " ")

	return Command{
		cmd:  cmdWords[0],
		args: cmdWords[1:],

		env: make(map[string]string),
	}
}

func (c Command) SetEnv(key, value string) Command {
	c.env[key] = value
	return c
}

func (c Command) SetParam(name, value string) Command {
	name = strings.Trim(name, "-")
	c.args = append(c.args, "-"+name, value)
	return c
}

func (c Command) WithFlag(flag string) Command {
	flag = strings.Trim(flag, "-")
	c.args = append(c.args, "-"+flag)
	return c
}

func (c Command) WithArg(arg string) Command {
	c.args = append(c.args, arg)
	return c
}

func (c Command) LogTo(w io.Writer) Command {
	c.logTo = w
	return c
}

func (c Command) CD(dir string) Command {
	c.dir = dir
	return c
}

func (c Command) Do() error {
	cmd := exec.Command(c.cmd, c.args...)
	cmd.Dir = c.dir

	for k, v := range c.env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	if c.logTo != nil {
		stdOut, err := cmd.StdoutPipe()
		if err != nil {
			return err
		}

		stdErr, err := cmd.StderrPipe()
		if err != nil {
			return err
		}

		go io.Copy(c.logTo, stdOut)
		go io.Copy(c.logTo, stdErr)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
