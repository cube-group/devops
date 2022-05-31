package task

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestRunnerAsync_Start(t *testing.T) {

	//开启多核
	runtime.GOMAXPROCS(runtime.NumCPU())

	//创建runner对象，设置超时时间
	runner := NewRunnerAsync(8 * time.Second)
	//添加运行的任务
	runner.Add(
		createTaskAsync(),
		createTaskAsync(),
		createTaskAsync(),
		createTaskAsync(),
		createTaskAsync(),
		createTaskAsync(),
		createTaskAsync(),
		createTaskAsync(),
		createTaskAsync(),
		createTaskAsync(),
		createTaskAsync(),
		createTaskAsync(),
		createTaskAsync(),
	)

	fmt.Println("同步执行任务")

	//开始执行任务
	if err := runner.Start(); err != nil {
		switch err {
		case ErrTimeout:
			fmt.Println("执行超时")
			os.Exit(1)
		case ErrInterrupt:
			fmt.Println("任务被中断")
			os.Exit(2)
		}
	}

	t.Log("执行结束")

}

//创建要执行的任务
func createTaskAsync() func(id int) {
	return func(id int) {
		fmt.Printf("正在执行%v个任务\n", id)
		//模拟任务执行,sleep两秒
		//time.Sleep(1 * time.Second)
	}
}
