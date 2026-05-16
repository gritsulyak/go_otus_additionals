// main.go
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1) // один процессор, чтобы переключение было заметнее

	var wg sync.WaitGroup
	done := make(chan bool)

	// Горутина, которая "кооперативно" уступает
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			fmt.Println("Cooperative: iteration", i)
			time.Sleep(100 * time.Millisecond) // Sleep вызывает переключение
		}
		done <- true
	}()

	// Горутина, которая "жадная" – почти не уступает
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Бесконечный цикл, но с небольшим вытеснением
		for {
			select {
			case <-done:
				fmt.Println("Greedy goroutine exiting")
				return
			default:
				fmt.Println("defualt greed")
				// без Gosched() она бы почти не отдавала CPU
				runtime.Gosched() // кооперативно уступаем
			}
		}
	}()

	wg.Wait()
	fmt.Println("Finished")
}
