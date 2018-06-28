package cmd

import (
	"bufio"
	"io"
	"log"
)

type LogWriter struct {
	prefix string

	pipeReader *io.PipeReader
	pipeWriter *io.PipeWriter
	doneCh     chan error
}

// LogWriter returns an io.WriteCloser that log.Prints everything coming through.
func NewLogWriter(prefix string) LogWriter {
	r, w := io.Pipe()

	logWriter := LogWriter{
		prefix:     prefix,
		pipeReader: r,
		pipeWriter: w,
		doneCh:     make(chan error),
	}

	go logWriter.copy()
	return logWriter
}

func (lw LogWriter) copy() {
	bufReader := bufio.NewReader(lw.pipeReader)

	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			lw.doneCh <- err
			return
		}

		log.Print(lw.prefix + line)
	}
}

func (lw LogWriter) Write(data []byte) (int, error) {
	return lw.pipeWriter.Write(data)
}

func (lw LogWriter) Close() error {
	err := lw.pipeWriter.Close()
	if err != nil {
		return err
	}

	return <-lw.doneCh
}
