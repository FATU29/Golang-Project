package main

import (
	"fmt"
)

func maxArea(heights []int) int {
	left := 0
	right := len(heights) - 1
	maxWater := 0

	for left < right {
		h := min(heights[left], heights[right])

		width := right - left
		currentWater := h * width

		if currentWater > maxWater {
			maxWater = currentWater
		}

		if heights[left] < heights[right] {
			left++
		} else {
			right--
		}
	}

	return maxWater

}

func main() {
	fmt.Println("Hello World")
	height := []int{1, 7, 2, 5, 4, 7, 3, 6}
	fmt.Println(maxArea(height))

}
