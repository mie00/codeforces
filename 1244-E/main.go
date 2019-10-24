package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"

	"github.com/mie00/codeforces/codeforcesutils"
)

func main() {
	var n, k int64
	_, err := fmt.Scanf("%d %d", &n, &k)
	if err != nil {
		panic(err)
	}
	codeforcesutils.Signal()
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	i2 := bytes.Split(bytes.TrimSpace(in), []byte(" "))
	codeforcesutils.Signal()
	inp := make([]int, len(i2))
	for i := range i2 {
		inp[i], err = strconv.Atoi(string(i2[i]))
		if err != nil {
			panic(err)
		}
	}
	codeforcesutils.Signal()
	sort.Ints(inp)
	codeforcesutils.Signal()
	i := int64(0)
	j := n - 1
	minus := int64(0)
	for k > 0 && (k >= i+1 || k >= n-j) && i < j {
		if inp[i] == inp[i+1] {
			i++
			continue
		}
		if inp[j] == inp[j-1] {
			j--
			continue
		}
		if i+1 <= n-j {
			need := int64(inp[i+1] - inp[i])
			if need*(i+1) <= k {
				k -= need * (i + 1)
				i++
				continue
			} else {
				minus = k / (i + 1)
				k -= minus * (i + 1)
				break
			}
		} else {
			need := int64(inp[j] - inp[j-1])
			if need*(n-j) <= k {
				k -= need * (n - j)
				j--
				continue
			} else {
				minus = k / (n - j)
				k -= minus * (n - j)
				break
			}
		}
	}
	fmt.Printf("%d", inp[j]-inp[i]-int(minus))
}
