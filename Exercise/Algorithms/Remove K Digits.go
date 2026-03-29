package main

import (
	"fmt"
	"strings"
)

func removeKdigits(num string, k int) string {
	if len(num) == k {
		return "0"
	}

	stack := make([]byte, 0, len(num))

	for i := 0; i < len(num); i++ {
		for k > 0 && len(stack) > 0 && num[i] < stack[len(stack)-1] {
			stack = stack[:len(stack)-1]
			k--
		}
		stack = append(stack, num[i])
	}

	stack = stack[:len(stack)-k]
	res := string(stack)

	res = strings.TrimLeft(res, "0")
	if res == "" {
		return "0"
	}
	return res
}

func main() {
	fmt.Println("Hello World")
	num := "1432219"
	k := 3
	fmt.Println(removeKdigits(num, k))
}
