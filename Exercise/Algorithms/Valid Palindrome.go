// A phrase is a palindrome if, after converting all uppercase letters into lowercase letters and removing all non-alphanumeric characters, it reads the same forward and backward. Alphanumeric characters include letters and numbers.
//
// Given a string s, return true if it is a palindrome, or false otherwise.
//
// Example 1:
//
// Input: s = "A man, a plan, a canal: Panama"
// Output: true
// Explanation: "amanaplanacanalpanama" is a palindrome.
// Example 2:
//
// Input: s = "race a car"
// Output: false
// Explanation: "raceacar" is not a palindrome.
// Example 3:
//
// Input: s = " "
// Output: true
// Explanation: s is an empty string "" after removing non-alphanumeric characters.
// Since an empty string reads the same forward and backward, it is a palindrome.
package main

import (
	"fmt"
	"unicode"
)

func isAlphanumberic(s rune) bool {
	return unicode.IsLetter(s) || unicode.IsNumber(s)

}

func isPalindrome(s string) bool {
	left, right := 0, len(s)-1

	runes := []rune(s)
	result := true

	for left < right {
		if !isAlphanumberic(runes[left]) {
			left++
			continue
		}

		if !isAlphanumberic(runes[right]) {
			right--
			continue
		}

		if unicode.ToLower(runes[left]) == unicode.ToLower(runes[right]) {
			left++
			right--
		} else {
			result = false
			break
		}
	}

	return result
}

func main() {
	fmt.Println("Hello World")

	s := "A man, a plan, a canal: Panama"
	fmt.Println(isPalindrome(s))

}
