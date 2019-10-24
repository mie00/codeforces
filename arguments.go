package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
)

var file = map[string]string{
	"go":         "main.go",
	"python":     "main.py",
	"javascript": "index.js",
	"c":          "main.c",
	"c++":        "main.cpp",
}

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
	if strings.Contains(name, "/") {
		_, err = fmt.Sscanf(name, "%d/%s", &num, &class)
	} else {
		class = name[0:1]
		num, err = strconv.Atoi(name[1:])
		if err != nil {
			return "", nil
		}
	}
	return fmt.Sprintf("%d-%s", num, class), nil
}

type Config struct {
	Run              bool `docopt:"run"`
	MatchFirstLine   bool
	Cmd              string
	BuildCmd         string
	Stdin            bool
	StdinOne         bool
	Timeout          string
	timeout          time.Duration
	BuildTimeout     string
	buildTimeout     time.Duration
	ForceDownload    bool
	Lang             string
	ExitOnFailure    bool
	Verbose          bool
	Quite            bool
	StrictEllipsis   bool
	Only             int
	NoPercentage     bool
	NoWindowsNewline bool
	New              bool `docopt:"new"`
	ForceWrite       bool
	Filename         string
	Show             bool   `docopt:"show"`
	Examples         bool   `docopt:"examples"`
	ListLangs        bool   `docopt:"list-langs"`
	Name             string `docopt:"<name>"`
}

func arguments() (*Config, error) {
	usage := `codeforces test runner.

Usage:
  codeforces run <name> [--match-first-line] [--cmd=<cmd>] [--build-cmd=<cmd>] [--stdin] [--stdin-one] [--timeout=<timeout>] [--build-timeout=<timeout>] [--force-download] [--lang=<lang>] [--exit-on-failure] [--verbose] [--quite] [--strict-ellipsis] [--only=<n>] [--no-percentage] [--no-windows-newline]
  codeforces new <name> [--force-write] [--lang=<lang>] [--filename=<fname>] [--force-download]
  codeforces show <name> [--lang=<lang>] [--filename=<fname>]
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
  --filename=<fname>                                   Default file name to show [default: ].
  --stdin                                              Get examples from stdin [default: false].
  --stdin-one                                          Get a single input from stdin. [default: false].
  --timeout=<timeout>                                  Timeout for a single case [default: 1s].
  --build-timeout=<timeout>                            Timeout for build [default: 10s].
  --force-download                                     Force download examples [default: false]
  --exit-on-failure                                    Exit on the first failed example [default: false].
  --verbose                                            Always show input/expected/output [default: false].
  --quite                                              Never show input/expected/output [default: false].
  --strict-ellipsis                                    Treat ellipsis (...) in output as is [default: false].
  --only=<n>                                           run only a specific test case, 0 means all [default: 0].
  --no-percentage                                      Show total time instead of percentage for steps instead of time [default: false].
  --no-windows-newline                                 For input do not use windows' newline (\r\n) and use (\n) instead [default: false].
  --force-write                                        Force overwrite the target file if exists [default: false]
`

	arguments, err := docopt.ParseDoc(usage)
	if err != nil {
		panic(err)
	}
	var config Config
	err = arguments.Bind(&config)
	if err != nil {
		return nil, err
	}

	config.Name, err = normalizeName(config.Name)
	if err != nil {
		return &config, err
	}
	lang, ok := languages[config.Lang]
	if !ok {
		return &config, fmt.Errorf("unknown language %s", config.Lang)
	}
	config.Lang = lang

	if config.Cmd == "" {
		langRunCmd, ok := run[lang]
		if !ok {
			return &config, fmt.Errorf("unknown language run %s", lang)
		}
		config.Cmd = langRunCmd

	}

	if langBuildCmd, ok := build[lang]; ok && config.BuildCmd == "" {
		config.BuildCmd = langBuildCmd
	}

	if config.Stdin && config.StdinOne {
		return &config, fmt.Errorf("cannot receive --stdin and --stdin-one at the same time")
	}

	config.timeout, err = time.ParseDuration(config.Timeout)
	if err != nil {
		panic(err)
	}
	config.buildTimeout, err = time.ParseDuration(config.BuildTimeout)
	if err != nil {
		panic(err)
	}
	if fn, ok := file[lang]; ok && config.Filename == "" {
		config.Filename = fn
	}

	if config.Filename == "" && (config.Show || config.New) {
		panic(fmt.Errorf("cannot find filename for lang %s, please provide one", config.Lang))
	}

	return &config, nil
}
