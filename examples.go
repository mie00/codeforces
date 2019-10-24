package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func getExamples(config Config) (Examples, error) {
	var examples Examples
	var err error
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
			return nil, err
		}
		defer file.Close()

		_, err = file.Write(buf)
		if err != nil {
			return nil, err
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

	return examples, nil
}
