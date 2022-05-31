package task

import (
	"os"
	"os/signal"
	"sync"
	"time"
)

//异步执行任务
type Runner struct {
	//操作系统的信号检测
	interrupt chan os.Signal

	//记录执行完成的状态
	complete chan error

	//超时检测
	timeout <-chan time.Time

	//保存所有要执行的任务,顺序执行
	tasks []func(id int) error

	waitGroup sync.WaitGroup

	lock sync.Mutex

	errs []error
}

//new一个Runner对象
func NewRunner(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
		waitGroup: sync.WaitGroup{},
		lock:      sync.Mutex{},
	}
}

//添加一个任务
func (this *Runner) Add(tasks ...func(id int) error) {
	this.tasks = append(this.tasks, tasks...)
}

//启动Runner，监听错误信息
func (this *Runner) Start() error {

	//接收操作系统信号
	signal.Notify(this.interrupt, os.Interrupt)

	//并发执行任务
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

//异步执行所有的任务
func (this *Runner) Run() error {
	for id, task := range this.tasks {
		if this.gotInterrupt() {
			return ErrInterrupt
		}

		this.waitGroup.Add(1)
		go func(id int) {
			this.lock.Lock()

			//执行任务
			err := task(id)
			//加锁保存到结果集中
			this.errs = append(this.errs, err)

			this.lock.Unlock()
			this.waitGroup.Done()
		}(id)
	}
	this.waitGroup.Wait()

	return nil
}

//判断是否接收到操作系统中断信号
func (this *Runner) gotInterrupt() bool {
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

//获取执行完的error
func (this *Runner) GetErrs() []error {
	return this.errs
}
