package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Hello world")
	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		for i := 0; i < 5; i++ {
			if i == 4 {
				fmt.Println(i + 1)
				break
			}
			fmt.Printf("%d -> ", i+1)
		}
		defer wg.Done()
	}()

	wg.Wait()

	fmt.Println("End Of Main")

}
