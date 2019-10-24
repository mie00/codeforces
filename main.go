package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
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
		printFile(config.Filename, config.Name, config.Lang)
		return
	}

	examples, err := getExamples(*config)

	if config.New {
		err = newFile(*config)
		if err != nil {
			panic(err)
		}
		return
	}

	if config.Examples {
		fmt.Println(examples.String())
		return
	}

	signals := processSignals()

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
			color.Green("build succeeded. took: %s\n", took)
		} else {
			color.Red("build failed. took: %s\n", took)
			color.Magenta(stdout.String())
			color.Red(stderr.String())
			if err != nil {
				color.RedString(err.Error())
			}
			return
		}
	}
	if config.Only != 0 && config.Only > len(examples) {
		panic(fmt.Errorf("have %d test cases, wanted to run case number %d", len(examples), config.Only))
	}
	for i, el := range examples {
		if config.Only != 0 && config.Only != i+1 {
			continue
		}
		if !config.StrictEllipsis && strings.HasSuffix(strings.TrimSpace(el.Input), "...") {
			color.White("case %d skipped due to incomplete input.\n", i+1)
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
			g := color.New(color.FgGreen)
			g.Printf("case %d completed successfully. took: %s", i+1, took)
		} else {
			ru := color.New(color.FgRed, color.Underline)
			ru.Printf("case %d failed. took: %s", i+1, took)
			if err != nil {
				rb := color.New(color.FgRed, color.Underline)
				rb.Printf("returned error: %s", err.Error())
			}
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
			color.Red(out)
		}
		if failed && config.ExitOnFailure {
			return
		}
	}
}

func (e *Example) String() string {
	ret := fmt.Sprintf("input\n%s", color.MagentaString(e.Input))
	if !e.noOut {
		ret += fmt.Sprintf("\nexpected\n%s", color.BlueString(e.Output))
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
