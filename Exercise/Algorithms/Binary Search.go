package main

import "fmt"

func search(nums []int, target int) int {
	i := 0
	j := len(nums) - 1

	for i <= j {
		mid := (i + j) / 2

		if nums[mid] > target {
			j = mid - 1
			continue
		}
		if nums[mid] < target {
			i = mid + 1
			continue
		}

		if nums[mid] == target {
			return mid
		}
	}

	return -1

}

func main() {
	fmt.Println("Hello World")
	nums := []int{-1, 0, 2, 4, 6, 8}
	target := 4

	fmt.Println(search(nums, target))

}
