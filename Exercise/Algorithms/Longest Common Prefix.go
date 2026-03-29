//
//Write a function to find the longest common prefix string amongst an array of strings.
//
//If there is no common prefix, return an empty string "".
//
//
//
//Example 1:
//
//Input: strs = ["flower","flow","flight"]
//Output: "fl"
//Example 2:
//
//Input: strs = ["dog","racecar","car"]
//Output: ""
//Explanation: There is no common prefix among the input strings.

package main

import (
	"fmt"
	"math"
)

func longestCommonPrefix(strs []string) string {
	minLen := math.MaxInt32
	idxSkip := 0

	for idx, str := range strs {
		if len(str) < minLen {
			minLen = len(str)
			idxSkip = idx
		}
	}

	for i := 0; i < minLen; i++ {
		for j, str := range strs {
			if idxSkip == j {
				continue
			}

			if str[i] != strs[idxSkip][i] {
				return strs[idxSkip][:i]
			}
		}
	}
	return strs[idxSkip]
}

func main() {
	fmt.Println("Hello World")
	strs := []string{"flower", "flow", "flight"}

	fmt.Println(longestCommonPrefix(strs))
}
