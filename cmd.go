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
	"strings"
	"sync"
	"time"
)

const VERSION = "1.3.2"
const URL = "https://github.com/jftuga/timeit"
const TMPFILE = ".timeit.start.tmp"

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

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func createStartFile() {
	f, err := os.Create(TMPFILE)
	check(err)
	f.WriteString(fmt.Sprintf("%s", time.Now().Format("2006-01-02 15:04:05.000000 -0700 MST")))
	f.Close()
}

func getStartFileTime() string {
	f, err := os.Open(TMPFILE)
	check(err)
	b1 := make([]byte, 100)
	n1, err := f.Read(b1)
	check(err)
	f.Close()
	//fmt.Printf("%d bytes: %s\n", n1, string(b1[:n1]))
	return string(b1[:n1])
}

func getElapsedTime(startTime string) time.Duration {
	timeT, err := time.Parse("2006-01-02 15:04:05.000000 -0700 MST", strings.Split(startTime, " m=")[0])
	check(err)
	elapsedTime := time.Since(timeT)
	err = os.Remove(TMPFILE)
	check(err)
	return elapsedTime
}

func version() {
	fmt.Fprintf(os.Stderr, "timeit v%s\n", VERSION)
	fmt.Fprintf(os.Stderr, "%s\n", URL)
}

func usage() {
	fmt.Fprintf(os.Stderr, "\ntimeit v%s\n\n", VERSION)
	fmt.Fprintf(os.Stderr, "A cross-platform CLI tool used to time the duration of the given command\n\n")
	fmt.Fprintf(os.Stderr, "Usage: %s [cmd] [args...]\n", filepath.Base(os.Args[0]))
	fmt.Fprintf(os.Stderr, "You may need to surround args within double-quotes\n\n")
	fmt.Fprintf(os.Stderr, "Examples:\n")
	fmt.Fprintf(os.Stderr, "    timeit wget https://example.com/file.tar.gz\n")
	fmt.Fprintf(os.Stderr, "    timeit gzip -t \"file with spaces.gz\"\n\n")
	fmt.Fprintf(os.Stderr, "For built-in Windows 'cmd' commands:\n")
	fmt.Fprintf(os.Stderr, "    timeit cmd /c \"dir c:\\ /s/b > list.txt\"\n")
	fmt.Fprintf(os.Stderr, "    timeit cmd /c dir /s \"c:\\Program Files\"\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Run 'timeit _start' to create (or overwrite) this file containing the current time: %s\n", TMPFILE)
	fmt.Fprintf(os.Stderr, "Run 'timeit _end' to read (and then delete) that file.  The elapsed time will then be displayed.\n")
	fmt.Fprintf(os.Stderr, "This can be useful for timing multiple, long-running commands.\n")
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	if len(os.Args) == 1 {
		usage()
		os.Exit(0)
	}

	if len(os.Args) == 2 && os.Args[1] == "-v" {
		version()
		os.Exit(0)
	}
	if strings.ToLower(os.Args[1]) == "_start" {
		createStartFile()
		return
	}

	if strings.ToLower(os.Args[1]) == "_end" {
		elapsedTime := getElapsedTime(getStartFileTime())
		fmt.Fprintln(os.Stderr, elapsedTime)
		return
	}

	cmd := exec.Command(os.Args[1], os.Args[2:len(os.Args)]...)
	outBuf, errBuf := ioSetup(cmd)
	ctrlCHandler()
	run(cmd, outBuf, errBuf)

	fmt.Fprintln(os.Stderr, time.Since(timeStart))
}
