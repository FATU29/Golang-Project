//You are given an array of integers temperatures where temperatures[i] represents the daily temperatures on the ith day.
//
//Return an array result where result[i] is the number of days after the ith day
//before a warmer temperature appears on a future day. If there is no day in the future
//where a warmer temperature will appear for the ith day, set result[i] to 0 instead.

package main

import "fmt"

func dailyTemperatures(temperatures []int) []int {
	var stack []int
	res := make([]int, len(temperatures))

	for idx, temp := range temperatures {
		for len(stack) > 0 && temp > temperatures[stack[len(stack)-1]] {
			val := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			res[val] = idx - val
		}
		stack = append(stack, idx)
	}

	return res

}

func main() {
	fmt.Println("Hello World")

	temperatures := []int{30, 38, 30, 36, 35, 40, 28} //Output: [1,4,1,2,1,0,0]
	fmt.Println(dailyTemperatures(temperatures))

}
