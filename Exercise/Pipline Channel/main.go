package main

import (
	"fmt"
)

func Generator(nums ...int) <-chan int {
	out := make(chan int)

	go func() {
		for _, num := range nums {
			out <- num
		}
		close(out)
	}()

	return out
}

func Square(nums <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for num := range nums {
			out <- num * num
		}
		close(out)

	}()

	return out
}

func PrintNum(nums <-chan int) {
	for num := range nums {
		fmt.Println(num)
	}
}

func main() {
	fmt.Println("Hello World")
	gen := Generator(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	square := Square(gen)
	PrintNum(square)
}
