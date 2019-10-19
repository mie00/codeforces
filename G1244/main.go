package main

import (
	"fmt"
	"strconv"
	"strings"
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

func main() {
	var n, k int
	fmt.Scanf("%d %d", &n, &k)
	arr1 := make([]int, n)
	arr2 := make([]int, n)
	for i := 0; i < n; i++ {
		arr1[i] = i + 1
		arr2[i] = i + 1
	}
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
	res := -1
	res1 := make([]int, n)
	res2 := make([]int, n)
	e1 := arr1
	permutations(arr2, func(e2 []int) bool {
		sum := 0
		for i := 0; i < len(e1); i++ {
			if e1[i] < e2[i] {
				sum += e2[i]
			} else {
				sum += e1[i]
			}
		}
		if sum > res && sum <= k {
			res = sum
			copy(res1, e1)
			copy(res2, e2)
			if res == k {
				return true
			}
		}
		return false
	})
	fmt.Println(res)
	if res != -1 {
		printIntArr(res1)
		printIntArr(res2)
	}
}
