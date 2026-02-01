package main

import (
	"fmt"
	"time"
)

func SlowWork(c chan<- string) {
	time.Sleep(2 * time.Second)
	c <- "Received Result SlowWork"
}

func FastWork(c chan<- string) {
	time.Sleep(1 * time.Second)
	c <- "Received Result FastWork"
}

func main() {
	fmt.Println("Hello World")

	c := make(chan string)

	go SlowWork(c)
	go FastWork(c)

	select {
	case msg := <-c:
		fmt.Println(msg)
	case msg2 := <-c:
		fmt.Println(msg2)
	}

	fmt.Println("End of Main")

}
