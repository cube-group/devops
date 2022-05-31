package g

import (
	"app/library/log"
	"time"
)

type CatchFunc func()

func CatchRun(f CatchFunc, dt time.Duration) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				log.StdFatal("CatchRun", e)
			}
		}()

		if dt > 0 {
			for {
				f()
				time.Sleep(dt)
			}
		} else {
			f()
		}
	}()
}
