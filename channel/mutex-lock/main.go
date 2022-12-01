package main

import (
	"fmt"
	"time"
)

type MutexLock struct {
	ch chan struct{}
}

func NewMutexLock() *MutexLock {
	mu := &MutexLock{ch: make(chan struct{}, 1)}
	mu.ch <- struct{}{}
	return mu
}

func (m *MutexLock) Lock() bool {
	select {
	case <-m.ch:
		return true
	default:
	}
	return false
}

func (m *MutexLock) Unlock() bool {
	select {
	case m.ch <- struct{}{}:
		return true
	default:
	}
	return false
}

func (m *MutexLock) IsLocked() bool {
	return len(m.ch) == 0
}

func echoNum(name string, mu *MutexLock) {
	i := 1
	for {
		if mu.Lock() {
			fmt.Printf("%s: %d\n", name, i)
			i += 1
			mu.Unlock()
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	m := NewMutexLock()
	go echoNum("A", m)
	go echoNum("B", m)
	select {}
}
