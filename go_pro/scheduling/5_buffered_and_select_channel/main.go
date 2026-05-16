// main.go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Буферизированный канал на 2 элемента
	ch := make(chan string, 2)
	ch <- "first"
	ch <- "second"
	// Если раскомментировать следующую строку – программа зависнет (канал полон)
	// ch <- "third"

	// Неблокирующее чтение с select
	select {
	case msg := <-ch:
		fmt.Println("Read from channel:", msg)
	default:
		fmt.Println("No data available")
	}

	// Пример таймаута
	ch2 := make(chan int)
	go func() {
		time.Sleep(2 * time.Second)
		ch2 <- 100
	}()

	select {
	case val := <-ch2:
		fmt.Println("Received:", val)
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout! No data after 1 second")
	}

	// Небуферизированный канал – писатель блокируется до появления читателя
	ch3 := make(chan bool)
	go func() {
		fmt.Println("Before send on unbuffered")
		ch3 <- true // блокируется, пока кто-то не прочитает
		fmt.Println("After send (unbuffered)")
	}()

	time.Sleep(500 * time.Millisecond)
	fmt.Println("Main ready to receive")
	<-ch3
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Done")
}
