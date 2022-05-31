package util

import "time"

func TaskRetry(times uint, dt time.Duration, f func() error) error {
	var try uint = 0
	for {
		try++
		err := f();
		if err == nil {
			return nil
		}
		if try >= times {
			return err
		}
		time.Sleep(dt)
	}
}
