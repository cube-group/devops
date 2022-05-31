package ctx

import (
	"app/library/uuid"
	"context"
	"time"
)

type RoutineContext struct {
	Id      string
	Cancel  context.Context
	Timeout context.Context

	cancelFunc        context.CancelFunc
	cancelTimeoutFunc context.CancelFunc
}

//创建一个协程关联实例
//timeout 超时秒数
func New(timeout time.Duration) *RoutineContext {
	cancelCtx, cancelFunc := context.WithCancel(context.TODO())
	timeoutCtx, cancelTimeoutFunc := context.WithTimeout(context.TODO(), timeout)
	return &RoutineContext{
		uuid.GetUUID(),
		cancelCtx,
		timeoutCtx,
		cancelFunc,
		cancelTimeoutFunc,
	}
}

//是否被删除
func (t *RoutineContext) IsCancel() <-chan struct{} {
	return t.Cancel.Done()
}

//是否超时
func (t *RoutineContext) IsTimeout() <-chan struct{} {
	return t.Timeout.Done()
}

//强制执行关闭
func (t *RoutineContext) Close() {
	t.cancelFunc()
	t.cancelTimeoutFunc()
}
