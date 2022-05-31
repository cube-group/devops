package task

import (
	"fmt"
	"testing"
	"time"
)

func TestAsync_Run(t *testing.T) {
	async := NewAsync()
	async.Run(
		func() interface{} {
			time.Sleep(1 * time.Second)
			fmt.Println("1 finished")
			return true
		},
		func() interface{} {
			time.Sleep(1 * time.Second)
			fmt.Println("2 finished")
			return false
		},
		func() interface{} {
			time.Sleep(1 * time.Second)
			fmt.Println("3 finished")
			return false
		},
	)
	fmt.Println(async.Results)
}
