package e

import (
	"errors"
	"fmt"
)

func TryCatch(f func()) error {
	var err error
	func() {
		defer func() {
			if e := recover(); e != nil {
				err = errors.New(fmt.Sprintf("%v", e))
			}
		}()
		f()
	}()
	return err
}
