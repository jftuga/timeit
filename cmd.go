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
	"sync"
	"time"
)

const version = "1.2.0"

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
	cmd.Start()

	var wg sync.WaitGroup
	wg.Add(2)
	go checkBuf(outBuf, &wg)
	go checkBuf(errBuf, &wg)
	wg.Wait()

}

func setup(cmd *exec.Cmd) (*bufio.Reader, *bufio.Reader) {
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

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "\ntimeit v%s\n", version)
		fmt.Fprintf(os.Stderr, "\nUsage: %s [cmd] [args...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "You may need to surround arguments in double-quotes\n")
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[1], os.Args[2:len(os.Args)]...)
	outBuf, errBuf := setup(cmd)
	ctrlCHandler()
	timeStart = time.Now()
	run(cmd, outBuf, errBuf)
	elapsed := time.Since(timeStart)
	fmt.Fprintln(os.Stderr, elapsed)
}
