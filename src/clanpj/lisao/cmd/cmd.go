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

	stdOutPipeReader *io.PipeReader
	stdOutPipeWriter *io.PipeWriter

	stdInPipeReader *io.PipeReader
	stdInPipeWriter *io.PipeWriter
}

func NewCommand(cmd string) Command {
	cmdWords := strings.Split(cmd, " ")
	stdOutPR, stdOutPW := io.Pipe()
	stdInPR, stdInPW := io.Pipe()

	return Command{
		cmd:  cmdWords[0],
		args: cmdWords[1:],

		env: make(map[string]string),

		stdOutPipeReader: stdOutPR,
		stdOutPipeWriter: stdOutPW,

		stdInPipeReader: stdInPR,
		stdInPipeWriter: stdInPW,
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

func (c Command) GetStdOut() *io.PipeReader {
	return c.stdOutPipeReader
}

func (c Command) GetStdIn() *io.PipeWriter {
	return c.stdInPipeWriter
}

func (c Command) Do() error {
	defer func() {
		// TODO(guy) handle errors here
		c.stdInPipeWriter.Close()
		c.stdOutPipeReader.Close()
	}()

	cmd := exec.Command(c.cmd, c.args...)
	cmd.Dir = c.dir

	for k, v := range c.env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdIn, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdOutSrc := stdOut.(io.Reader)
	if c.logTo != nil {
		go io.Copy(c.logTo, stdErr)

		stdOutSrc = io.TeeReader(stdOutSrc, c.logTo)
	}

	go io.Copy(stdIn, c.stdInPipeReader)
	go io.Copy(c.stdOutPipeWriter, stdOutSrc)

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}
