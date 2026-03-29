//Given an array of integers numbers that is sorted in non-decreasing order.
//
//Return the indices (1-indexed) of two numbers, [index1, index2], such that they add up to a given target number target and index1 < index2. Note that index1 and index2 cannot be equal, therefore you may not use the same element twice.
//
//There will always be exactly one valid solution.
//
//Your solution must use
//O
//(
//1
//)
//O(1) additional space.
//
//Example 1:
//
//Input: numbers = [1,2,3,4], target = 3
//
//Output: [1,2]
//Explanation:
//The sum of 1 and 2 is 3. Since we are assuming a 1-indexed array, index1 = 1, index2 = 2. We return [1, 2].

package main

import "fmt"

func twoSum(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1

	for left < right {
		res := numbers[left] + numbers[right]

		if res > target {
			right--
			continue
		}

		if res < target {
			left++
			continue
		}

		if res == target {
			return []int{left + 1, right + 1}
		}

	}

	return []int{}
}

func main() {
	fmt.Println("Hello World")
	nums := []int{1, 2, 3, 4}
	target := 3
	fmt.Println(twoSum(nums, target))
}
