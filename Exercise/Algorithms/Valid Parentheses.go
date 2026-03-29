//Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.
//
//An input string is valid if:
//
//Open brackets must be closed by the same type of brackets.
//Open brackets must be closed in the correct order.
//Every close bracket has a corresponding open bracket of the same type.
//
//Example 1:
//
//Input: s = "()"
//
//Output: true
//
//Example 2:
//
//Input: s = "()[]{}"
//
//Output: true
//
//Example 3:
//
//Input: s = "(]"
//
//Output: false

package main

import "fmt"

func isValid(s string) bool {
	if len(s)%2 != 0 {
		return false
	}

	stack := make([]rune, 0)
	hashMap := make(map[rune]rune)

	hashMap['('] = ')'
	hashMap['['] = ']'
	hashMap['{'] = '}'

	for _, val := range s {
		lenStack := len(stack)

		if val == ')' || val == ']' || val == '}' {

			if lenStack > 0 && hashMap[stack[lenStack-1]] == val {
				stack = stack[0 : lenStack-1]
			} else {
				return false
			}
		} else {
			stack = append(stack, val)
		}

	}

	if len(stack) > 0 {
		return false
	}

	return true

}

func main() {
	fmt.Println("Hello World")

	fmt.Println(isValid("(("))
}
