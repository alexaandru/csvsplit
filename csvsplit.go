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
	flag.IntVar(&(*opts).limit, "limit", 2, "Number of rows per output file (not counting repeated/header rows)")
	flag.Parse()

	if opts.source != "" {
		prefix = strings.TrimSuffix(opts.source, filepath.Ext(opts.source))
	} else {
		prefix = "file"
	}
}

func stdinReady() bool {
	if fi, err := os.Stdin.Stat(); err != nil {
		return false
	} else {
		return fi.Mode()&os.ModeNamedPipe != 0 || fi.Size() > 0
	}
}

func scan(what io.Reader, opts *options) {
	defer curFile.Close()
	scanner, i, repeats := bufio.NewScanner(what), 0, []string{}
	for i = 0; i < opts.repeatRows && scanner.Scan(); i++ {
		repeats = append(repeats, scanner.Text())
	}

	i = 0
	for scanner.Scan() {
		write(scanner.Text(), i, &repeats, opts.limit)
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func write(line string, i int, repeats *[]string, limit int) {
	if i%limit == 0 {
		if curFile != nil {
			curFile.Close()
		}

		curFile = createFile(fmt.Sprintf("%s_%02d.csv", prefix, (i/limit)+1))
		for _, repeat := range *repeats {
			fmt.Fprintln(curFile, repeat)
		}
	}

	fmt.Fprintln(curFile, line)
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
