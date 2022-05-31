package e

import (
	"fmt"
	"runtime"
)

func Trace() {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("%s\n", string(buf[:n]))
}

func TraceString() string{
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return fmt.Sprintf("%s\n", string(buf[:n]))
}
