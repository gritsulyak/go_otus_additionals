// main.go
package main

import (
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func kernelSleep(id int, wg *sync.WaitGroup, duration time.Duration) {

	defer wg.Done()

	// lock on thread level
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	fmt.Printf("go %d: sleep in kernel for %v...\n", id, duration)

	// TASK_INTERRUPTIBLE in kernel.
	tv := syscall.NsecToTimespec(duration.Nanoseconds())
	err := syscall.Nanosleep(&tv, nil)

	if err != nil {
		fmt.Printf("Error in goroutine %d: %v\n", id, err)
	}
	fmt.Printf("Goroutine %d: awake !\n", id)

}

func main() {

	// P (logical procs) = 2
	runtime.GOMAXPROCS(2)
	fmt.Printf("GOMAXPROCS = %d\n", runtime.GOMAXPROCS(0))

	var wg sync.WaitGroup
	numCoroutines := 10
	start := time.Now()
	for i := 1; i <= numCoroutines; i++ {
		wg.Add(1)
		go kernelSleep(i, &wg, 1*time.Second)
	}

	wg.Wait()
	fmt.Println("All done. Total time:", time.Since(start))
}
