package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Download(url string, res chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	duration := time.Duration(rand.Intn(3)) * time.Second
	time.Sleep(duration)
	fmt.Println("Downloading", url, "to", duration)
	res <- url
}

func main() {
	fmt.Println("Hello World")

	urls := []string{"movie.mp4", "music.mp3", "document.pdf"}

	res := make(chan string, 3)

	var wg sync.WaitGroup

	for i := 0; i < len(urls); i++ {
		wg.Add(1)
		go Download(urls[i], res, &wg)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	for val := range res {
		fmt.Printf("Downloaded %s \n", val)
	}

}
