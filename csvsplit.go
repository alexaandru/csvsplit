package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type options struct {
	repeatRows, limit int
	source            string
}

var (
	curFile *os.File
	prefix  string
)

func processCmdLineFlags(opts *options) {
	flag.StringVar(&(*opts).source, "source", "", "Name of the source file (if not using stdin)")
	flag.IntVar(&(*opts).repeatRows, "repeat", 1, "How many rows to repeat on subsequent files")
	flag.IntVar(&(*opts).limit, "limit", 1000, "Number of rows per output file (not counting repeated/header rows)")
	flag.Parse()

	if opts.source != "" {
		prefix = strings.TrimSuffix(opts.source, filepath.Ext(opts.source))
	} else {
		prefix = "file"
	}
}

func scan(what io.Reader, opts *options) {
	defer curFile.Close()
	scanner, lineno, repeats := bufio.NewScanner(what), 0, []string{}
	for j := 0; j < opts.repeatRows && scanner.Scan(); j++ {
		repeats = append(repeats, scanner.Text())
	}

	for scanner.Scan() {
		write(scanner.Text(), lineno, &repeats, opts.limit)
		lineno++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func write(line string, lineno int, repeats *[]string, limit int) {
	if lineno%limit == 0 {
		if curFile != nil {
			curFile.Close()
		}

		curFile = createFile(fmt.Sprintf("%s_%02d.csv", prefix, (lineno/limit)+1))
		for _, repeat := range *repeats {
			fmt.Fprintln(curFile, repeat)
		}
	}

	fmt.Fprintln(curFile, line)
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
