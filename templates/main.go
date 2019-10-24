package main

import (
	"fmt"
	"io/ioutil"
	"os"

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
	fmt.Println(n)
	fmt.Printf(string(in))
}
