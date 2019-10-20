package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func printFile(filename, dir, lang string) {
	file, err := os.Open(dir + "/" + filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		t := scanner.Text()
		if !strings.Contains(t, "codeforcesutils") {
			fmt.Println(t)
		}
	}
}
