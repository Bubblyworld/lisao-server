package cmd

import (
	"bytes"
	"log"
	"regexp"
	"testing"
)

func TestLogWriter(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})
	log.SetOutput(buffer)

	logWriter := NewLogWriter("prefix: ")
	logWriter.Write([]byte("Testing.\n"))
	logWriter.Close()

	// Check it got written.
	reg := regexp.MustCompile("Testing.")
	if !reg.Match(buffer.Bytes()) {
		t.Errorf("Log didn't write correct string.")
	}
}
