package main

import (
	"fmt"
	"sync"
)

func DownloadFile(fileNum int, wg *sync.WaitGroup) {
	sem := make(chan int, 10)

	for i := 0; i < fileNum; i++ {
		wg.Add(1)
		go func(val int) {
			sem <- 1
			fmt.Printf("Downloading %d\n", val)
			defer func() {
				<-sem
				wg.Done()
			}()
		}(i + 1)
	}
}

func main() {
	fmt.Println("Hello World")
	var wg sync.WaitGroup

	DownloadFile(1000, &wg)

	wg.Wait()
	fmt.Println("Done")
}
