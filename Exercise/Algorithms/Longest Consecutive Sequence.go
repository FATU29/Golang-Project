package main

import (
	"fmt"
	"sort"
)

func longestConsecutive(nums []int) int {
	size := len(nums)
	hashSet := make(map[int]struct{})

	for _, num := range nums {
		hashSet[num] = struct{}{}
	}

	longest := 0

	for i := 0; i < size; i++ {
		if _, ok := hashSet[nums[i]-1]; !ok {
			length := 1
			for {
				if _, found := hashSet[nums[i]+length]; found {
					length++
				} else {
					break
				}
			}
			if length > longest {
				longest = length
			}
		}
	}

	return longest
}

func longestConsecutiveV2(nums []int) int {
	size := len(nums)

	if size < 1 {
		return 0
	}

	sort.Ints(nums)

	longest := 1
	currentWindow := 1

	for i := 1; i < size; i++ {
		if nums[i] == nums[i-1] {
			continue
		}

		if nums[i]-nums[i-1] == 1 {
			currentWindow++
		} else {
			if currentWindow > longest {
				longest = currentWindow
			}
			currentWindow = 1
		}
	}

	if currentWindow > longest {
		longest = currentWindow
	}

	return longest

}

func main() {
	fmt.Println("Hello World")
	nums := []int{2, 20, 4, 10, 3, 4, 5}
	fmt.Println(longestConsecutive(nums))

}
