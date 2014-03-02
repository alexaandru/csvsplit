package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func stdinReady() bool {
	fi, e := os.Stdin.Stat()
	return e == nil && (fi.Mode()&os.ModeNamedPipe != 0 || fi.Size() > 0)
}

func openFile(fname string) *os.File {
	if f, err := os.Open(fname); err == nil {
		return f
	} else {
		log.Fatal(err)
		return nil
	}
}

func createFile(fname string) *os.File {
	if f, err := os.Create(fname); err == nil {
		return f
	} else {
		log.Fatal(err)
		return nil
	}
}

func determineFilePrefix(opts *options) string {
	if opts.source != "" {
		return strings.TrimSuffix(opts.source, filepath.Ext(opts.source))
	}

	return "file"
}

func bufWrite(buf *bufio.Writer, line string) {
	if _, err := buf.WriteString(line); err != nil {
		log.Fatal(err)
	}
}

func bufWriteln(buf *bufio.Writer, line string) {
	bufWrite(buf, line)
	bufWrite(buf, "\n")
}
