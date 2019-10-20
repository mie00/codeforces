package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/mie00/codeforces/codeforcesutils"
)

func main() {
	var n int
	_, err := fmt.Scanf("%d", &n)
	if err != nil {
		panic(err)
	}
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	codeforcesutils.Signal()
	inp := strings.Split(strings.Replace(string(bytes.TrimSpace(in)), "\r", "", -1), "\n")
	codeforcesutils.Signal()
	for i := range inp {
		cc := strings.Split(inp[i], " ")
		a, _ := strconv.Atoi(cc[0])
		b, _ := strconv.Atoi(cc[1])
		c, _ := strconv.Atoi(cc[2])
		d, _ := strconv.Atoi(cc[3])
		k, _ := strconv.Atoi(cc[4])
		x := (c + a - 1) / c
		y := (d + b - 1) / d
		if x+y <= k {
			fmt.Printf("%d %d\n", k-y, y)
		} else {
			fmt.Println("-1")
		}
	}
}
