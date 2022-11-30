package main

import (
	"fmt"
)

func product(ch chan int) {
	for i := 0; i < 100; i++ {
		ch <- i
	}
	close(ch)
}

func consumer(ch chan int) {
	for {
		v, ok := <-ch
		if ok {
			fmt.Println(v)
		}
	}
}

func main() {
	ch := make(chan int)
	go product(ch)
	go consumer(ch)
	go consumer(ch)
	select {}
}
