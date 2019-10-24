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
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const escape rune = 27
const esc string = string(escape)

type signalTimes struct {
	signal os.Signal
	time   time.Time
}

func processSignals() chan signalTimes {
	c := make(chan os.Signal, 10)
	ret := make(chan signalTimes, 10)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			ret <- signalTimes{
				signal: s,
				time:   time.Now(),
			}
		}
	}()
	return ret
}

func main() {
	signals := processSignals()
	config, err := arguments()
	if err != nil {
		panic(err)
	}
	if config.ListLangs {
		for k := range languages {
			fmt.Println(k)
		}
		return
	}
	if config.Show {
		if config.Filename == "" {
			if fn, ok := file[config.Lang]; ok {
				config.Filename = fn
			}
		}
		if config.Filename == "" {
			panic(fmt.Errorf("cannot find filename for lang %s, please provide one", config.Lang))
		}
		printFile(config.Filename, config.Name, config.Lang)
		return
	}

	var examples Examples
	var ir io.Reader
	fname := config.Name + "/io.txt"
	writeInp := false
	extracted := false
	var buf []byte
	if config.Stdin {
		fmt.Fprintln(os.Stderr, "Waiting for stdin test cases...")
		buf, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		ir = bytes.NewReader(buf)
	} else if _, err := os.Stat(fname); !config.ForceDownload && err == nil {
		ir, err = os.Open(fname)
		if err != nil {
			panic(err)
		}
	} else if config.ForceDownload || os.IsNotExist(err) {
		writeInp = true
		exs, err := extractString(config.Name)
		buf = []byte(exs)
		if err != nil {
			panic(err)
		}
		extracted = true
		ir = strings.NewReader(exs)
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
	if config.Examples {
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

	if config.StdinOne {
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		examples = append(examples, Example{
			Input:  string(buf),
			Output: "",
			noOut:  true,
		})
	}

	if config.BuildCmd != "" {
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		ctx, cancel := context.WithTimeout(context.Background(), config.buildTimeout)
		cmd := exec.CommandContext(ctx, "bash", "-c", config.BuildCmd)
		cmd.Dir = config.Name
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
	if config.Only != 0 && config.Only > len(examples) {
		panic(fmt.Sprintf("have %d test cases, wanted to run case number %d", len(examples), config.Only))
	}
	for i, el := range examples {
		if config.Only != 0 && config.Only != i+1 {
			continue
		}
		if !config.StrictEllipsis && strings.HasSuffix(strings.TrimSpace(el.Input), "...") {
			fmt.Printf(esc + "[37m")
			fmt.Printf("case %d skipped due to incomplete input.\n", i+1)
			fmt.Printf(esc + "[0m")
			continue
		}
		stdout := bytes.Buffer{}
		stderr := bytes.Buffer{}
		ctx, cancel := context.WithTimeout(context.Background(), config.timeout)
		cmd := exec.CommandContext(ctx, "bash", "-c", config.Cmd)
		cmd.Dir = config.Name
		inp := el.Input
		if config.NoWindowsNewline && strings.Contains(inp, "\r\n") {
			inp = strings.Replace(inp, "\r\n", "\n", -1)
		} else if !config.NoWindowsNewline && !strings.Contains(inp, "\r\n") {
			inp = strings.Replace(inp, "\n", "\r\n", -1)
		}
		cmd.Stdin = strings.NewReader(inp)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		start := time.Now()
		err := cmd.Run()
		end := time.Now()
		took := end.Sub(start)
		cancel()
		out := strings.TrimSpace(stdout.String())
		failed := false
		if err != nil || !cmd.ProcessState.Success() {
			failed = true
		} else if !el.noOut {
			oo := strings.Split(out, "\n")
			ee := strings.Split(el.Output, "\n")
			if len(oo) < len(ee) || (len(oo) != len(ee) && !(!config.StrictEllipsis && strings.HasSuffix(strings.TrimSpace(el.Output), "..."))) {
				failed = true
			} else {
				for i := range oo {
					if !config.StrictEllipsis && strings.HasSuffix(strings.TrimSpace(ee[i]), "...") {
						eee := strings.TrimSpace(ee[i])
						if !strings.HasPrefix(strings.TrimSpace(oo[i]), eee[:len(eee)-3]) {
							failed = true
							break
						}
						break
					} else {
						if strings.TrimSpace(oo[i]) != strings.TrimSpace(ee[i]) {
							failed = true
							break
						}
					}
					if config.MatchFirstLine {
						break
					}
				}
			}
		}
		if !failed {
			fmt.Printf(esc + "[32m")
			fmt.Printf("case %d completed successfully. took: %s", i+1, took)
			fmt.Printf(esc + "[0m")
		} else {
			fmt.Printf(esc + "[4;31m")
			fmt.Printf("case %d failed. took: %s", i+1, took)
			fmt.Printf(esc + "[0m")
			fmt.Printf(esc + "[1;31m")
			if err != nil {
				fmt.Printf("returned error: %s", err.Error())
			}
			fmt.Printf(esc + "[0m")
		}
		if len(signals) > 0 {
			totalTime := end.Sub(start)
			prev := start
			done := false
			for !done {
				select {
				case v := <-signals:
					if config.NoPercentage {
						fmt.Printf("\t%s", v.time.Sub(prev))
					} else {
						fmt.Printf("\t%02.2f%%", float64(v.time.Sub(prev))/float64(totalTime)*100)
					}
					prev = v.time
				case <-time.After(10 * time.Millisecond):
					done = true
				}
			}
			if config.NoPercentage {
				fmt.Printf("\t%s", end.Sub(prev))
			} else {
				fmt.Printf("\t%02.2f%%", float64(end.Sub(prev))/float64(totalTime)*100)
			}
		}
		fmt.Printf("\n")
		if !config.Quite && (failed || config.Verbose) {
			if stderr.Len() != 0 {
				fmt.Printf(stderr.String())
			}
			fmt.Println(el.String())
			fmt.Println("output")
			fmt.Printf(esc + "[31m")
			fmt.Println(out)
			fmt.Printf(esc + "[0m")
		}
		if failed && config.ExitOnFailure {
			return
		}
	}
}

func (e *Example) String() string {
	ret := "input\n" + esc + "[35m" + e.Input + esc + "[0m"
	if !e.noOut {
		ret += "\nexpected\n" + esc + "[34m" + e.Output + esc + "[0m"
	}
	return ret
}

func (e *Examples) String() string {
	s := ""
	for _, i := range *e {
		s += i.String()
		s += "\n"
	}
	return s
}
