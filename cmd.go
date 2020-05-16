/*
timeit.go
-John Taylor
May-16-2020

A CLI tool used to time the duration of the given command
The result is sent to STDERR
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const version = "1.0.1"

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "\ntimeit v%s\n", version)
		fmt.Fprintf(os.Stderr, "\nUsage: %s [cmd] [args...]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "You may need to surround arguments in double-quotes\n")
		os.Exit(1)
	}
	timeStart := time.Now()
	out, err := exec.Command(os.Args[1], os.Args[2:len(os.Args)]...).CombinedOutput()
	elapsed := time.Since(timeStart)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Command finished with error: %v\n", err)
	}
	fmt.Printf("%s", out)
	fmt.Fprintln(os.Stderr, elapsed)
}
