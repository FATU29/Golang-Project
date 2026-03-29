//Given two strings s and p, return an array of all the start indices of p's anagrams in s. You may return the answer in any order.
//
//Example 1:
//
//Input: s = "cbaebabacd", p = "abc"
//Output: [0,6]
//Explanation:
//The substring wsith start index = 0 is "cba", which is an anagram of "abc".
//The substring with start index = 6 is "bac", which is an anagram of "abc".
//Example 2:
//
//Input: s = "abab", p = "ab"
//Output: [0,1,2]
//Explanation:
//The substring with start index = 0 is "ab", which is an anagram of "ab".
//The substring with start index = 1 is "ba", which is an anagram of "ab".
//The substring with start index = 2 is "ab", which is an anagram of "ab".

package main

import "fmt"

func findAnagrams(s string, p string) []int {
	sizeP := len(p)
	sizeS := len(s)

	if sizeS < sizeP {
		return []int{}
	}

	var mapS, mapP [26]int

	for idx, r := range p {
		mapP[r-'a']++
		mapS[s[idx]-'a']++
	}

	var res []int
	if mapS == mapP {
		res = append(res, 0)
	}

	for i := sizeP; i < sizeS; i++ {
		mapS[s[i]-'a']++
		mapS[s[i-sizeP]-'a']--

		if mapS == mapP {
			res = append(res, i-sizeP+1)
		}
	}

	return res
}

func main() {
	fmt.Println("Hello World")
	s := "cbaebabacd"
	p := "abc"
	fmt.Println(findAnagrams(s, p))
}
