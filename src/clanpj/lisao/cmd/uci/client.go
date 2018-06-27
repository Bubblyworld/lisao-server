package uci

import (
	"clanpj/lisao/cmd"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type Client struct {
	pathToBinary string

	mu             sync.Mutex
	command        *exec.Cmd
	stdInPipe      *io.PipeWriter
	hasBeenStarted bool
}

func NewClient(pathToBinary string) Client {
	return Client{
		pathToBinary: pathToBinary,
	}
}

func (c *Client) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.hasBeenStarted {
		return errors.New("uci: Tried to start a client that had already been started.")
	}

	// TODO(guy) check that the given file actually exists
	logWriter := cmd.NewLogWriter()
	c.command = cmd.NewCommand(c.pathToBinary)
	c.command.Stdout = logWriter
	c.command.Stderr = logWriter

	pipeReader, pipeWriter := io.Pipe()
	c.stdInPipe = pipeWriter
	c.command.Stdin = pipeReader

	err := c.command.Start()
	if err != nil {
		return err
	}

	c.hasBeenStarted = true
	return nil
}

func (c *Client) Stop() error {
	if c.stdInPipe == nil || c.command == nil {
		return errors.New("uci: Tried to stop a client that had never been running.")
	}

	err := c.stdInPipe.Close()
	if err != nil {
		return err
	}

	return c.command.Wait()
}

func (c *Client) sendMessage(msg string) error {
	n, err := c.stdInPipe.Write([]byte(msg + "\n"))
	if err != nil {
		return err
	}

	if n != len(msg)+1 {
		return fmt.Errorf("uci: Failed to write all of \"%s\" to engine.", msg)
	}

	return nil
}
