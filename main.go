package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
)

const escape rune = 27
const esc string = string(escape)

func normalizeName(name string) (string, error) {
	var class string
	var num int
	var err error
	if name[len(name)-2] == '/' {
		_, err = fmt.Sscanf(name, "%d/%s", &num, &class)
	} else {
		class = name[0:1]
		num, err = strconv.Atoi(name[1:])
		if err != nil {
			return "", nil
		}
	}
	return fmt.Sprintf("%s%d", class, num), nil
}

func main() {
	usage := `codeforces test runner.

Usage:
  codeforces <name> [--match-first-line] [--cmd=<cmd>] [--stdin] [--timeout=<timeout>]
  codeforces -h | --help
  codeforces --version

Options:
  -h --help              Show this screen.
  --version              Show version.
  --match-first-line     Only match output first line [default: false].
  --cmd=<cmd>            Command to execute the program [default: go build -o main && ./main]
  --stdin                Get input from stdin [default: false].
  --timeout=<timeout>    Timeout for a single case [default: 1s].
`

	arguments, _ := docopt.ParseDoc(usage)
	name, err := arguments.String("<name>")
	if err != nil {
		panic(err)
	}
	name, err = normalizeName(name)
	if err != nil {
		panic(err)
	}
	firstLine, err := arguments.Bool("--match-first-line")
	if err != nil {
		panic(err)
	}
	cmd, err := arguments.String("--cmd")
	if err != nil {
		panic(err)
	}
	stdin, err := arguments.Bool("--stdin")
	if err != nil {
		panic(err)
	}
	timeoutString, err := arguments.String("--timeout")
	if err != nil {
		panic(err)
	}
	timeout, err := time.ParseDuration(timeoutString)
	if err != nil {
		panic(err)
	}

	var examples Examples
	var ir io.Reader
	fname := name + "/io.txt"
	writeInp := false
	var buf []byte
	if _, err := os.Stat(fname); err == nil {
		ir, err = os.Open(fname)
		if err != nil {
			panic(err)
		}
	} else if os.IsNotExist(err) {
		writeInp = true
		if stdin {
			fmt.Fprintln(os.Stderr, "Waiting for stdin test cases...")
			buf, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}
			ir = bytes.NewReader(buf)
		} else {
			exs, err := extractString(name)
			buf = []byte(exs)
			if err != nil {
				panic(err)
			}
			ir = strings.NewReader(exs)
		}
	} else {
		panic(err)
	}
	dec := json.NewDecoder(ir)
	err = dec.Decode(&examples)
	if err != nil {
		panic(err)
	}
	if writeInp {
		file, err := os.Create(fname)
		if err != nil {
			return
		}
		defer file.Close()

		_, err = file.Write(buf)
		if err != nil {
			panic(err)
		}
	}

	for i, el := range examples {
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		cmd := exec.CommandContext(ctx, "bash", "-c", cmd)
		cmd.Dir = name
		cmd.Stdin = strings.NewReader(el.Input)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		cancel()
		if err == nil && ((firstLine && strings.Split(stdout.String(), "\n")[0] == strings.Split(el.Output, "\n")[0]) || stdout.String() == el.Output) {
			fmt.Printf(esc + "[32m")
			fmt.Printf("case %d completed successfully\n", i)
			fmt.Printf(esc + "[0m")
		} else {
			fmt.Printf(esc + "[4;31m")
			fmt.Printf("case %d failed\n", i)
			fmt.Printf(esc + "[0m")
			fmt.Printf(esc + "[1;31m")
			if err != nil {
				fmt.Printf("returned error: %s\n", err.Error())
			}
			if stderr.Len() != 0 {
				fmt.Printf(stderr.String())
			}
			fmt.Printf(esc + "[0m")
			fmt.Println("input")
			fmt.Println(el.Input)
			fmt.Println("expected")
			fmt.Printf(esc + "[34m")
			fmt.Println(el.Output)
			fmt.Printf(esc + "[0m")
			fmt.Println("output")
			fmt.Printf(esc + "[31m")
			fmt.Println(stdout.String())
			fmt.Printf(esc + "[0m")
		}
	}
}
