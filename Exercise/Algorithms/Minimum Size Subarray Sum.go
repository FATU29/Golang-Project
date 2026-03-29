//Given an array of positive integers nums and a positive integer target, return the minimal length of a subarray whose sum is greater than or equal to target. If there is no such subarray, return 0 instead.

//Example 1:
//
//Input: target = 7, nums = [2,3,1,2,4,3]
//Output: 2
//Explanation: The subarray [4,3] has the minimal length under the problem constraint.
//Example 2:
//
//Input: target = 4, nums = [1,4,4]
//Output: 1
//Example 3:
//
//Input: target = 11, nums = [1,1,1,1,1,1,1,1]
//Output: 0

package main

import (
	"fmt"
	"math"
)

func minSubArrayLen(target int, nums []int) int {
	left := 0
	total := 0

	minLen := math.MaxInt32

	for idx, val := range nums {
		total += val

		for total >= target {
			if minLen > idx-left+1 {
				minLen = idx - left + 1
			}
			total -= nums[left]
			left++
		}
	}
	if minLen == math.MaxInt32 {
		return 0
	}
	return minLen
}

func main() {
	fmt.Println("Hello World")
	target := 7
	nums := []int{2, 3, 1, 2, 4, 3}

	fmt.Println(minSubArrayLen(target, nums))

}
