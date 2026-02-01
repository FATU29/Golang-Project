package main

import (
	"fmt"
	"sync"
)

func DoJobs(jobs <-chan int, res chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for value := range jobs {
		res <- value
		fmt.Printf("Doing Job %d \n", value)
	}
}

func main() {
	fmt.Println("Hello World")
	jobsLength := 10
	maxWorker := 3

	jobs := make(chan int, jobsLength)
	results := make(chan int, jobsLength)

	var wg sync.WaitGroup

	for i := 0; i < maxWorker; i++ {
		wg.Add(1)
		go DoJobs(jobs, results, &wg)
	}

	for i := 0; i <= jobsLength; i++ {
		jobs <- i
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for value := range results {
		fmt.Printf("Result: %d\n", value)
	}

	fmt.Println("Done")

}
