package task

import (
	"os"
	"os/signal"
	"time"
)

//同步执行任务
type RunnerAsync struct {
	//操作系统的信号检测
	interrupt chan os.Signal

	//记录执行完成的状态
	complete chan error

	//超时检测
	timeout <-chan time.Time

	//保存所有要执行的任务,顺序执行
	tasks []func(id int)
}

//new一个RunnerAsync对象
func NewRunnerAsync(d time.Duration) *RunnerAsync {
	return &RunnerAsync{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
	}
}

//添加一个任务
func (this *RunnerAsync) Add(tasks ...func(id int)) {
	this.tasks = append(this.tasks, tasks...)
}

//启动RunnerAsync，监听错误信息
func (this *RunnerAsync) Start() error {

	//接收操作系统信号
	signal.Notify(this.interrupt, os.Interrupt)

	//执行任务
	go func() {
		this.complete <- this.Run()
	}()

	select {
	//返回执行结果
	case err := <-this.complete:
		return err

		//超时返回
	case <-this.timeout:
		return ErrTimeout
	}
}

//顺序执行所有的任务
func (this *RunnerAsync) Run() error {
	for id, task := range this.tasks {
		if this.gotInterrupt() {
			return ErrInterrupt
		}
		//执行任务
		task(id)
	}
	return nil
}

//判断是否接收到操作系统中断信号
func (this *RunnerAsync) gotInterrupt() bool {
	select {
	case <-this.interrupt:
		//停止接收别的信号
		signal.Stop(this.interrupt)
		return true
		//正常执行
	default:
		return false
	}
}
