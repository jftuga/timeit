/*
timeit.go
-John Taylor
May-16-2020

A CLI tool used to time the duration of the given command
The result is sent to STDERR
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"time"
)

const version = "1.2.2"

var timeStart time.Time

func checkBuf(buf *bufio.Reader, wg *sync.WaitGroup) {
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if "EOF" != err.Error() {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			break
		}
		fmt.Printf("%s", line)
	}
	wg.Done()
}

func run(cmd *exec.Cmd, outBuf, errBuf *bufio.Reader) {
	timeStart = time.Now()
	cmd.Start()

	var wg sync.WaitGroup
	wg.Add(2)
	go checkBuf(outBuf, &wg)
	go checkBuf(errBuf, &wg)
	wg.Wait()

}

func ioSetup(cmd *exec.Cmd) (*bufio.Reader, *bufio.Reader) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	outBuf := bufio.NewReader(stdout)
	errBuf := bufio.NewReader(stderr)
	return outBuf, errBuf
}

func ctrlCHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			elapsed := time.Since(timeStart)
			fmt.Fprintln(os.Stderr, sig)
			fmt.Fprintln(os.Stderr, elapsed)
			os.Exit(1)
		}
	}()
}

func usage() {
	fmt.Fprintf(os.Stderr, "\ntimeit v%s\n", version)
	fmt.Fprintf(os.Stderr, "https://github.com/jftuga/timeit\n")
	fmt.Fprintf(os.Stderr, "A CLI tool used to time the duration of the given cmd\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [cmd] [args...]\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "You may need to surround args within double-quotes\n")
}

func main() {
	if len(os.Args) == 1 {
		usage()
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[1], os.Args[2:len(os.Args)]...)
	outBuf, errBuf := ioSetup(cmd)
	ctrlCHandler()
	run(cmd, outBuf, errBuf)
	elapsed := time.Since(timeStart)
	fmt.Fprintln(os.Stderr, elapsed)
}
