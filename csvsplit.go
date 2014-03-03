package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type options struct {
	repeatRows, limit int
	source            string
	headerRows        []string
}

func processCmdLineFlags(opts *options) {
	flag.StringVar(&(*opts).source, "source", "", "Name of the source file (if not using stdin)")
	flag.IntVar(&(*opts).repeatRows, "repeat", 1, "How many rows to repeat on subsequent files")
	flag.IntVar(&(*opts).limit, "limit", 1000, "Number of rows per output file (not counting repeated/header rows)")
	flag.Parse()
}

func scan(what io.Reader, opts *options) {
	var f *os.File
	scanner, write := bufio.NewScanner(what), makeWriter(opts)
	for j := 0; j < opts.repeatRows && scanner.Scan(); j++ {
		opts.headerRows = append(opts.headerRows, scanner.Text())
	}

	for scanner.Scan() {
		f = write(scanner.Text())
	}
	f.Close()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func makeWriter(opts *options) func(string) *os.File {
	var curFile *os.File
	var curBuf *bufio.Writer
	lineno, prefix := 0, determineFilePrefix(opts)

	return func(row string) *os.File {
		if lineno%opts.limit == 0 { // Time to change the file
			if curFile != nil {
				if err := curBuf.Flush(); err != nil {
					log.Fatal(err)
				}
				curFile.Close()
			}

			curFile = createFile(fmt.Sprintf("%s_%02d.csv", prefix, (lineno/opts.limit)+1))
			curBuf = bufio.NewWriter(curFile)

			for _, headerRow := range opts.headerRows {
				bufWriteln(curBuf, headerRow)
			}
		}

		bufWriteln(curBuf, row)
		lineno++

		return curFile
	}
}

func main() {
	opts := new(options)
	if processCmdLineFlags(opts); stdinReady() {
		scan(os.Stdin, opts)
	} else if opts.source != "" {
		f := openFile(opts.source)
		defer f.Close()
		scan(f, opts)
	} else {
		fmt.Println("No source available, please see help (also available with -h):")
		flag.PrintDefaults()
	}
}
