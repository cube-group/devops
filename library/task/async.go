package task

import (
	"sync"
)

type AsyncFunc func() interface{}

type Async struct {
	locker  sync.Mutex
	Results []interface{}
}

func NewAsync() *Async {
	a := new(Async)
	return a
}

func (this *Async) Run(tasks ...AsyncFunc) {
	this.Results = make([]interface{}, len(tasks))
	wg := sync.WaitGroup{}
	for i, _ := range tasks {
		wg.Add(1)
		go func(index int) {
			defer func() {
				wg.Done()
				if e := recover(); e != nil {
					//pass catch err
				}
			}()

			this.locker.Lock()
			this.Results[index] = tasks[index]()
			this.locker.Unlock()
		}(i)
	}
	wg.Wait()
}
