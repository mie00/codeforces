package main

import (
	"fmt"

	"github.com/mie00/codeforces/codeforcesutils"
)

type Cycle []byte

func (c Cycle) flip(i int) {
	if i < 0 {
		i += len(c)
	} else if i >= len(c) {
		i -= len(c)
	}
	if c[i] == 'W' {
		c[i] = 'B'
	} else {
		c[i] = 'W'
	}
}

func (c Cycle) flipall() {
	for i := range c {
		if c[i] == 'W' {
			c[i] = 'B'
		} else {
			c[i] = 'W'
		}
	}
}

func (c Cycle) get(n int) byte {
	if n < 0 {
		n += len(c)
	} else if n >= len(c) {
		n -= len(c)
	}
	return c[n]
}

func (c Cycle) set(n int, b byte) {
	if n < 0 {
		n += len(c)
	} else if n >= len(c) {
		n -= len(c)
	}
	c[n] = b
}

func (c Cycle) diff(i, j int) int {
	x := i - j
	for x < 0 {
		x += len(c)
	}
	for x >= len(c) {
		x -= len(c)
	}
	return x
}

func (c Cycle) print() {
	fmt.Println(string(c))
}

func main() {
	var n int
	var k int64
	var nk int
	_, err := fmt.Scanf("%d %d", &n, &k)
	if err != nil {
		panic(err)
	}
	nk = n
	if k < int64(n) {
		nk = int(k)
	}
	var s Cycle
	codeforcesutils.Signal()
	_, err = fmt.Scan(&s)
	if err != nil {
		panic(err)
	}
	codeforcesutils.Signal()
	firstBeforeStartPos := -1
	var before byte
	finalBefore := false
	firstAfterStartPos := -1
	var after byte
	startPos := -1
	alternative := true
	for i := range s {
		if startPos == -1 && s.get(i) != s.get(i-1) && s.get(i) != s.get(i+1) {
			startPos = i
			if firstBeforeStartPos != -1 {
				finalBefore = true
			}
		}
		if s.get(i) == s.get(i-1) || s.get(i) == s.get(i+1) {
			alternative = false
		}
		if s.get(i) == s.get(i+1) {
			if startPos != -1 && firstAfterStartPos == -1 {
				firstAfterStartPos = i
				after = s.get(i)
			}
		}
		if s.get(i) == s.get(i-1) {
			if !finalBefore {
				firstBeforeStartPos = i
				before = s.get(i)
			}
		}
	}
	if startPos == -1 {
		s.print()
	} else if alternative {
		if k%2 != 0 {
			s.flipall()
		}
		s.print()
	} else {
		for i := startPos; i < startPos+n; {
			if i == firstAfterStartPos {
				for i < startPos+n && (s.get(i) == s.get(i-1) || s.get(i+1) == s.get(i)) {
					i++
				}
				firstBeforeStartPos = i - 1
				before = s.get(firstBeforeStartPos)
				for j := i + 1; j <= startPos+n; j++ {
					if s.get(j) == s.get(j+1) {
						firstAfterStartPos = j
						after = s.get(firstAfterStartPos)
						break
					}
				}
			}
			if s.diff(i, firstBeforeStartPos) <= s.diff(firstAfterStartPos, i) {
				if s.diff(i, firstBeforeStartPos) <= nk {
					s.set(i, before)
				} else if nk%2 == 1 {
					s.flip(i)
				}
			} else {
				if s.diff(firstAfterStartPos, i) <= nk {
					s.set(i, after)
				} else if nk%2 == 1 {
					s.flip(i)
				}
			}
			i++
		}
		s.print()
	}
}
