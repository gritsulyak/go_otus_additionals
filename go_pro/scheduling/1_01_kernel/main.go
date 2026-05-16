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

	// 1. Привязываем горутину к конкретному системному потоку (M).
	// Это гарантирует, что системный вызов заблокирует весь поток целиком.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	fmt.Printf("Горутина %d: засыпаю в ядре на %v...\n", id, duration)

	// 2. Выполняем прямой системный вызов (Syscall).
	// Процесс перейдет в состояние TASK_INTERRUPTIBLE на уровне ядра.
	tv := syscall.NsecToTimespec(duration.Nanoseconds())
	err := syscall.Nanosleep(&tv, nil)

	if err != nil {
		fmt.Printf("Ошибка в горутине %d: %v\n", id, err)
	}
	fmt.Printf("Горутина %d: проснулась!\n", id)
}

func main() {

	// Устанавливаем количество P (логических процессоров) = 2
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
