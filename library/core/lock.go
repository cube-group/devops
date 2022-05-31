package core

import "sync"

func Lock(f func()) {
	var locker sync.Mutex
	locker.Lock()
	defer locker.Unlock()
	f()
}
