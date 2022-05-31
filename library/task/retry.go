package task

import (
	"fmt"
	"time"
)

//重试函数
func Retry(msg string, tryTimes int, sleep time.Duration, callback func() error) (err error) {
	for i := 1; i <= tryTimes; i++ {
		err = callback()
		if err == nil {
			return nil
		}
		fmt.Printf("%v失败，第%v次重试， 错误信息:%s \n", msg, i, err)
		time.Sleep(sleep)
	}
	err = fmt.Errorf("%v失败，共重试%d次, 最近一次错误:\n %s", msg, tryTimes, err)
	fmt.Println(err)
	return err
}

//
func TetryDuring(duration time.Duration, sleep time.Duration, callback func() error) (err error) {
	t0 := time.Now()
	i := 0
	for {
		i++

		err = callback()
		if err == nil {
			return
		}

		delta := time.Now().Sub(t0)
		if delta > duration {
			return fmt.Errorf("after %d times (during %s), last error: %s", i, delta, err)
		}

		time.Sleep(sleep)
		fmt.Println("retrying after error:", err)
	}
}
