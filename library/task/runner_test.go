package task

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestRunner_Start(t *testing.T) {
	//开启多核心
	runtime.GOMAXPROCS(runtime.NumCPU())

	//创建runner对象，设置超时时间
	runner := NewRunner(18 * time.Second)
	//添加运行的任务
	runner.Add(
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
		createTask(),
	)

	fmt.Println("异步执行任务")

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

	t.Log(runner.GetErrs())

}

//创建要执行的任务
func createTask() func(id int) error {
	return func(id int) error {
		fmt.Printf("正在执行%v个任务\n", id)
		//模拟任务执行,sleep
		//time.Sleep(1 * time.Second)
		return nil
	}
}
