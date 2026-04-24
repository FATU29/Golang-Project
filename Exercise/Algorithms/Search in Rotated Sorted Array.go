package main

import "fmt"

func searchRotatedSortArray(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := (left + right) / 2

		numMid := nums[mid]
		numRight := nums[right]
		numLeft := nums[left]

		if target == numMid {
			return mid
		}

		if numLeft <= numMid {
			if target >= numLeft && target < numMid {
				right = mid - 1
			} else {
				left = mid + 1
			}
		} else {
			if target <= numRight && target > numMid {
				left = mid + 1
			} else {
				right = mid - 1
			}

		}

	}

	return -1
}

func main() {
	fmt.Println("Hello World")

	nums := []int{3, 4, 5, 6, 1, 2}
	target := 1

	fmt.Println(searchRotatedSortArray(nums, target))
}
