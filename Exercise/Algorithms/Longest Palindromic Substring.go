//Given a string s, return the longest palindromic substring in s.

//Example 1:
//
//Input: s = "babad"
//Output: "bab"
//Explanation: "aba" is also a valid answer.
//Example 2:
//
//Input: s = "cbbd"
//Output: "bb"
//
//
//Constraints:
//
//1 <= s.length <= 1000
//s consist of only digits and English letters.

package main

import "fmt"

func expand(s string, left, right int) string {
	for left >= 0 && right < len(s) && s[left] == s[right] {
		left--
		right++
	}
	return s[left+1 : right]
}

func longestPalindrome(s string) string {
	if len(s) < 2 {
		return s
	}

	longest := ""
	for i := 0; i < len(s); i++ {
		resOdd := expand(s, i, i)
		if len(resOdd) > len(longest) {
			longest = resOdd
		}

		resEven := expand(s, i, i+1)
		if len(resEven) > len(longest) {
			longest = resEven
		}
	}
	return longest
}

func main() {
	fmt.Println("Hello World")
	s := "babad"
	fmt.Println(longestPalindrome(s))
}
