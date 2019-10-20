package codeforcesutils

import (
	"fmt"
	"os"
	"syscall"
)

var last = syscall.SIGUSR2

func Signal() {
	if last == syscall.SIGUSR2 {
		last = syscall.SIGUSR1
	} else {
		last = syscall.SIGUSR2
	}
	SomeSignal(last)
}

func SomeSignal(signal os.Signal) {
	pid := os.Getppid()
	p, err := os.FindProcess(pid)
	if err == nil {
		p.Signal(signal)
	} else {
		fmt.Fprintf(os.Stderr, "unable to find parent process")
	}
}
