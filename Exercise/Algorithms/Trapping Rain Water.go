package main

import (
	"fmt"
	"math"
)

func trap(height []int) int {
	left := 0
	var res float64
	var plus float64

	for i := 0; i < len(height); i++ {
		if height[i] == 0 {
			continue
		}

		plus = 0
		findBlockBack := 0

		for j := left + 1; j < i; j++ {
			findBlockBack += height[j]
		}

		minVal := min(height[left], height[i])
		plus = float64(minVal)*math.Abs(float64(i-left-1)) - float64(findBlockBack)

		if height[i] >= height[left] {
			left = i
			res += plus
		}

	}

	return int(res)
}

func main() {
	fmt.Println("Hello World")
	height := []int{0, 2, 0, 3, 1, 0, 1, 3, 2, 1}
	fmt.Println(trap(height))
}
