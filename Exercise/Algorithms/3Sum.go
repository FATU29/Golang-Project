package main

import (
	"fmt"
	"sort"
)

func threeSum(nums []int) [][]int {
	sort.Ints(nums)

	var res [][]int
	for i := 0; i < len(nums)-2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		left := i + 1
		right := len(nums) - 1

		for left < right {
			val := nums[i] + nums[left] + nums[right]

			if val > 0 {
				right--
				continue
			}

			if val < 0 {
				left++
				continue
			}

			arr := []int{nums[i], nums[left], nums[right]}
			res = append(res, arr)

			for left < right && nums[left] == nums[left+1] {
				left++
			}
			for left < right && nums[right] == nums[right-1] {
				right--
			}
			left++
			right--
		}
	}
	return res
}

func main() {
	fmt.Println("Hello World")
	nums := []int{-1, 0, 1, 2, -1, -4}
	fmt.Println(threeSum(nums))
}
