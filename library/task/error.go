package task

import "errors"

//超时错误
var ErrTimeout = errors.New("received timeout")

//操作系统系统中断错误
var ErrInterrupt = errors.New("received interrupt")
