package uci

import (
	"bufio"
	"clanpj/lisao/cmd"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Client struct {
	pathToBinary string

	mu             sync.Mutex
	command        *exec.Cmd
	stdInPipe      *io.PipeWriter
	stdOutPipe     *io.PipeReader
	stdOutBuf      *bufio.Reader
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

	if !fileExists(c.pathToBinary) {
		return fmt.Errorf("uci: Path to binary isn't a file: %s", c.pathToBinary)
	}
	c.command = cmd.NewCommand(c.pathToBinary)

	pipeReader, pipeWriter := io.Pipe()
	logWriter := cmd.NewLogWriter("uci/client: ")
	c.stdOutPipe = pipeReader
	c.stdOutBuf = bufio.NewReader(pipeReader)
	c.command.Stderr = logWriter
	c.command.Stdout = io.MultiWriter(logWriter, pipeWriter)

	pipeReader, pipeWriter = io.Pipe()
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

	err = c.stdOutPipe.Close()
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

func (c *Client) getLine() (string, error) {
	msg, err := c.stdOutBuf.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(msg), nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == os.ErrNotExist {
		return false
	}

	return true
}
