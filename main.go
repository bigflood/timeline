package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	startTime time.Time
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Printf("Usage: %s command args ..\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	var wg sync.WaitGroup

	cmd := exec.Command(args[0], args[1:]...)

	if stdin, err := cmd.StdinPipe(); err != nil {
		fmt.Fprintf(os.Stderr, "StdinPipe failed: %v\n", err)
		os.Exit(2)
	} else {
		go io.Copy(stdin, os.Stdin)
		//go outToLog("I> ", os.Stdin, stdin, nil, nil)
	}

	if stdout, err := cmd.StdoutPipe(); err != nil {
		fmt.Fprintf(os.Stderr, "StdoutPipe failed: %v\n", err)
		os.Exit(2)
	} else {
		wg.Add(1)
		go outToLog("O> ", stdout, os.Stdout, os.Stdout, &wg)
	}

	if stderr, err := cmd.StderrPipe(); err != nil {
		fmt.Fprintf(os.Stderr, "StdoutPipe failed: %v\n", err)
		os.Exit(2)
	} else {
		wg.Add(1)
		go outToLog("E> ", stderr, os.Stderr, os.Stderr, &wg)
	}

	startTime = time.Now()
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Start failed: %v\n", err)
		os.Exit(2)
	}

	err := cmd.Wait()
	exitTime := getTimeStr()
	wg.Wait()

	fmt.Println()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Fprintf(os.Stderr, "%v\n", exitErr)
		} else {
			fmt.Fprintf(os.Stderr, "Wait failed: %v\n", err)
		}
	}

	fmt.Printf("Command time: %s \n", exitTime)
	fmt.Printf("System CPU time: %s \n", cmd.ProcessState.SystemTime())
	fmt.Printf("User CPU time: %s \n", cmd.ProcessState.UserTime())

}

func outToLog(prefix string, input io.Reader, output io.Writer, timeOutput io.Writer, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	prefixBytes := []byte(prefix)
	buf := make([]byte, 1024*10)
	sep := []byte("\n")
	space := []byte(" ")
	beginLine := true

	for {
		n, err := input.Read(buf)
		if err != nil {
			break
		}

		if n == 0 {
			continue
		}

		slices := bytes.Split(buf[:n], sep)

		for i, s := range slices {
			if i > 0 {
				output.Write(sep)
				beginLine = true
			}

			if len(s) == 0 {
				beginLine = true
				continue
			}

			if beginLine && timeOutput != nil {
				timeOutput.Write(prefixBytes)
				timeOutput.Write([]byte(getTimeStr()))
				timeOutput.Write(space)
			}

			output.Write(s)
			beginLine = false
		}
	}
}

func getTimeStr() string {
	if !startTime.IsZero() {
		return fmt.Sprint(time.Since(startTime), " ")
	}

	return ""
}

func Log(s string) {
	s = strings.TrimRight(s, " \n\r\t")
	ts := getTimeStr()

	if ts != "" {
		fmt.Printf("%s %s\n", ts, s)
	} else {
		fmt.Printf("%s\n", s)
	}
}

func Logf(f string, args ...interface{}) {
	s := fmt.Sprintf(f, args...)
	Log(s)
}
