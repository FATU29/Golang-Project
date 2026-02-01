package main

import (
	"fmt"
	"sync"
)

func Ping(pingChan <-chan string, pongChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for val := range pingChan {
		fmt.Printf("Ping -> %s\n", val)
		pongChan <- "Ping"
	}
	close(pongChan)
}

func Pong(pingChan chan<- string, pongChan <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0

	for range pongChan {

		count++
		if count == 5 {
			fmt.Println("🔵 Pong: Đủ 5 lần rồi, đóng pingChan!")
			close(pingChan)
			return
		}

		pingChan <- "Pong"

	}

}

func main() {
	fmt.Println("Hello World")
	pingChan := make(chan string)
	pongChan := make(chan string)

	var wg sync.WaitGroup

	wg.Add(2)

	go Ping(pingChan, pongChan, &wg)
	go Pong(pingChan, pongChan, &wg)

	pingChan <- "Pong"
	wg.Wait()

	fmt.Println("End Of Main")

}
