//Given a string s, find the length of the longest substring without duplicate characters.
//
//
//
//Example 1:
//
//Input: s = "abcabcbb"
//Output: 3
//Explanation: The answer is "abc", with the length of 3. Note that "bca" and "cab" are also correct answers.
//Example 2:
//
//Input: s = "bbbbb"
//Output: 1
//Explanation: The answer is "b", with the length of 1.
//Example 3:
//
//Input: s = "pwwkew"
//Output: 3
//Explanation: The answer is "wke", with the length of 3.
//Notice that the answer must be a substring, "pwke" is a subsequence and not a substring.
//
//
//Constraints:
//
//0 <= s.length <= 5 * 104
//s consists of English letters, digits, symbols and spaces.

package main

import (
	"fmt"
	"strings"
)

func isContain(s string, r rune) bool {
	return strings.ContainsRune(s, r)
}

func lengthOfLongestSubstring(s string) int {
	var tmpS string
	maxLen := 0
	for _, value := range s {
		for isContain(tmpS, value) {
			tmpS = tmpS[1:]
		}

		tmpS += string(value)

		if len(tmpS) > maxLen {
			maxLen = len(tmpS)
		}
	}

	return maxLen
}

func lengthOfLongestSubstringV2(s string) int {
	savedMap := make(map[rune]int)
	maxLen := 0
	left := 0

	for right, r := range s {
		if val, ok := savedMap[r]; ok && val >= left {
			left = val + 1
		}

		savedMap[r] = right
		tmp := right - left + 1

		if tmp > maxLen {
			maxLen = tmp
		}
	}

	return maxLen
}

func main() {
	fmt.Println("Hello World")
	s := "pwwkew"

	fmt.Println(lengthOfLongestSubstringV2(s))
}
