// main.go
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	// set P (logical processors) = 2
	runtime.GOMAXPROCS(2)
	fmt.Printf("GOMAXPROCS = %d\n", runtime.GOMAXPROCS(0))

	var wg sync.WaitGroup
	start := time.Now()

	// Run 10 goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// each goroutine 1 
			fmt.Printf("[%d] started at %v\n", id, time.Since(start))
			time.Sleep(time.Second)
			fmt.Printf("[%d] finished at %v\n", id, time.Since(start))
		}(i)
	}

	wg.Wait()
	fmt.Println("All done. Total time:", time.Since(start))
}
