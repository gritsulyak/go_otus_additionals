// main.go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int) // небуферизированный канал

	// Читатель
	go func() {
		val := <-ch
		fmt.Printf("Reader received: %d\n", val)
	}()

	// Писатель (с задержкой)
	go func() {
		time.Sleep(2 * time.Second)
		ch <- 42
		fmt.Println("Writer sent 42")
	}()

	fmt.Println("Main waiting...")
	time.Sleep(3 * time.Second) // даём горутинам время выполниться
	fmt.Println("Done")
}

