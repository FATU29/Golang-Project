package main

import "fmt"

func findMin(nums []int) int {
	left := 0
	right := len(nums) - 1

	for left < right {
		mid := (left + right) / 2
		midNum := nums[mid]
		midRight := nums[right]

		if midNum > midRight {
			left = mid + 1
		} else {
			right = mid
		}
	}

	return nums[left]
}

func main() {
	fmt.Println("Hello World")
	nums := []int{3, 4, 5, 6, 1, 2}

	fmt.Println(findMin(nums))
}
