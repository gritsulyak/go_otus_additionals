// main.go
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup

	// Запускаем 100 горутин, каждая делает небольшую работу
	start := time.Now()
	for i := 0; i < 100_000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// "Тяжёлые" вычисления – имитируем нагрузку
			sum := 0
			for j := 0; j < 1_000_000; j++ {
				sum += j
			}
			// Выводим только первые и последние несколько
			if id < 5 || id >= 95 {
				fmt.Printf("Goroutine %d finished, sum=%d\n", id, sum)
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("All done in %v\n", time.Since(start))
}

