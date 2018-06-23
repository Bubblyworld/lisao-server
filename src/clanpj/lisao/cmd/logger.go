package cmd

import (
	"bufio"
	"io"
	"log"
)

type LogWriter struct {
	pipeReader *io.PipeReader
	pipeWriter *io.PipeWriter

	doneCh chan bool
}

// LogWriter returns an io.WriteCloser that log.Prints everything coming through.
func NewLogWriter() LogWriter {
	r, w := io.Pipe()

	logWriter := LogWriter{
		pipeReader: r,
		pipeWriter: w,
		doneCh:     make(chan bool),
	}

	go logWriter.copy()
	return logWriter
}

func (lw LogWriter) copy() {
	bufReader := bufio.NewReader(lw.pipeReader)

	// TODO(guy) don't swallow non-io.EOF errors. (or non ErrClosedPipe)
	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			lw.doneCh <- true
			return
		}

		log.Println(line)
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

	<-lw.doneCh
	return nil
}
