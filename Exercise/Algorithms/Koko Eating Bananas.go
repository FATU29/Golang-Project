package main

import (
	"fmt"
	"math"
)

func minEatingSpeed(piles []int, h int) int {
	maxK := 0

	for i := 0; i < len(piles); i++ {
		if piles[i] > maxK {
			maxK = piles[i]
		}
	}

	left := 1
	right := maxK

	for left < right {
		mid := (right + left) / 2

		flag := false
		count := 0

		for _, pile := range piles {
			count += int(math.Ceil(float64(pile) / float64(mid)))
		}

		if count <= h {
			flag = true
		}

		if flag {
			right = mid
		} else {
			left = mid + 1
		}
	}

	return left
}

func main() {
	fmt.Println("Hello World")

	piles := []int{1, 4, 3, 2}
	h := 9

	fmt.Println(minEatingSpeed(piles, h))
}
