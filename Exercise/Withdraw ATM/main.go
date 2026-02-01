package main

import (
	"fmt"
	"sync"
)

var (
	balance = 1000
	mu      sync.Mutex
)

func Withdraw(name string, amount int, wg *sync.WaitGroup) {

	mu.Lock()

	if balance >= amount {
		balance -= amount
		fmt.Printf("%s withdrawn %d\n", name, amount)
	} else {
		fmt.Printf("%s withdrawn Error %d\n", name, amount)
	}

	defer func() {
		mu.Unlock()
		wg.Done()
	}()

}

func main() {
	fmt.Println("Hello World")
	var wg sync.WaitGroup

	wg.Add(2)

	go Withdraw("Wife", 700, &wg)
	go Withdraw("Husband", 500, &wg)

	wg.Wait()

	fmt.Println("Done")
}
