package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func permutations(a []int, fn func([]int) bool) {
	if len(a) == 0 {
		return
	}
	generate(len(a), a, fn)
}

func permutationsClose(a []int, fn func([]int) bool, done func()) {
	permutations(a, fn)
	done()
}

func generate(k int, a []int, fn func([]int) bool) bool {
	if k == 1 {
		return fn(a)
	}
	done := generate(k-1, a, fn)
	if done {
		return true
	}
	for i := 0; i < k-1; i++ {
		if k%2 == 0 {
			a[i], a[k-1] = a[k-1], a[i]
		} else {
			a[0], a[k-1] = a[k-1], a[0]
		}
		done = generate(k-1, a, fn)
		if done {
			return true
		}
	}
	return false
}

func intToStrArr(inp []int) []string {
	res := make([]string, len(inp))
	for i := 0; i < len(inp); i++ {
		res[i] = strconv.Itoa(inp[i])
	}
	return res
}

func printStrArr(inp []string) {
	fmt.Println(strings.Join(inp, " "))
}

func printIntArr(inp []int) {
	fmt.Println(strings.Join(intToStrArr(inp), " "))
}

func len1N(n int) int {
	if n <= 0 {
		return 0
	}
	l := int(math.Log10(float64(n)))
	nn := int(math.Pow10(l) - 1)
	return len1N(nn) + (n-nn)*(l+2)
}

func main() {
	var n, k int
	fmt.Scanf("%d %d", &n, &k)
	min := n * (n + 1) / 2
	var max int
	if n%2 == 0 {
		max = (3*n + 2) * n / 4
	} else {
		max = (3*n - 1) * (n + 1) / 4
	}
	if k < min {
		fmt.Println("-1")
		return
	}
	if k > max {
		k = max
	}
	fmt.Println(k)
	res := make([]byte, 0, len1N(n))
	for i := 1; i <= n; i++ {
		res = append(res, []byte(fmt.Sprintf("%d ", i))...)
	}
	spew.Dump(len(res), len1N(n))
	fmt.Printf("%s\n", res)
	res = res[:0]
	diff := k - min
	mmax := n
	rem := map[int]int{}
	for i := 1; i <= n; i++ {
		if i >= mmax || mmax-i > diff {
			if nnn, ok := rem[i]; ok {
				res = append(res, []byte(fmt.Sprintf("%d ", nnn))...)
			} else {
				res = append(res, []byte(fmt.Sprintf("%d ", i))...)
			}
		} else {
			res = append(res, []byte(fmt.Sprintf("%d ", mmax))...)
			diff -= mmax - i
			rem[mmax] = i
			mmax--
		}
	}
	if diff != 0 {
		panic("Programming error")
	}
}
