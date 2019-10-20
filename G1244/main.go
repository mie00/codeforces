package main

import (
	"fmt"
	"math"
	"strconv"
)

func len1N(n int) int {
	if n <= 0 {
		return 0
	}
	l := int(math.Log10(float64(n)))
	nn := int(math.Pow10(l) - 1)
	return len1N(nn) + (n-nn)*(l+2)
}

func main() {
	var n, k int64
	fmt.Scanf("%d %d", &n, &k)
	min := n * (n + 1) / 2
	var max int64
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
	res := make([]byte, 0, len1N(int(n)))
	for i := int64(1); i <= n; i++ {
		res = strconv.AppendInt(res, i, 10)
		res = append(res, ' ')
	}
	fmt.Printf("%s\n", res)
	res = res[:0]
	diff := k - min
	mmax := n
	rem := map[int64]int64{}
	for i := int64(1); i <= n; i++ {
		if i >= mmax || mmax-i > diff {
			if nnn, ok := rem[i]; ok {
				res = strconv.AppendInt(res, nnn, 10)
				res = append(res, ' ')
			} else {
				res = strconv.AppendInt(res, i, 10)
				res = append(res, ' ')
			}
		} else {
			res = strconv.AppendInt(res, mmax, 10)
			res = append(res, ' ')
			diff -= mmax - i
			rem[mmax] = i
			mmax--
		}
	}
	fmt.Printf("%s\n", res)
	if diff != 0 {
		panic("Programming error")
	}
}
