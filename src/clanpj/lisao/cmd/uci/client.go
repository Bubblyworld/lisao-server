package uci

import (
	"bufio"
	"clanpj/lisao/cmd"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type Client struct {
	pathToBinary string
	logTo        io.Writer

	mu         sync.Mutex
	command    *exec.Cmd
	stdInPipe  *io.PipeWriter
	stdOutPipe *io.PipeReader
	stdOutBuf  *bufio.Reader
	isRunning  bool
}

func NewClient(pathToBinary string, logTo io.Writer) Client {
	return Client{
		pathToBinary: pathToBinary,
	}
}

func (c *Client) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isRunning {
		return errors.New("uci: tried to start a client that is running")
	}

	if !fileExists(c.pathToBinary) {
		log.Printf("uci: Path to binary isn't a file: %s", c.pathToBinary)
		return fmt.Errorf("uci: path to binary isn't a file")
	}
	c.command = cmd.NewCommand(c.pathToBinary)

	pipeReader, pipeWriter := io.Pipe()
	c.stdOutPipe = pipeReader
	c.stdOutBuf = bufio.NewReader(pipeReader)

	if c.logTo != nil {
		c.command.Stderr = c.logTo
	}

	c.command.Stdout = pipeWriter
	if c.logTo != nil {
		c.command.Stdout = io.MultiWriter(c.logTo, c.command.Stdout)
	}

	pipeReader, pipeWriter = io.Pipe()
	c.stdInPipe = pipeWriter
	c.command.Stdin = pipeReader

	err := c.command.Start()
	if err != nil {
		return err
	}

	c.isRunning = true
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

	err = c.command.Wait()
	if err != nil {
		return err
	}

	c.isRunning = false
	return nil
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
