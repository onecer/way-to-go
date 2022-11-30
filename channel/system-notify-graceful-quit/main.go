package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func cleanUp(closed chan struct{}) {
	rand.Seed(time.Now().UnixNano())
	randCleanUpTime := rand.Intn(5) + 1
	fmt.Println("clean up start")
	time.Sleep(time.Second * time.Duration(randCleanUpTime))
	close(closed)
}

func main() {
	closing := make(chan struct{})
	closed := make(chan struct{})
	chSignal := make(chan os.Signal)
	go func() {
		i := 0
		for {
			select {
			case <-closing:
				fmt.Println("receive closing.")
				return
			default:
				i += 1
				fmt.Printf("Running %d secs.\n", i)
				time.Sleep(time.Second)
			}
		}
	}()

	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	<-chSignal
	close(closing)
	go cleanUp(closed)
	select {
	case <-time.After(time.Second * 4):
		fmt.Println("clean up timeout, force exit")
	case <-closed:
		fmt.Println("clean up done")
	}
	fmt.Println("graceful exit")
}
