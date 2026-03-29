package main

import (
	"fmt"
	"strconv"
)

func evalRPN(tokens []string) int {
	var stack []int
	result := 0

	if len(tokens) == 1 {
		val, _ := strconv.Atoi(tokens[0])
		return val
	}

	for _, token := range tokens {
		if val, err := strconv.Atoi(token); err == nil {
			stack = append(stack, val)
		} else {
			firstNumber := stack[len(stack)-1]
			secondNumber := stack[len(stack)-2]
			if token == "+" {
				result = firstNumber + secondNumber
			} else if token == "-" {
				result = secondNumber - firstNumber
			} else if token == "*" {
				result = firstNumber * secondNumber
			} else if token == "/" {
				result = secondNumber / firstNumber
			}

			stack = stack[0 : len(stack)-2]
			stack = append(stack, result)
		}
	}
	return result
}

func main() {
	fmt.Println("Hello World")
	tokens := []string{"1", "2", "+", "3", "*", "4", "-"}
	fmt.Println(evalRPN(tokens))
}
