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

var build = map[string]string{
	"go":     "go build -o main",
	"python": "python -m py_compile main.py",
	"c":      "gcc main.c -o main -Wall",
	"c++":    "g++ main.cpp -Wall -o main",
}

var run = map[string]string{
	"go":         "./main",
	"python":     "python main.py",
	"javascript": "node index.js",
	"c":          "./main",
	"c++":        "./main",
}

var languages = map[string]string{
	"go":         "go",
	"python":     "python",
	"javascript": "javascript",
	"node":       "javascript",
	"c":          "c",
	"c++":        "c++",
	"cpp":        "c++",
}

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
  codeforces run <name> [--match-first-line] [--cmd=<cmd>] [--build-cmd=<cmd>] [--stdin] [--timeout=<timeout>] [--build-timeout=<timeout>] [--force-download] [--lang=<lang>] [--exit-on-failure] [--verbose]
  codeforces examples <name> [--force-download]
  codeforces list-langs
  codeforces -h | --help
  codeforces --version

Options:
  -h --help                                            Show this screen.
  --version                                            Show version.
  --match-first-line                                   Only match output first line [default: false].
  --lang=<lang>                                        Source code language use "codeforces list-langs" to list languages [default: go].
  --build-cmd=<cmd>                                    Command to execute the program, overrides lang [default: ].
  --cmd=<cmd>                                          Command to execute the program, overrides lang [default: ].
  --stdin                                              Get input from stdin [default: false].
  --timeout=<timeout>                                  Timeout for a single case [default: 1s].
  --build-timeout=<timeout>                            Timeout for build [default: 10s].
  --force-download                                     Force download examples [default: false]
  --exit-on-failure                                    Exit on the first failed example [default: false].
  --verbose                                            Always show input/expected/output [default: false].
`

	arguments, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}
	examplesOnly, err := arguments.Bool("examples")
	if err != nil {
		panic(err)
	}
	listLangs, err := arguments.Bool("list-langs")
	if err != nil {
		panic(err)
	}
	if listLangs {
		for k := range languages {
			fmt.Println(k)
		}
		return
	}
	forceDownload, err := arguments.Bool("--force-download")
	if err != nil {
		panic(err)
	}
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
	runCmd, err := arguments.String("--cmd")
	if err != nil {
		panic(err)
	}
	buildCmd, err := arguments.String("--build-cmd")
	if err != nil {
		panic(err)
	}
	lang, err := arguments.String("--lang")
	if err != nil {
		panic(err)
	}
	lang, ok := languages[lang]
	if !ok {
		panic(fmt.Errorf("unknown language %s", lang))
	}
	langRunCmd, ok := run[lang]
	if !ok {
		panic(fmt.Errorf("unknown language run %s", lang))
	}
	if runCmd == "" {
		runCmd = langRunCmd
	}
	langBuildCmd, ok := build[lang]
	if ok {
		if buildCmd == "" {
			buildCmd = langBuildCmd
		}
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
	buildTimeoutString, err := arguments.String("--build-timeout")
	if err != nil {
		panic(err)
	}
	buildTimeout, err := time.ParseDuration(buildTimeoutString)
	if err != nil {
		panic(err)
	}
	exitOnFailure, err := arguments.Bool("--exit-on-failure")
	if err != nil {
		panic(err)
	}
	verbose, err := arguments.Bool("--verbose")
	if err != nil {
		panic(err)
	}

	var examples Examples
	var ir io.Reader
	fname := name + "/io.txt"
	writeInp := false
	extracted := false
	var buf []byte
	if _, err := os.Stat(fname); !forceDownload && err == nil {
		ir, err = os.Open(fname)
		if err != nil {
			panic(err)
		}
	} else if forceDownload || os.IsNotExist(err) {
		writeInp = true
		if !forceDownload && stdin {
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
			extracted = true
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
	if examplesOnly {
		if extracted {
			fmt.Println("extract successful")
		} else {
			fmt.Println("not extracted")
		}
		if writeInp {
			fmt.Println("examples saved")
		} else {
			fmt.Println("examples not saved")
		}
		fmt.Println(examples.String())
		return
	}
	if buildCmd != "" {
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		ctx, cancel := context.WithTimeout(context.Background(), buildTimeout)
		cmd := exec.CommandContext(ctx, "bash", "-c", buildCmd)
		cmd.Dir = name
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		start := time.Now()
		err := cmd.Run()
		took := time.Now().Sub(start)
		cancel()
		if err == nil && cmd.ProcessState.Success() {
			fmt.Printf(esc + "[32m")
			fmt.Printf("build succeeded. took: %s\n", took)
			fmt.Printf(esc + "[0m")
		} else {
			fmt.Printf(esc + "[4;31m")
			fmt.Printf("build failed. took: %s\n", took)
			fmt.Printf(esc + "[0m")
			fmt.Printf(esc + "[35m")
			fmt.Printf(stdout.String())
			fmt.Printf(esc + "[0m")
			fmt.Printf(esc + "[31m")
			fmt.Printf(stderr.String())
			fmt.Printf(esc + "[0m")
			if err != nil {
				fmt.Println(err)
			}
			return
		}
	}
	for i, el := range examples {
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		cmd := exec.CommandContext(ctx, "bash", "-c", runCmd)
		cmd.Dir = name
		cmd.Stdin = strings.NewReader(el.Input)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		start := time.Now()
		err := cmd.Run()
		took := time.Now().Sub(start)
		cancel()
		out := strings.TrimSpace(stdout.String())
		failed := false
		if err != nil || !cmd.ProcessState.Success() {
			failed = true
		} else {
			oo := strings.Split(out, "\n")
			ee := strings.Split(el.Output, "\n")
			for i := range oo {
				if strings.TrimSpace(oo[i]) != strings.TrimSpace(ee[i]) {
					failed = true
					break
				}
				if firstLine {
					break
				}
			}
		}
		if !failed {
			fmt.Printf(esc + "[32m")
			fmt.Printf("case %d completed successfully. took: %s\n", i+1, took)
			fmt.Printf(esc + "[0m")
		} else {
			fmt.Printf(esc + "[4;31m")
			fmt.Printf("case %d failed. took: %s\n", i+1, took)
			fmt.Printf(esc + "[0m")
			fmt.Printf(esc + "[1;31m")
			if err != nil {
				fmt.Printf("returned error: %s\n", err.Error())
			}
		}
		if failed || verbose {
			if stderr.Len() != 0 {
				fmt.Printf(stderr.String())
			}
			fmt.Printf(esc + "[0m")
			fmt.Println(el.String())
			fmt.Println("output")
			fmt.Printf(esc + "[31m")
			fmt.Println(out)
			fmt.Printf(esc + "[0m")
		}
		if failed && exitOnFailure {
			return
		}
	}
}

func (e *Example) String() string {
	return "input\n" + esc + "[35m" + e.Input + "\n" + esc + "[0mexpected\n" + esc + "[34m" + e.Output + esc + "[0m"
}

func (e *Examples) String() string {
	s := ""
	for _, i := range *e {
		s += i.String()
		s += "\n"
	}
	return s
}
