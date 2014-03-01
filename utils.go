package main

import (
	"log"
	"os"
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
