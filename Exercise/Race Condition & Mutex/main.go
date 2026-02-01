package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Hello World")

	count := 0

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			count++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Println(count)
	fmt.Println("End of Main")
}
