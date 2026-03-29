//You are given an integer array nums consisting of n elements, and an integer k.
//
//Find a contiguous subarray whose length is equal to k that has the maximum average value and return this value. Any answer with a calculation error less than 10-5 will be accepted.
//
//
//
//Example 1:
//
//Input: nums = [1,12,-5,-6,50,3], k = 4
//Output: 12.75000
//Explanation: Maximum average is (12 - 5 - 6 + 50) / 4 = 51 / 4 = 12.75
//Example 2:
//
//Input: nums = [5], k = 1
//Output: 5.00000

package main

import (
	"fmt"
	"math"
)

func findMaxAverage(nums []int, k int) float64 {
	var maxAvg float64
	maxAvg = math.MinInt16
	total := 0
	left := 0

	for idx, num := range nums {
		total += num
		if idx-left+1 == k {
			avg := float64(total) / float64(k)
			if avg > maxAvg {
				maxAvg = avg
			}
			total -= nums[left]
			left++
		}
	}
	return maxAvg

}

func main() {
	fmt.Println("Hello World")

	nums := []int{-1}
	k := 1

	fmt.Println(findMaxAverage(nums, k))
}
